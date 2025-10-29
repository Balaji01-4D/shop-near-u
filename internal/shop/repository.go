package shop

import (
	"errors"
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

func (r *Repository) FindByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	result := r.DB.First(&shop, id)
	return &shop, result.Error
}

func (r *Repository) UpdateShopStatus(shopID uint, status bool) error {
	return r.DB.Model(&models.Shop{}).Where("id = ?", shopID).Update("is_open", status).Error
}

func (r *Repository) SubscribeShop(shopID uint, userID uint) (uint, error) {
	// Use a transaction to ensure data consistency
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return 0, tx.Error
	}

	// Check if already subscribed
	var existingCount int64
	if err := tx.Model(&models.ShopSubscription{}).Where("shop_id = ? AND user_id = ?", shopID, userID).Count(&existingCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if existingCount > 0 {
		tx.Rollback()
		return 0, errors.New("already subscribed")
	}

	// Create new subscription
	subscription := &models.ShopSubscription{
		ShopID: shopID,
		UserID: userID,
	}
	if err := tx.Create(subscription).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Recalculate subscriber count
	var subscriberCount int64
	if err := tx.Model(&models.ShopSubscription{}).Where("shop_id = ?", shopID).Count(&subscriberCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update shop's subscriber count
	if err := tx.Model(&models.Shop{}).Where("id = ?", shopID).Update("subscriber_count", subscriberCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return uint(subscriberCount), nil
}

func (r *Repository) UnsubscribeShop(shopID uint, userID uint) (uint, error) {
	// Use a transaction to ensure data consistency
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return 0, tx.Error
	}

	// Check if subscription exists
	var existingCount int64
	if err := tx.Model(&models.ShopSubscription{}).Where("shop_id = ? AND user_id = ?", shopID, userID).Count(&existingCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if existingCount == 0 {
		tx.Rollback()
		return 0, errors.New("not subscribed")
	}

	// Delete the subscription
	if err := tx.Where("shop_id = ? AND user_id = ?", shopID, userID).Delete(&models.ShopSubscription{}).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Recalculate subscriber count
	var subscriberCount int64
	if err := tx.Model(&models.ShopSubscription{}).Where("shop_id = ?", shopID).Count(&subscriberCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update shop's subscriber count
	if err := tx.Model(&models.Shop{}).Where("id = ?", shopID).Update("subscriber_count", subscriberCount).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return uint(subscriberCount), nil
}

func (r *Repository) IsAlreadySubscribed(shopID uint, userID uint) bool {
	var count int64
	result := r.DB.Model(&models.ShopSubscription{}).Where("shop_id = ? AND user_id = ?", shopID, userID).Count(&count)
	if result.Error != nil {
		return false
	}
	return count > 0
}

func (r *Repository) GetShopDetails(shopID uint, userID uint) (*models.Shop, bool, error) {
	var shop models.Shop
	if err := r.DB.First(&shop, shopID).Error; err != nil {
		return nil, false, err
	}

	// Check if user is subscribed (if userID is provided)
	var isSubscribed bool
	if userID > 0 {
		var count int64
		if err := r.DB.Model(&models.ShopSubscription{}).Where("user_id = ? AND shop_id = ?", userID, shopID).Count(&count).Error; err != nil {
			return &shop, false, err
		}
		isSubscribed = count > 0
	}

	return &shop, isSubscribed, nil
}
