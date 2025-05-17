package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetKamarSantriByID(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var kamarSantri []models.KamarSantri
	if err := db.Preload("Kamar").Preload("Kamar.Asrama").First(&kamarSantri, "santri_id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kamar Santri tidak ditemukan"})
		return
	}

	// Buat response custom
	var response []models.KamarSantriResponse
	for _, ks := range kamarSantri {
		response = append(response, models.KamarSantriResponse{
			ID:          ks.ID,
			KamarID:     ks.KamarID,
			AsramaID:    ks.Kamar.AsramaID,
			NamaKamar:   ks.Kamar.Nama,
			NamaAsrama:  ks.Kamar.Asrama.Nama,
			StatusAktif: ks.StatusAktif,
			Tahun:       ks.Tahun,
		})
	}

	c.JSON(http.StatusOK, response)
}
