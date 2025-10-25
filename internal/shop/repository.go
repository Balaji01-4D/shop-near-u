package shop

import (
	"shop-near-u/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) FindNearbyShops(lat float64, lon float64, radius float64, limit int) ([]NearByShopsDTORespone, error) {
	var shops []NearByShopsDTORespone

	query := `
        SELECT 
            id, 
            name, 
            address, 
            latitude, 
            longitude, 
            ST_Distance(location, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography) AS distance
        FROM shops
        WHERE ST_DWithin(location, ST_SetSRID(ST_MakePoint(?, ?), 4326)::geography, ?)
        ORDER BY distance
        LIMIT ?
    `

	result := r.DB.Raw(query, lon, lat, lon, lat, radius, limit).Scan(&shops)

	if result.Error != nil {
		return nil, result.Error
	}

	return shops, nil
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
