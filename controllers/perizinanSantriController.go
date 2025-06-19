// controllers/pelanggaran_controller.go
package controllers

import (
	"fmt"
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"
	"strconv"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func generateKodeIzin() string {
	rand.Seed(time.Now().UnixNano())
	timestamp := time.Now().Format("0601021504") // Format: YYMMDDHHMM
	randomNumber := rand.Intn(900) + 100         // Random 3 digit (100–999)
	return fmt.Sprintf("%s%d", timestamp, randomNumber)
}
func GetPerizinanSantri(c *gin.Context) {
	db := database.DB

	var data []models.PerizinanSantri

	if err := db.Preload("Santri").Order("id DESC").Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
func AddPerizinanSantri(c *gin.Context) {
	db := database.DB

	// Ambil form-data
	santriID, _ := strconv.Atoi(c.PostForm("santri_id"))
	tanggalIzin := c.PostForm("tanggal_izin")
	tanggalKembali := c.PostForm("tanggal_kembali")
	alasan := c.PostForm("alasan")
	jenisIzin := c.PostForm("jenis_izin")
	statusDokumen := c.PostForm("status_dokumen")

	// Validasi awal bisa ditambah di sini jika diperlukan

	// Buat struct dan isi datanya
	perizinan := models.PerizinanSantri{
		SantriID:       uint(santriID),
		TanggalIzin:    tanggalIzin,
		TanggalKembali: tanggalKembali,
		Alasan:         alasan,
		JenisIzin:      jenisIzin,
		Kode:           generateKodeIzin(), // contoh fungsi kode unik
		StatusAktif:    1,
		StatusDokumen:  statusDokumen, // status awal
	}

	// Simpan ke database
	if err := db.Create(&perizinan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan perizinan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perizinan berhasil disimpan",
		"data":    perizinan,
	})
}
