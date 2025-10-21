package productcatlog

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service *Service
}

func NewController(s *Service) *Controller {
	return &Controller{service: s}
}

func (ctrl *Controller) CreateCatalogProduct(c *gin.Context) {
	var dto CreateCatalogProductDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.service.CreateCatalogProduct(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created Successfully",
		"name":    dto.Name,
	})

}

func (ctrl *Controller) SuggestCatalogProducts(c *gin.Context) {

	keyword := c.Query("keyword")
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10 // default limit
	}

	products, err := ctrl.service.SuggestCatalogProducts(keyword, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	ctrl := NewController(svc)
	productCatlogGroup := r.Group("/api/catalog-products")
	{
		productCatlogGroup.POST("/", ctrl.CreateCatalogProduct)
		productCatlogGroup.GET("/suggest", ctrl.SuggestCatalogProducts)
	}
}
