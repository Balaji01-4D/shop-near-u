package models

import (
	"time"
)

type CatalogProduct struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string `gorm:"type:varchar(100);not null;index" json:"name"`
	Brand      string `gorm:"type:varchar(100);index" json:"brand"`
	Category   string `gorm:"type:varchar(100);index" json:"category"`
	Desciption string `gorm:"type:text" json:"description"`

	ImageURL string `gorm:"type:varchar(255)" json:"image_url"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	ShopProducts []ShopProduct `gorm:"foreignKey:CatalogID" json:"shop_products"`
}

type ShopProduct struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID    uint `gorm:"not null;index" json:"shop_id"`
	CatalogID uint `gorm:"not null;index" json:"catalog_id"`

	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int     `gorm:"not null" json:"stock"`
	IsAvailable bool    `gorm:"type:boolean;default:true" json:"is_available"`
	Discount    float64 `gorm:"type:decimal(5,2);default:0" json:"discount"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Shop           Shop           `gorm:"foreignKey:ShopID" json:"shop"`
	CatalogProduct CatalogProduct `gorm:"foreignKey:CatalogID" json:"catalog_product"`
}
