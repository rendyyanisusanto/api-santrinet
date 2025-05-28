package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetTatib(c *gin.Context) {
	db := database.DB

	var Tatib []models.Tatib
	if err := db.Limit(20).Find(&Tatib).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, Tatib)
}
