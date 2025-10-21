package product

import (
	"shop-near-u/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) AddProduct(product *models.ShopProduct) error {
	return r.DB.Create(product).Error
}

func (r *Repository) GetProductsByShopID(shopID uint) ([]models.ShopProduct, error) {
	var products []models.ShopProduct
	result := r.DB.Where("shop_id = ?", shopID).Find(&products)
	return products, result.Error
}

func (r *Repository) GetProductByID(productID uint) (*models.ShopProduct, error) {
	var product models.ShopProduct
	result := r.DB.First(&product, productID)
	return &product, result.Error
}

func (r *Repository) UpdateProduct(product *models.ShopProduct) error {
	return r.DB.Save(product).Error
}

func (r *Repository) DeleteProduct(productID uint) error {
	return r.DB.Delete(&models.ShopProduct{}, productID).Error
}
