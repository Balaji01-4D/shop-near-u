package productcatlog

import (
	"shop-near-u/internal/models"
	"strings"

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

	// Convert keyword to lowercase for case-insensitive search
	searchPattern := "%" + strings.ToLower(keyword) + "%"

	// Query with enhanced search across multiple fields
	result := r.DB.
		Limit(limit).
		Where("LOWER(name) LIKE ? OR LOWER(brand) LIKE ? OR LOWER(category) LIKE ? OR LOWER(desciption) LIKE ?",
					searchPattern, searchPattern, searchPattern, searchPattern).
		Order("name ASC"). // Order by name for consistent results
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}
	return &products, nil
}
