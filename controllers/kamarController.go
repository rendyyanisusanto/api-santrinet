package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetKamar(c *gin.Context) {
	db := database.DB

	var Kamar []models.Kamar
	if err := db.Limit(100).Find(&Kamar).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, Kamar)
}
