package models

import (
	"time"
	"github.com/restayway/gogis"
)


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
	Location      gogis.Point `gorm:"type:geometry(POINT,4326);" json:"location"`
	SubscriberCount uint        `gorm:"type:int;default:0" json:"subscriber_count"`
	IsOpen        bool          `gorm:"type:boolean;default:true" json:"is_open"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	ShopProducts []ShopProduct `gorm:"foreignKey:ShopID" json:"shop_products"`
}


type ShopSubscription struct {
	ID 	  uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShopID uint      `gorm:"not null;index" json:"shop_id"`
	UserID uint      `gorm:"not null;index" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}