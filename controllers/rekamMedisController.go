package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetRekamMedisBySantriId(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var rekamMedis []models.RekamMedis
	if err := db.Preload("ObatList").Preload("KeluhanList").Preload("StatusRekamMedis").Preload("Perawat").Find(&rekamMedis, "santri_id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekam Medis tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, rekamMedis)
}

func GetRekamMedisById(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var rekamMedis models.RekamMedis
	if err := db.Preload("ObatList").Preload("KeluhanList").Preload("StatusRekamMedis").Preload("Perawat").Last(&rekamMedis, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rekam Medis tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, rekamMedis)
}
