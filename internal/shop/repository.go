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
