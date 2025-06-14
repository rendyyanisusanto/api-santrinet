package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetPelanggaranById(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var pelanggaran []models.Pelanggaran
	if err := db.
		Preload("Santri").
		Where("santri_id = ? AND (status_pengajuan = ? OR status_pengajuan = ?)", id, "DITERIMA", "BUKAN PENGAJUAN").
		Order("id DESC").
		Find(&pelanggaran).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggaran Santri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pelanggaran)
}

func GetDetailPelanggaranById(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var pelanggaran models.Pelanggaran
	if err := db.Preload("Santri").Preload("Tatib.KategoriTatib").Preload("Tatib").Preload("Pengurus").Preload("Pengurus.Santri").Preload("Pengurus.LembagaPengurus").Order("id DESC").First(&pelanggaran, "id= ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggaran Santri tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, pelanggaran)
}

func GetTotalPelanggaranBySantriID(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	var count int64
	if err := db.Model(&models.Pelanggaran{}).
		Where("santri_id = ? AND (status_pengajuan = ? OR status_pengajuan = ?)", id, "DITERIMA", "BUKAN PENGAJUAN").
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung pelanggaran"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"santri_id": id,
		"total":     count,
	})
}

func GetTotalPelanggaranHariIniBySantriID(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	today := time.Now().Format("2006-01-02") // format yyyy-mm-dd

	var count int64
	if err := db.Model(&models.Pelanggaran{}).
		Where("santri_id = ? AND DATE(tanggal) = ? AND (status_pengajuan = ? OR status_pengajuan = ?)", id, today, "DITERIMA", "BUKAN PENGAJUAN").
		Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung pelanggaran hari ini"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"santri_id": id,
		"tanggal":   today,
		"total":     count,
	})
}
