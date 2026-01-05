package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"web-photobooth/backend/internal/middleware"
	"web-photobooth/backend/internal/models"
)

type Handler struct {
	DB        *gorm.DB
	S3Client  *s3.Client
	Bucket    string
	JWTSecret string
}

func NewHandler(db *gorm.DB, s3Client *s3.Client, bucket string, jwtSecret string) *Handler {
	return &Handler{
		DB:        db,
		S3Client:  s3Client,
		Bucket:    bucket,
		JWTSecret: jwtSecret,
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
		strips.Use(middleware.AuthMiddleware(h.JWTSecret))
		{
			strips.POST("/save", h.SaveStrip)
			strips.GET("/my-strips", h.GetMyStrips)
			strips.PATCH("/:id", h.UpdateStrip)
			strips.DELETE("/:id", h.DeleteStrip)
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

	// 2. Upload to S3 (DigitalOcean Spaces)
	fileName := fmt.Sprintf("strips/%s/%s.png", userID, uuid.New().String())
	
	_, err = h.S3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.Bucket),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(imgBytes),
		ContentType: aws.String("image/png"),
		ACL:         "public-read",
	})

	if err != nil {
		log.Printf("S3 Upload Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to storage"})
		return
	}

	// 3. Construct File URL using CDN
	fileURL := fmt.Sprintf("https://%s.sgp1.cdn.digitaloceanspaces.com/%s", h.Bucket, fileName)

	// 4. Save to DB
	strip := models.Strip{
		UserID:    userID,
		Title:     req.Title,
		FileURL:   fileURL,
		Caption:   req.Caption,
		CreatedAt: time.Now(),
	}

	if err := h.DB.Create(&strip).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	// Return ID so frontend can update it later
	c.JSON(http.StatusOK, gin.H{
		"message":  "Strip saved successfully",
		"file_url": fileURL,
		"id":       strip.ID,
	})
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
	if err := h.DB.Where("id = ? AND user_id = ?", stripID, userID).First(&strip).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Strip not found"})
		return
	}

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

	if err := h.DB.Where("id = ? AND user_id = ?", stripID, userID).Delete(&models.Strip{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete strip"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Strip deleted"})
}
