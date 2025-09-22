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

// Struktur untuk response permission yang terorganisir
type MenuPermission struct {
	MenuID   uint                `json:"menu_id"`
	MenuName string              `json:"menu_name"`
	Submenus []SubmenuPermission `json:"submenus"`
}

type SubmenuPermission struct {
	SubmenuID   uint   `json:"submenu_id"`
	SubmenuName string `json:"submenu_name"`
	Permissions CRUD   `json:"permissions"`
}

type CRUD struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

// Fungsi untuk mengorganisir permission berdasarkan menu dan submenu
func organizePermissions(groups []models.Group) []MenuPermission {
	menuMap := make(map[uint]*MenuPermission)

	for _, group := range groups {
		for _, groupDetail := range group.GroupsDetail {
			menuID := groupDetail.Submenu.MenuID

			// Buat menu baru jika belum ada
			if _, exists := menuMap[menuID]; !exists {
				menuMap[menuID] = &MenuPermission{
					MenuID:   menuID,
					MenuName: groupDetail.Submenu.Menu.Menu,
					Submenus: []SubmenuPermission{},
				}
			}

			// Tambahkan submenu permission
			submenuPerm := SubmenuPermission{
				SubmenuID:   groupDetail.SubmenuID,
				SubmenuName: groupDetail.Submenu.Submenu,
				Permissions: CRUD{
					Create: groupDetail.C == 1,
					Read:   groupDetail.R == 1,
					Update: groupDetail.U == 1,
					Delete: groupDetail.D == 1,
				},
			}

			menuMap[menuID].Submenus = append(menuMap[menuID].Submenus, submenuPerm)
		}
	}

	// Convert map ke slice
	var result []MenuPermission
	for _, menu := range menuMap {
		result = append(result, *menu)
	}

	return result
}

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

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
	})
}

func GetGroupPermissions(c *gin.Context) {
	db := database.DB

	// Ambil group_id dari parameter URL
	groupIDParam := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Group ID tidak valid"})
		return
	}

	// Ambil group dari database dengan preload relasi
	var group models.Group
	if err := db.Preload("GroupsDetail.Submenu.Menu").First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group tidak ditemukan"})
		return
	}

	// Filter hanya GroupsDetail yang memiliki permission (bukan 0 semua)
	var filteredGroupsDetail []models.GroupsDetail
	for _, detail := range group.GroupsDetail {
		// Hanya ambil jika minimal ada satu permission yang true (1)
		if detail.C == 1 || detail.R == 1 || detail.U == 1 || detail.D == 1 {
			filteredGroupsDetail = append(filteredGroupsDetail, detail)
		}
	}

	// Update group dengan GroupsDetail yang sudah difilter
	group.GroupsDetail = filteredGroupsDetail

	// Struktur permission data berdasarkan menu dan submenu untuk group ini
	groups := []models.Group{group}
	permissions := organizePermissions(groups)

	c.JSON(http.StatusOK, gin.H{
		"group_id":    group.ID,
		"group_name":  group.Name,
		"permissions": permissions,
	})
}
