package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Latitude  float64   `gorm:"type:decimal(10,8);" json:"latitude"`
	Longitude float64   `gorm:"type:decimal(10,8);" json:"longitude"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
