package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"web-photobooth/backend/internal/middleware"
	"web-photobooth/backend/internal/models"
)

type Handler struct {
	DB                  *gorm.DB
	S3Client            *s3.Client
	Bucket              string
	JWTSecret           string
	GuestExpirationDays int
}

func NewHandler(db *gorm.DB, s3Client *s3.Client, bucket string, jwtSecret string, guestDays int) *Handler {
	return &Handler{
		DB:                  db,
		S3Client:            s3Client,
		Bucket:              bucket,
		JWTSecret:           jwtSecret,
		GuestExpirationDays: guestDays,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/signup", h.Signup)
			auth.POST("/login", h.Login)
			auth.POST("/sync", h.SyncUser)
		}

		strips := api.Group("/strips")
		{
			// Public routes (no auth)
			strips.POST("/guest-save", h.GuestSaveStrip)
			strips.GET("/public/:id", h.GetPublicStrip)

			// Protected routes
			protected := strips.Group("/")
			protected.Use(middleware.AuthMiddleware(h.JWTSecret))
			{
				protected.POST("/save", h.SaveStrip)
				protected.GET("/my-strips", h.GetMyStrips)
				protected.PATCH("/:id", h.UpdateStrip)
				protected.DELETE("/:id", h.DeleteStrip)
			}
		}
	}
}

func (h *Handler) Signup(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Double check duplicates
	var existing models.User
	if err := h.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}
	if err := h.DB.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel password"})
		return
	}

	user := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.DB.Where("email = ? OR username = ?", req.Identifier, req.Identifier).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *Handler) SyncUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User sync successful"})
}

func (h *Handler) SaveStrip(c *gin.Context) {
	userID := c.GetString("user_id")
	var req struct {
		ID      string `json:"id"`
		Image   string `json:"image"` // Base64
		Title   string `json:"title"`
		Caption string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Decode
	b64data := req.Image
	if idx := strings.Index(b64data, ","); idx != -1 {
		b64data = b64data[idx+1:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
		return
	}

	// 2. Prepare IDs and Filename
	stripID := req.ID
	if stripID == "" {
		stripID = uuid.New().String()
	}
	fileName := fmt.Sprintf("strips/%s/%s.png", userID, stripID)

	// 3. Upload to S3 (DigitalOcean Spaces)
	_, err = h.S3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imgBytes),
		ContentType: aws.String("image/png"),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		log.Printf("S3 Upload Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage"})
		return
	}

	// 3. Construct File URL using CDN
	fileURL := fmt.Sprintf("https://%s.sgp1.cdn.digitaloceanspaces.com/%s", h.Bucket, fileName)

	// 4. Save to DB (already have stripID)
	log.Printf("SaveStrip: userID=%s, stripID=%s, fileURL=%s", userID, stripID, fileURL)
	// 4. Save to DB (already have stripID)
	uid := userID
	log.Printf("SaveStrip: Writing to DB - userID=%s, stripID=%s", uid, stripID)
	
	strip := models.Strip{
		ID:        stripID,
		UserID:    &uid,
		Title:     req.Title,
		FileURL:   fileURL,
		Caption:   req.Caption,
		CreatedAt: time.Now(),
	}

	if err := h.DB.Create(&strip).Error; err != nil {
		log.Printf("SaveStrip DB Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}
	
	uidVal := "NIL"
	if strip.UserID != nil {
		uidVal = *strip.UserID
	}
	log.Printf("SaveStrip SUCCESS: id=%s, stored_user_id=%s", strip.ID, uidVal)

	// Return ID so frontend can update it later
	c.JSON(http.StatusOK, gin.H{
		"message":  "Strip saved successfully",
		"file_url": fileURL,
		"id":       strip.ID,
	})
}

func (h *Handler) GuestSaveStrip(c *gin.Context) {
	var req struct {
		ID      string `json:"id"`
		Image   string `json:"image"` // Base64
		Title   string `json:"title"`
		Caption string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 1. Decode
	b64data := req.Image
	if idx := strings.Index(b64data, ","); idx != -1 {
		b64data = b64data[idx+1:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to decode image"})
		return
	}

	// 2. Prepare IDs and Filename
	log.Printf("Incoming request ID: %s", req.ID)
	stripID := req.ID
	if stripID == "" {
		stripID = uuid.New().String()
		log.Printf("No ID provided, generated new: %s", stripID)
	}
	fileName := fmt.Sprintf("strips/guest/%s.png", stripID)
	log.Printf("Final fileName for upload: %s", fileName)

	// 3. Upload to S3
	_, err = h.S3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imgBytes),
		ContentType: aws.String("image/png"),
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage"})
		return
	}

	// 3. Construct File URL
	fileURL := fmt.Sprintf("https://%s.sgp1.cdn.digitaloceanspaces.com/%s", h.Bucket, fileName)

	// 4. Expiration
	expiresAt := time.Now().AddDate(0, 0, h.GuestExpirationDays)

	// 6. Save to DB (already have stripID)

	strip := models.Strip{
		ID:        stripID,
		UserID:    nil,
		Title:     req.Title,
		FileURL:   fileURL,
		Caption:   req.Caption,
		IsGuest:   true,
		ExpiresAt: &expiresAt,
		CreatedAt: time.Now(),
	}

	if err := h.DB.Create(&strip).Error; err != nil {
		log.Printf("DATABASE ERROR in GuestSaveStrip: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Guest strip saved successfully",
		"file_url":   fileURL,
		"id":         strip.ID,
		"expires_at": expiresAt,
	})
}

func (h *Handler) GetPublicStrip(c *gin.Context) {
	id := c.Param("id")
	var strip models.Strip

	if err := h.DB.First(&strip, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Memory not found"})
		return
	}

	// Double check expiration for guests
	if strip.IsGuest && strip.ExpiresAt != nil && time.Now().After(*strip.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "This memory has expired"})
		return
	}

	c.JSON(http.StatusOK, strip)
}

