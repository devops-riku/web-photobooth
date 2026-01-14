package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex" json:"username"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

type Strip struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	UserID    *string    `gorm:"index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;references:ID" json:"-"`
	Title     string     `json:"title"`
	FileURL   string     `json:"file_url"`
	Caption   string     `json:"caption"`
	IsGuest   bool       `json:"is_guest"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Strip{})
}
