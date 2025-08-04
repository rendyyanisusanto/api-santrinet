package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllPengurus(c *gin.Context) {
	db := database.DB

	var pengurus []models.Pengurus
	if err := db.
		Preload("Santri").
		Preload("LembagaPengurus").
		Find(&pengurus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data pengurus"})
		return
	}

	c.JSON(http.StatusOK, pengurus)
}
