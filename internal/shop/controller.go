package shop

import (
	"errors"
	"net/http"
	"os"
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
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shop, err := ctrl.shopService.RegisterShop(&dto)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(409, gin.H{"error": "shop already exists"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if shop == nil {
		c.JSON(500, gin.H{"error": "failed to create shop"})
		return
	}

	token, err := utils.GenerateAccessToken(shop.ID, models.RoleShopOwner)

	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate access token"})
		return
	}

	res := ShopRegisterDTOResponse{
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
	}

	domain := os.Getenv("COOKIE_DOMAIN")

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "/", domain, false, true)

	c.JSON(http.StatusCreated, res)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var dto ShopLoginDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shop, err := ctrl.shopService.AuthenticateShop(&dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if shop == nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := utils.GenerateAccessToken(shop.ID, models.RoleShopOwner)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate access token"})
		return
	}

	res := ShopRegisterDTOResponse{
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
	}

	domain := os.Getenv("COOKIE_DOMAIN")

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24*30, "/", domain, false, true)

	c.JSON(http.StatusOK, res)

}

func (ctrl *Controller) GetShopProfile(c *gin.Context) {
	shopInterface, exists := c.Get("shop")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		c.JSON(500, gin.H{"error": "failed to parse shop data"})
		return
	}

	res := ShopRegisterDTOResponse{
		ID:        shop.ID,
		Name:      shop.Name,
		OwnerName: shop.OwnerName,
		Type:      shop.Type,
		Email:     shop.Email,
		Mobile:    shop.Mobile,
		Address:   shop.Address,
		Latitude:  shop.Latitude,
		Longitude: shop.Longitude,
	}

	c.JSON(http.StatusOK, res)
}

func (ctrl *Controller) AddProduct(c *gin.Context) {
	var dto product.AddProductDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shopInterface, exists := c.Get("shop")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		c.JSON(500, gin.H{"error": "failed to parse shop data"})
		return
	}

	err := ctrl.productService.AddProduct(&dto, shop.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product added successfully"})
}

func (ctrl *Controller) GetAllProducts(c *gin.Context) {
	shopInterface, exists := c.Get("shop")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	shop, ok := shopInterface.(models.Shop)
	if !ok {
		c.JSON(500, gin.H{"error": "failed to parse shop data"})
		return
	}

	products, err := ctrl.productService.GetProductsByShopID(shop.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ctrl *Controller) GetProductByID(c *gin.Context) {
	productIDParam := c.Param("id")
	if productIDParam == "" {
		c.JSON(400, gin.H{"error": "product ID is required"})
		return
	}

	productID, err := utils.ParseUintParam(productIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := ctrl.productService.GetProductByID(productID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if product == nil {
		c.JSON(404, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ctrl *Controller) UpdateProduct(c *gin.Context) {
	var dto product.ProductUpdateDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.productService.UpdateProduct(&dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

func (ctrl *Controller) DeleteProduct(c *gin.Context) {
	productIDParam := c.Param("id")
	if productIDParam == "" {
		c.JSON(400, gin.H{"error": "product ID is required"})
		return
	}

	productID, err := utils.ParseUintParam(productIDParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid product ID"})
		return
	}

	err = ctrl.productService.DeleteProduct(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "product not found"})
			return
		}
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	shopService := NewService(repo)
	productService := product.NewService(product.NewRepository(db))
	ctrl := NewController(shopService, productService)

	shops := r.Group("/shop")
	{
		shops.POST("/register", ctrl.RegisterShop)
		shops.POST("/login", ctrl.Login)
		shops.GET("/profile", middlewares.RequireShopOwnerAuth(db), ctrl.GetShopProfile)

	}

	products := r.Group("/shop/products")
	products.Use(middlewares.RequireShopOwnerAuth(db))
	{
		products.POST("/", ctrl.AddProduct)
		products.GET("/", ctrl.GetAllProducts)
		products.GET("/:id", ctrl.GetProductByID)
		products.PUT("/", ctrl.UpdateProduct)
		products.DELETE("/:id", ctrl.DeleteProduct)
	}
}
