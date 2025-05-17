package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

func GetAllSantri(c *gin.Context) {
	db := database.DB

	var santri []models.Santri
	if err := db.Limit(20).Find(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, santri)
}

func GetSantriByID(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var santri models.Santri
	if err := db.First(&santri, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Santri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, santri)
}

func GetSantriDetailByID(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	// 1. Ambil data santri
	var santri models.Santri
	if err := db.First(&santri, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Santri tidak ditemukan",
		})
		return
	}

	// 2. Ambil satu data kamar santri (misalnya status aktif)
	var kamarSantri models.KamarSantri
	if err := db.Preload("Kamar").Preload("Kamar.Asrama").
		Where("santri_id = ? AND status_aktif = ?", id, "1").
		First(&kamarSantri).Error; err != nil {
		// Kamar santri tidak wajib? Bisa kosong
		c.JSON(http.StatusOK, gin.H{
			"santri":       santri,
			"kamar_santri": nil,
		})
		return
	}

	// 3. Format ke struct ringan (KamarSantriResponse)
	response := models.KamarSantriResponse{
		ID:          kamarSantri.ID,
		KamarID:     kamarSantri.KamarID,
		AsramaID:    kamarSantri.Kamar.AsramaID,
		NamaKamar:   kamarSantri.Kamar.Nama,
		NamaAsrama:  kamarSantri.Kamar.Asrama.Nama,
		StatusAktif: kamarSantri.StatusAktif,
		Tahun:       kamarSantri.Tahun,
	}

	// 4. Return JSON gabungan
	c.JSON(http.StatusOK, gin.H{
		"santri":       santri,
		"kamar_santri": response,
	})
}
