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
	CreatedAt time.Time `json:"created_at"`
}

type Strip struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Title     string    `json:"title"`
	FileURL   string    `json:"file_url"`
	Caption   string    `json:"caption"`
	CreatedAt time.Time `json:"created_at"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Strip{})
}
