package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllPengasuh(c *gin.Context) {
	db := database.DB

	var pengasuh []models.Pengasuh
	if err := db.Find(&pengasuh).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, pengasuh)
}

func GetPengasuhByID(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var pengasuh models.Pengasuh
	if err := db.First(&pengasuh, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengasuh tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pengasuh)
}
