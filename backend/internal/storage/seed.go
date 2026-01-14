package storage

import (
	"log"
	"web-photobooth/backend/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedSuperUser(db *gorm.DB) {
	var count int64
	// Check if admin exists
	db.Model(&models.User{}).Where("is_admin = ?", true).Count(&count)
	if count > 0 {
		return // Admin already exists
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte("Admin123"), bcrypt.DefaultCost)

	admin := models.User{
		ID:       uuid.New().String(),
		Username: "superuser",
		Email:    "wuby@superuser.com",
		Password: string(hashed),
		IsAdmin:  true,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to seed superuser: %v", err)
	} else {
		log.Println("Superuser seeded successfully: wuby@superuser.com / Admin123")
	}
}
