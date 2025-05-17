package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"

	"github.com/gin-gonic/gin"
)

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
