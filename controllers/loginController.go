package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/middleware"
	"santrinet-api/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	db := database.DB

	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
		return
	}

	var user models.Users
	if err := db.Preload("Groups").Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username tidak ditemukan"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	role := "user"
	if len(user.Groups) > 0 {
		role = user.Groups[0].Name
	}

	// Generate Access Token
	accessToken, err := middleware.GenerateAccessToken(user.ID, role, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token akses"})
		return
	}

	// Generate Refresh Token
	refreshToken, err := middleware.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat refresh token"})
		return
	}

	// Simpan refresh token di cookie
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "", false, true) // secure=true kalau pakai HTTPS

	c.JSON(http.StatusOK, gin.H{
		"token":   accessToken,
		"user":    user,
		"group":   role,
		"user_id": user.ID,
	})
}
func RefreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	// Parse token dari cookie
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return middleware.SecretKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Ambil user ID dari subject
	userID, _ := strconv.Atoi(claims.Subject)

	// Ambil user dari database
	db := database.DB
	var user models.Users
	if err := db.Preload("Groups").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Ambil role dan username dari user
	role := "user"
	if len(user.Groups) > 0 {
		role = user.Groups[0].Name
	}
	username := user.Username

	// Generate access token baru
	newToken, err := middleware.GenerateAccessToken(uint(userID), role, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newToken})
}
