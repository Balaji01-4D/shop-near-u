package shop

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (ctrl *Controller) RegisterShop(c *gin.Context) {
	var dto ShopRegisterDTORequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shop, err := ctrl.service.RegisterShop(&dto)
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

	c.JSON(201, res)
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	ctrl := NewController(svc)

	shops := r.Group("/shop")
	{
		shops.POST("", ctrl.RegisterShop)
	}
}
