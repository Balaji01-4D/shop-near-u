package shop

import (
	"errors"
	"net/http"
	"shop-near-u/internal/middlewares"
	"shop-near-u/internal/models"
	"shop-near-u/internal/product"
	"shop-near-u/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	shopService    *Service
	productService *product.Service
}

func NewController(s *Service, p *product.Service) *Controller {
	return &Controller{shopService: s, productService: p}
}

func (ctrl *Controller) RegisterShop(c *gin.Context) {
	var dto ShopRegisterDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponseSimple(c, 400, err.Error())
		return
	}

	shop, err := ctrl.shopService.RegisterShop(&dto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.ErrorResponseSimple(c, 409, "shop already exists")
			return
		}
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	if shop == nil {
		utils.ErrorResponseSimple(c, 500, "failed to create shop")
		return
	}

	token, err := utils.GenerateAccessToken(shop.ID, models.RoleShopOwner)

	if err != nil {
		utils.ErrorResponseSimple(c, 500, "failed to generate access token")
		return
	}

	utils.SetCookie(token, 3600*24*30, c)

	utils.SuccessResponse(c, http.StatusCreated, "Shop registered successfully", ShopRegisterDTOResponse{
		ID:        shop.ID,
		Name:      shop.Name,
		OwnerName: shop.OwnerName,
		Type:      shop.Type,
		Email:     shop.Email,
		Mobile:    shop.Mobile,
		Address:   shop.Address,
		Latitude:  shop.Latitude,
		Longitude: shop.Longitude,
		Token:     token,
	})
}

func (ctrl *Controller) Login(c *gin.Context) {
	var dto ShopLoginDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponseSimple(c, 400, err.Error())
		return
	}

	shop, err := ctrl.shopService.AuthenticateShop(&dto)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	if shop == nil {
		utils.ErrorResponseSimple(c, 401, "invalid credentials")
		return
	}

	token, err := utils.GenerateAccessToken(shop.ID, models.RoleShopOwner)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, "failed to generate access token")
		return
	}

	utils.SetCookie(token, 3600*24*30, c)

	utils.SuccessResponse(c, http.StatusOK, "Shop logged in successfully", ShopRegisterDTOResponse{
		ID:        shop.ID,
		Name:      shop.Name,
		OwnerName: shop.OwnerName,
		Type:      shop.Type,
		Email:     shop.Email,
		Mobile:    shop.Mobile,
		Address:   shop.Address,
		Latitude:  shop.Latitude,
		Longitude: shop.Longitude,
		Token:     token,
		IsOpen:    shop.IsOpen,
	})

}

func (ctrl *Controller) GetShopProfile(c *gin.Context) {
	shopInterface, exists := c.Get("shop")
	if !exists {
		utils.ErrorResponseSimple(c, 401, "unauthorized")
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		utils.ErrorResponseSimple(c, 500, "failed to parse shop data")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Shop profile retrieved successfully", ShopRegisterDTOResponse{
		ID:              shop.ID,
		Name:            shop.Name,
		OwnerName:       shop.OwnerName,
		Type:            shop.Type,
		Email:           shop.Email,
		Mobile:          shop.Mobile,
		Address:         shop.Address,
		Latitude:        shop.Latitude,
		Longitude:       shop.Longitude,
		SubscriberCount: shop.SubscriberCount,
		IsOpen:          shop.IsOpen,
	})
}

func (ctrl *Controller) AddProduct(c *gin.Context) {
	var dto product.AddProductDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponseSimple(c, 400, err.Error())
		return
	}

	shopInterface, exists := c.Get("shop")
	if !exists {
		utils.ErrorResponseSimple(c, 401, "unauthorized")
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		utils.ErrorResponseSimple(c, 500, "failed to parse shop data")
		return
	}

	err := ctrl.productService.AddProduct(&dto, shop.ID)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product added successfully", nil)
}

