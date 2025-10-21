package shop

import (
	"gorm.io/gorm"
	"shop-near-u/internal/models"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}
func (r *Repository) Create(shop *models.Shop) error {
	return r.DB.Create(shop).Error
}

func (r *Repository) FindByEmail(email string) (*models.Shop, error) {
	var shop models.Shop
	result := r.DB.Where("email = ?", email).First(&shop)
	return &shop, result.Error
}
