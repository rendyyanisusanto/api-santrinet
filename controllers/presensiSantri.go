package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PresensiSantriResponse struct {
	SantriID        uint   `json:"santri_id"`
	NamaSantri      string `json:"nama_santri"`
	StatusKehadiran string `json:"status_kehadiran,omitempty"` // kosong jika belum absen
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
		response = append(response, PresensiSantriResponse{
			SantriID:        ks.SantriID,
			NamaSantri:      ks.Nama,
			StatusKehadiran: presensiMap[ks.SantriID], // otomatis kosong kalau belum ada
		})
	}

	c.JSON(http.StatusOK, response)
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
