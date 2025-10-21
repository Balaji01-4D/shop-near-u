package models

import "time"

type Shop struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Name             string `gorm:"type:varchar(100);not null" json:"name"`
	OwnerName        string `gorm:"type:varchar(100);not null" json:"owner_name"`
	Email            string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Mobile           string `gorm:"type:varchar(15);not null" json:"mobile"`
	Type             string `gorm:"type:varchar(50);not null" json:"type"`
	SupportsDelivery bool   `gorm:"type:boolean;default:false" json:"supports_delivery"`
	Password         string `gorm:"type:varchar(255);not null" json:"-"`

	Address   string  `gorm:"type:varchar(255);not null" json:"address"`
	Latitude  float64 `gorm:"type:decimal(10,8);" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(10,8);" json:"longitude"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	ShopProducts []ShopProduct `gorm:"foreignKey:ShopID" json:"shop_products"`
}
