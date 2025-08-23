package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PresensiSantriResponse struct {
	SantriID        uint                    `json:"santri_id"`
	NamaSantri      string                  `json:"nama_santri"`
	StatusKehadiran string                  `json:"status_kehadiran,omitempty"` // kosong jika belum absen
	Perizinan       *models.PerizinanSantri `json:"perizinan,omitempty"`
}
type UpdatePresensiRequest struct {
	SantriID         uint   `json:"santri_id" binding:"required"`
	Tanggal          string `json:"tanggal" binding:"required"`
	JadwalPresensiID uint   `json:"jadwal_presensi_id" binding:"required"`
	StatusKehadiran  string `json:"status_kehadiran" binding:"required"`
}

func GetPresensiSantriBySantriID(c *gin.Context) {
	db := database.DB

	id := c.Param("id")

	var hasil []models.PresensiJoin

	err := db.
		Table("presensi_santri").
		Select("presensi_santri.tanggal, jadwal_presensi.waktu_presensi AS nama_jadwal, presensi_santri.status_kehadiran").
		Joins("JOIN jadwal_presensi ON jadwal_presensi.id = presensi_santri.jadwal_presensi_id").
		Where("presensi_santri.santri_id = ?", id).
		Order("presensi_santri.tanggal DESC").
		Scan(&hasil).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data presensi"})
		return
	}

	grouped := make(map[string]map[string]string)

	for _, row := range hasil {
		if _, exists := grouped[row.Tanggal]; !exists {
			grouped[row.Tanggal] = make(map[string]string)
		}
		grouped[row.Tanggal][row.NamaJadwal] = row.StatusKehadiran
	}

	var response []map[string]interface{}

	for tanggal, presensiMap := range grouped {
		item := map[string]interface{}{"tanggal": tanggal}
		for namaJadwal, status := range presensiMap {
			item[namaJadwal] = status
		}
		response = append(response, item)
	}
	c.JSON(http.StatusOK, response)
}

func GetPresensiByKamar(c *gin.Context) {
	db := database.DB

	// Ambil query params
	kamarIDStr := c.Query("kamar_id")
	tanggal := c.Query("tanggal")
	jadwalIDStr := c.Query("jadwal_presensi_id")

	if kamarIDStr == "" || tanggal == "" || jadwalIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "kamar_id, tanggal dan jadwal_presensi_id wajib"})
		return
	}

	kamarID, _ := strconv.Atoi(kamarIDStr)
	jadwalID, _ := strconv.Atoi(jadwalIDStr)

	// Ambil semua santri di kamar yang aktif
	var kamarSantri []struct {
		SantriID uint
		Nama     string
	}

	err := db.Table("kamar_santri").
		Select("santri.id as santri_id, santri.nama").
		Joins("JOIN santri ON santri.id = kamar_santri.santri_id").
		Where("kamar_santri.kamar_id = ? AND kamar_santri.status_aktif = 1", kamarID).
		Scan(&kamarSantri).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil santri"})
		return
	}

	// Ambil presensi santri di tanggal & jadwal itu
	var presensi []models.PresensiSantri
	err = db.Where("tanggal = ? AND jadwal_presensi_id = ? AND santri_id IN (?)",
		tanggal, jadwalID, getSantriIDs(kamarSantri)).
		Find(&presensi).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data presensi"})
		return
	}

	// Map presensi untuk lookup cepat
	presensiMap := make(map[uint]string)
	for _, p := range presensi {
		presensiMap[p.SantriID] = p.StatusKehadiran
	}

	// Bangun response
	var response []PresensiSantriResponse
	for _, ks := range kamarSantri {
		// cek apakah ada perizinan aktif
		var perizinan models.PerizinanSantri
		err := db.Where("? BETWEEN tanggal_izin AND tanggal_kembali AND santri_id = ? AND status_aktif = 1",
			tanggal, ks.SantriID).
			First(&perizinan).Error

		var izin *models.PerizinanSantri
		if err == nil { // ada izin
			izin = &perizinan
		} else {
			izin = nil
		}

		response = append(response, PresensiSantriResponse{
			SantriID:        ks.SantriID,
			NamaSantri:      ks.Nama,
			StatusKehadiran: presensiMap[ks.SantriID], // otomatis kosong kalau belum ada
			Perizinan:       izin,
		})
	}

	c.JSON(http.StatusOK, response)
}

func UpdatePresensiSantri(c *gin.Context) {
	var req UpdatePresensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db := database.DB

	var presensi models.PresensiSantri
	// cek apakah sudah ada datanya
	err := db.Where("santri_id = ? AND tanggal = ? AND jadwal_presensi_id = ?",
		req.SantriID, req.Tanggal, req.JadwalPresensiID).
		First(&presensi).Error

	if err != nil {
		// kalau tidak ada → buat baru
		newPresensi := models.PresensiSantri{
			SantriID:         req.SantriID,
			Tanggal:          req.Tanggal,
			JadwalPresensiID: req.JadwalPresensiID,
			StatusKehadiran:  req.StatusKehadiran,
		}
		if err := db.Create(&newPresensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Gagal membuat presensi",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Presensi berhasil ditambahkan",
			"data":    newPresensi,
		})
		return
	}

	// kalau sudah ada → update status
	presensi.StatusKehadiran = req.StatusKehadiran
	if err := db.Save(&presensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal update presensi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Presensi berhasil diupdate",
		"data":    presensi,
	})
}

// helper ambil slice santri_id
func getSantriIDs(list []struct {
	SantriID uint
	Nama     string
}) []uint {
	var ids []uint
	for _, s := range list {
		ids = append(ids, s.SantriID)
	}
	return ids
}