func (ctrl *Controller) GetAllProducts(c *gin.Context) {
	shopInterface, exists := c.Get("shop")
	if !exists {
		utils.ErrorResponseSimple(c, 401, "unauthorized")
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		utils.ErrorResponseSimple(c, 500, "failed to parse shop data")
		return
	}

	products, err := ctrl.productService.GetProductsByShopID(shop.ID)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

func (ctrl *Controller) GetProductByID(c *gin.Context) {
	productIDParam := c.Param("id")
	if productIDParam == "" {
		utils.ErrorResponseSimple(c, 400, "product ID is required")
		return
	}

	productID, err := utils.ParseUintParam(productIDParam)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid product ID")
		return
	}

	product, err := ctrl.productService.GetProductByID(productID)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	if product == nil {
		utils.ErrorResponseSimple(c, 404, "product not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

func (ctrl *Controller) UpdateProduct(c *gin.Context) {
	var dto product.ProductUpdateDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponseSimple(c, 400, err.Error())
		return
	}

	err := ctrl.productService.UpdateProduct(&dto)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", nil)
}

func (ctrl *Controller) DeleteProduct(c *gin.Context) {
	productIDParam := c.Param("id")
	if productIDParam == "" {
		utils.ErrorResponseSimple(c, 400, "product ID is required")
		return
	}

	productID, err := utils.ParseUintParam(productIDParam)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid product ID")
		return
	}

	err = ctrl.productService.DeleteProduct(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponseSimple(c, 404, "product not found")
			return
		}
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

func (ctrl *Controller) IsShopOpen(c *gin.Context) {
	shopIDParam := c.Param("id")
	shopId, err := utils.ParseUintParam(shopIDParam)

	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid shop ID")
		return
	}

	shop, err := ctrl.shopService.GetShopByID(shopId)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Shop status retrieved successfully", shop.IsOpen)
}

func (ctrl *Controller) UpdateShopStatus(c *gin.Context) {
	shopInterface, exists := c.Get("shop")
	if !exists {
		utils.ErrorResponseSimple(c, 401, "unauthorized")
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		utils.ErrorResponseSimple(c, 500, "failed to parse shop data")
		return
	}

	status := c.Query("status")
	// status can be "open" or "closed"
	if status == "" {
		utils.ErrorResponseSimple(c, 400, "status is required")
		return
	}

	var isOpen bool
	switch status {
	case "open":
		isOpen = true
	case "closed":
		isOpen = false
	default:
		utils.ErrorResponseSimple(c, 400, "invalid status")
		return
	}

	err := ctrl.shopService.UpdateShopStatus(shop.ID, isOpen)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Shop status updated successfully", nil)
}

func (ctrl *Controller) NearByShop(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radStr := c.DefaultQuery("radius", "5000")
	limit := c.DefaultQuery("limit", "10")

	lat, err := utils.ParseFloatParam(latStr)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid latitude")
		return
	}

	lim, err := utils.ParseIntParam(limit)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid limit")
		return
	}

	lon, err := utils.ParseFloatParam(lonStr)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid longitude")
		return
	}

	radius, err := utils.ParseFloatParam(radStr)
	if err != nil {
		utils.ErrorResponseSimple(c, 400, "invalid radius")
		return
	}

	shops, err := ctrl.shopService.GetNearbyShops(lat, lon, radius, lim)
	if err != nil {
		utils.ErrorResponseSimple(c, 500, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Nearby shops retrieved successfully", shops)
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	shopService := NewService(repo)
	productService := product.NewService(product.NewRepository(db))
	ctrl := NewController(shopService, productService)

	shops := r.Group("/shops")
	{
		shops.POST("/register", ctrl.RegisterShop)
		shops.POST("/login", ctrl.Login)
		shops.GET("/profile", middlewares.RequireShopOwnerAuth(db), ctrl.GetShopProfile)
		shops.GET("", ctrl.NearByShop)
		shops.GET("/is_open/:id", ctrl.IsShopOpen)
		shops.PUT("/status", middlewares.RequireShopOwnerAuth(db), ctrl.UpdateShopStatus)

		shops.GET("/:id", middlewares.RequireUserAuth(db), ctrl.GetShopDetails)
		shops.POST("/:id/subscribe", middlewares.RequireUserAuth(db), ctrl.SubscribeShop)
		shops.POST("/:id/unsubscribe", middlewares.RequireUserAuth(db), ctrl.UnsubscribeShop)
	}

	products := r.Group("/shop/products")
	products.Use(middlewares.RequireShopOwnerAuth(db))
	{
		products.POST("", ctrl.AddProduct)
		products.GET("", ctrl.GetAllProducts)
		products.GET("/:id", ctrl.GetProductByID)
		products.PUT("", ctrl.UpdateProduct)
		products.DELETE("/:id", ctrl.DeleteProduct)
	}
}
