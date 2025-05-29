// controllers/pelanggaran_controller.go
package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"santrinet-api/database"
	"santrinet-api/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func generateKodePengajuan() string {
	return fmt.Sprintf("PEL-%s", uuid.New().String()[:8]) // Contoh: PEL-a1b2c3d4
}
func GetPengajuanPelanggaran(c *gin.Context) {
	db := database.DB

	var data []models.PengajuanPelanggaran

	if err := db.Preload("DetailSantri").Order("id DESC").Find(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, data)
}
func AddPengajuanPelanggaran(c *gin.Context) {
	db := database.DB

	file, err := c.FormFile("foto")
	var filePath string
	var relPath string
	if err == nil {
		uploadDir := "/var/www/simsmk/santrinet-uploads/pelanggaran"
		os.MkdirAll(uploadDir, os.ModePerm)

		// Dapatkan ekstensi file
		fileExt := filepath.Ext(file.Filename)

		// Buat nama file unik (misalnya timestamp + random angka)
		newFileName := fmt.Sprintf("%d_%d%s", time.Now().Unix(), rand.Intn(1000), fileExt)

		// Path lengkap untuk disimpan
		filePath = fmt.Sprintf("%s/%s", uploadDir, newFileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		// Simpan path relatif untuk database
		relPath = fmt.Sprintf("pelanggaran/%s", newFileName)
	}

	// Ambil form-data lainnya
	tatibID, _ := strconv.Atoi(c.PostForm("tatib_id"))
	pelaporID, _ := strconv.Atoi(c.PostForm("pelapor_id"))
	tanggal := c.PostForm("tanggal")
	kronologi := c.PostForm("kronologi")
	santriIDsStr := c.PostForm("santri_ids")

	// Parse santri IDs dari string ke slice
	var santriIDs []uint
	for _, idStr := range strings.Split(santriIDsStr, ",") {
		id, _ := strconv.Atoi(strings.TrimSpace(idStr))
		santriIDs = append(santriIDs, uint(id))
	}

	// Simpan data utama
	pengajuan := models.PengajuanPelanggaran{
		Kode:            generateKodePengajuan(),
		TatibID:         uint(tatibID),
		PelaporID:       uint(pelaporID),
		Tanggal:         tanggal,
		StatusAktif:     1,
		Kronologi:       kronologi,
		Foto:            relPath, // path file
		StatusPengajuan: "BELUM DIPROSES",
	}

	if err := db.Create(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengajuan"})
		return
	}

	// Simpan detail santri
	for _, santriID := range santriIDs {
		detail := models.DetailPengajuanPelanggaran{
			PengajuanPelanggaranID: pengajuan.ID,
			SantriID:               santriID,
		}
		db.Create(&detail)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengajuan berhasil disimpan", "data": pengajuan})
}
