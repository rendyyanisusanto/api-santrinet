package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllIjinPengurus(c *gin.Context) {
	db := database.DB

	var data []models.IjinPengurus

	if err := db.
		Preload("Santri").
		Order("id DESC").
		Find(&data).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data perizinan pengurus"})
		return
	}

	c.JSON(http.StatusOK, data)
}
func GetPengurusSedangKeluar(c *gin.Context) {
	db := database.DB

	var data []models.IjinPengurus

	if err := db.
		Preload("Santri").
		Where("waktu_kembali IS NULL").
		Where("status = ?", "KELUAR").
		Order("waktu_keluar DESC").
		Find(&data).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pengurus keluar"})
		return
	}

	c.JSON(http.StatusOK, data)
}
