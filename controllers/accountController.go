package controllers

import (
    "net/http"
    "santrinet-api/database"
    "santrinet-api/models"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

// GET /account - get current authenticated user's profile (without password)
func GetAccount(c *gin.Context) {
    db := database.DB

    // get user_id from auth middleware context
    uidVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.Users
    if err := db.First(&user, uidVal).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
        return
    }

    // fetch linked santri by anggota_id (if exists)
    var santri *models.Santri
    if user.AnggotaID != 0 {
        var s models.Santri
        if err := db.First(&s, user.AnggotaID).Error; err == nil {
            santri = &s
        }
    }

    // organize permissions using existing helper in loginController.go
    // perms := organizePermissions(user.Groups)

    c.JSON(http.StatusOK, gin.H{
        "user":        user,   // Password field already omitted via json:"-"
        "santri":      santri, // can be null if not linked
        // "permissions": perms,  // grouped by menu and submenu
    })
}

type UpdateAccountRequest struct {
    Email *string `json:"email"`
}

// PUT /account - update current authenticated user's profile (selected fields)
func UpdateAccount(c *gin.Context) {
    db := database.DB

    uidVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var req UpdateAccountRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
        return
    }

    var user models.Users
    if err := db.First(&user, uidVal).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
        return
    }

    // update only provided fields (email only as per schema)
    updates := make(map[string]interface{})
    if req.Email != nil { updates["email"] = *req.Email }

    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Tidak ada field yang diupdate"})
        return
    }

    if err := db.Model(&user).Updates(updates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate akun"})
        return
    }

    // reload with relations (optional)
    db.Preload("Groups.GroupsDetail.Submenu.Menu").First(&user, user.ID)

    perms := organizePermissions(user.Groups)

    c.JSON(http.StatusOK, gin.H{
        "message":     "Akun berhasil diupdate",
        "user":        user,
        "permissions": perms,
    })
}

type ChangePasswordRequest struct {
    CurrentPassword string `json:"current_password" binding:"required"`
    NewPassword     string `json:"new_password" binding:"required,min=6"`
}

// PUT /account/password - change current user's password
func ChangePassword(c *gin.Context) {
    db := database.DB

    uidVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var req ChangePasswordRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.Users
    if err := db.First(&user, uidVal).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
        return
    }

    // verify current password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Password saat ini salah"})
        return
    }

    // hash new password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah password"})
        return
    }

    if err := db.Model(&user).Update("password", string(hashed)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan password baru"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diubah"})
}
