package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetJadwalPresensi(c *gin.Context) {
	db := database.DB

	var JadwalPresensi []models.JadwalPresensi
	if err := db.Limit(20).Find(&JadwalPresensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, JadwalPresensi)
}
