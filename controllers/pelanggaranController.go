package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetPelanggaranById(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var pelanggaran []models.Pelanggaran
	if err := db.Preload("Santri").Preload("Tatib").Preload("Pengurus").Preload("Pengurus.Santri").Preload("Pengurus.LembagaPengurus").Order("id DESC").Find(&pelanggaran, "santri_id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggaran Santri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pelanggaran)
}

func GetDetailPelanggaranById(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var pelanggaran models.Pelanggaran
	if err := db.Preload("Santri").Preload("Tatib").Preload("Pengurus").Preload("Pengurus.Santri").Preload("Pengurus.LembagaPengurus").Order("id DESC").First(&pelanggaran, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggaran Santri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pelanggaran)
}