func (h *Handler) GetMyStrips(c *gin.Context) {
	userID := c.GetString("user_id")
	var strips []models.Strip
	if err := h.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&strips).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch strips"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"strips": strips})
}

func (h *Handler) UpdateStrip(c *gin.Context) {
	userID := c.GetString("user_id")
	stripID := c.Param("id")

	var req struct {
		Title   string `json:"title"`
		Caption string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var strip models.Strip
	log.Printf("UpdateStrip lookup: id=%s, user_id=%s", stripID, userID)
	
	// 1. Find the strip regardless of owner first
	if err := h.DB.Where("id = ?", stripID).First(&strip).Error; err != nil {
		log.Printf("UpdateStrip NOT FOUND: id=%s", stripID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Strip not found"})
		return
	}

	// 2. Ownership & Claiming Logic
	if strip.UserID != nil && *strip.UserID != userID {
		log.Printf("UpdateStrip FORBIDDEN: id=%s attempted update by %s but owned by %s", stripID, userID, *strip.UserID)
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to update this memory"})
		return
	}

	// 3. If it was a guest strip, claim it for this user
	if strip.UserID == nil {
		log.Printf("UpdateStrip: Claiming guest strip %s for user %s", stripID, userID)
		strip.UserID = &userID
		strip.IsGuest = false
		strip.ExpiresAt = nil // Permanent save
	}

	log.Printf("UpdateStrip PROCEEDING: id=%s, for user=%s", strip.ID, userID)

	if req.Title != "" {
		strip.Title = req.Title
	}
	if req.Caption != "" {
		strip.Caption = req.Caption
	}

	if err := h.DB.Save(&strip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update strip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strip updated", "strip": strip})
}

func (h *Handler) DeleteStrip(c *gin.Context) {
	userID := c.GetString("user_id")
	stripID := c.Param("id")

	// 1. Get the strip to find the file path
	var strip models.Strip
	if err := h.DB.Where("id = ? AND user_id = ?", stripID, userID).First(&strip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Strip not found"})
		return
	}

	// 2. Delete from S3 (DigitalOcean Spaces)
	if strip.FileURL != "" {
		// Extract the key from the URL
		// URL format: https://<bucket>.<endpoint>/<key>
		// We can try to parse the URL
		u, err := url.Parse(strip.FileURL)
		if err == nil {
			// Path usually starts with /, trim it
			key := strings.TrimPrefix(u.Path, "/")
			
			// Safety check: key shouldn't be empty
			if key != "" {
				_, err := h.S3Client.DeleteObject(c.Request.Context(), &s3.DeleteObjectInput{
					Bucket: aws.String(h.Bucket),
					Key:    aws.String(key),
				})
				if err != nil {
					log.Printf("Failed to delete S3 object %s: %v", key, err)
					// We continue to delete from DB even if S3 fails, to keep DB consistent with user intent
				} else {
					log.Printf("Deleted S3 object: %s", key)
				}
			}
		}
	}

	// 3. Delete from DB
	if err := h.DB.Delete(&strip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete strip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strip deleted"})
}

func (h *Handler) CleanupExpiredStrips() {
	var expiredStrips []models.Strip
	now := time.Now()

	// Find expired guest strips
	if err := h.DB.Where("is_guest = ? AND expires_at < ?", true, now).Find(&expiredStrips).Error; err != nil {
		log.Printf("Cleanup Error (Find): %v", err)
		return
	}

	if len(expiredStrips) == 0 {
		return
	}

	log.Printf("Found %d expired guest strips to clean up", len(expiredStrips))

	for _, strip := range expiredStrips {
		// 1. Extract Key from URL
		// URL Format: https://{bucket}.sgp1.cdn.digitaloceanspaces.com/strips/guest/{uuid}.png
		parts := strings.SplitN(strip.FileURL, ".com/", 2)
		if len(parts) < 2 {
			log.Printf("Cleanup Warning: Could not parse key from URL %s", strip.FileURL)
			continue
		}
		key := parts[1]

		// 2. Delete from S3
		_, err := h.S3Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
			Bucket: aws.String(h.Bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			log.Printf("Cleanup Error (S3 Delete %s): %v", key, err)
			// Continue to DB delete anyway to avoid getting stuck, or omit? 
			// Usually better to delete DB if S3 is gone or failed 404
		}

		// 3. Delete from DB
		if err := h.DB.Delete(&strip).Error; err != nil {
			log.Printf("Cleanup Error (DB Delete ID %s): %v", strip.ID, err)
		} else {
			log.Printf("Successfully cleaned up expired strip: %s", key)
		}
	}
}
