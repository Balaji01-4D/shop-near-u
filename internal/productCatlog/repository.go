package productcatlog

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

func (r *Repository) CreateCatalogProduct(product *models.CatalogProduct) error {
	return r.DB.Create(product).Error
}

func (r *Repository) Suggest(keyword string, limit int) (*[]models.CatalogProduct, error) {
	var products []models.CatalogProduct

	searchPattern := "%" + keyword + "%"
	result := r.DB.
		Limit(limit).
		Where("name LIKE ?", searchPattern).
		Or("brand LIKE ?", searchPattern).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}
	return &products, nil

}
