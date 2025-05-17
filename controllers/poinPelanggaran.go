package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

// GetSiswaByNISN godoc
// @Summary Get siswa by NISN
// @Description Ambil data siswa berdasarkan NISN
// @Tags Siswa
// @Accept  json
// @Produce  json
// @Param nisn path string true "NISN"
// @Success 200 {object} models.Siswa
// @Failure 404 {object} map[string]string
// @Router /siswa/{nisn} [get]
// @Security BearerAuth
func GetPoinPelanggaranbyId(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	var poinPelanggaran models.PoinPelanggaran
	if err := db.Where("id = ?", id).First(&poinPelanggaran).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Poin Pelanggaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, poinPelanggaran)
}

func GetAllPoinPelanggaran(c *gin.Context) {
	db := database.DB

	var poinPelanggaran []models.PoinPelanggaran
	if err := db.Limit(20).Find(&poinPelanggaran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, poinPelanggaran)
}
