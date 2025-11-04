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

// Struct untuk response laporan
type LaporanRekapHarian struct {
	Tanggal       string                     `json:"tanggal"`
	KamarID       *uint                      `json:"kamar_id,omitempty"`
	NamaKamar     string                     `json:"nama_kamar,omitempty"`
	WaktuPresensi map[string]StatusKehadiran `json:"waktu_presensi"`
}

type StatusKehadiran struct {
	Hadir  int `json:"hadir"`
	Alpha  int `json:"alpha"`
	Sakit  int `json:"sakit"`
	Izin   int `json:"izin"`
}

type LaporanBulananPerSantri struct {
	SantriID    uint   `json:"santri_id"`
	NamaSantri  string `json:"nama_santri"`
	Bulan       string `json:"bulan"`
	Hadir       int    `json:"hadir"`
	Alpha       int    `json:"alpha"`
	Sakit       int    `json:"sakit"`
	Izin        int    `json:"izin"`
	TotalHari   int    `json:"total_hari"`
}

type LaporanKehadiranPerWaktu struct {
	SantriID      uint   `json:"santri_id"`
	NamaSantri    string `json:"nama_santri"`
	WaktuPresensi string `json:"waktu_presensi"`
	TotalHadir    int    `json:"total_hadir"`
	TotalAlpha    int    `json:"total_alpha"`
	TotalSakit    int    `json:"total_sakit"`
	TotalIzin     int    `json:"total_izin"`
}

// GetLaporanRekapHarian - Laporan rekap harian mengelompokkan berdasarkan waktu presensi
func GetLaporanRekapHarian(c *gin.Context) {
	db := database.DB
	
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tanggal wajib diisi"})
		return
	}

	kamarIDStr := c.Query("kamar_id") // optional: filter per kamar

	var hasil []struct {
		WaktuPresensi   string `json:"waktu_presensi"`
		StatusKehadiran string `json:"status_kehadiran"`
		Jumlah          int    `json:"jumlah"`
	}

	query := db.Table("presensi_santri").
		Select("jadwal_presensi.waktu_presensi, presensi_santri.status_kehadiran, COUNT(*) as jumlah").
		Joins("JOIN jadwal_presensi ON jadwal_presensi.id = presensi_santri.jadwal_presensi_id").
		Where("presensi_santri.tanggal = ?", tanggal)

	if kamarIDStr != "" {
		// Filter santri berdasarkan kamar aktif saat ini
		query = query.Joins("JOIN kamar_santri ks ON ks.santri_id = presensi_santri.santri_id").
			Where("ks.kamar_id = ? AND ks.status_aktif = 1", kamarIDStr)
	}

	err := query.
		Group("jadwal_presensi.waktu_presensi, presensi_santri.status_kehadiran").
		Scan(&hasil).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan harian"})
		return
	}

	// Kelompokkan data berdasarkan waktu presensi
	rekapData := make(map[string]StatusKehadiran)
	
	for _, row := range hasil {
		if _, exists := rekapData[row.WaktuPresensi]; !exists {
			rekapData[row.WaktuPresensi] = StatusKehadiran{}
		}
		
		status := rekapData[row.WaktuPresensi]
		switch row.StatusKehadiran {
		case "HADIR":
			status.Hadir = row.Jumlah
		case "ALPHA":
			status.Alpha = row.Jumlah
		case "SAKIT":
			status.Sakit = row.Jumlah
		case "IZIN":
			status.Izin = row.Jumlah
		}
		rekapData[row.WaktuPresensi] = status
	}

	var kamarIDPtr *uint
	var namaKamar string
	if kamarIDStr != "" {
		if idParsed, errParse := strconv.Atoi(kamarIDStr); errParse == nil {
			idu := uint(idParsed)
			kamarIDPtr = &idu
			// ambil nama kamar untuk kenyamanan response
			var kamar models.Kamar
			if err := db.Table("kamar").Where("id = ?", idu).First(&kamar).Error; err == nil {
				namaKamar = kamar.Nama
			}
		}
	}

	response := LaporanRekapHarian{
		Tanggal:       tanggal,
		KamarID:       kamarIDPtr,
		NamaKamar:     namaKamar,
		WaktuPresensi: rekapData,
	}

	c.JSON(http.StatusOK, response)
}

// GetLaporanBulananPerSantri - Laporan rekap bulanan per santri
func GetLaporanBulananPerSantri(c *gin.Context) {
	db := database.DB
	
	bulan := c.Query("bulan") // format: 2025-11
	santriID := c.Query("santri_id") // optional, jika kosong tampilkan semua santri
	
	if bulan == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter bulan wajib diisi (format: YYYY-MM)"})
		return
	}

	query := db.Table("presensi_santri").
		Select(`
			presensi_santri.santri_id,
			santri.nama as nama_santri,
			? as bulan,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'HADIR' THEN 1 ELSE 0 END) as hadir,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'ALPHA' THEN 1 ELSE 0 END) as alpha,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'SAKIT' THEN 1 ELSE 0 END) as sakit,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'IZIN' THEN 1 ELSE 0 END) as izin,
			COUNT(DISTINCT presensi_santri.tanggal) as total_hari
		`, bulan).
		Joins("JOIN santri ON santri.id = presensi_santri.santri_id").
		Where("DATE_FORMAT(presensi_santri.tanggal, '%Y-%m') = ?", bulan).
		Group("presensi_santri.santri_id, santri.nama")

	if santriID != "" {
		query = query.Where("presensi_santri.santri_id = ?", santriID)
	}

	var hasil []LaporanBulananPerSantri
	err := query.Scan(&hasil).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan bulanan"})
		return
	}

	c.JSON(http.StatusOK, hasil)
}

// GetLaporanKehadiranPerWaktu - Laporan kehadiran per waktu presensi
func GetLaporanKehadiranPerWaktu(c *gin.Context) {
	db := database.DB
	
	tanggalMulai := c.Query("tanggal_mulai")
	tanggalSelesai := c.Query("tanggal_selesai")
	santriID := c.Query("santri_id") // optional
	waktuPresensi := c.Query("waktu_presensi") // optional: pagi, siang, malam
	
	if tanggalMulai == "" || tanggalSelesai == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter tanggal_mulai dan tanggal_selesai wajib diisi"})
		return
	}

	query := db.Table("presensi_santri").
		Select(`
			presensi_santri.santri_id,
			santri.nama as nama_santri,
			jadwal_presensi.waktu_presensi,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'HADIR' THEN 1 ELSE 0 END) as total_hadir,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'ALPHA' THEN 1 ELSE 0 END) as total_alpha,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'SAKIT' THEN 1 ELSE 0 END) as total_sakit,
			SUM(CASE WHEN presensi_santri.status_kehadiran = 'IZIN' THEN 1 ELSE 0 END) as total_izin
		`).
		Joins("JOIN santri ON santri.id = presensi_santri.santri_id").
		Joins("JOIN jadwal_presensi ON jadwal_presensi.id = presensi_santri.jadwal_presensi_id").
		Where("presensi_santri.tanggal BETWEEN ? AND ?", tanggalMulai, tanggalSelesai).
		Group("presensi_santri.santri_id, santri.nama, jadwal_presensi.waktu_presensi")

	if santriID != "" {
		query = query.Where("presensi_santri.santri_id = ?", santriID)
	}

	if waktuPresensi != "" {
		query = query.Where("jadwal_presensi.waktu_presensi = ?", waktuPresensi)
	}

	var hasil []LaporanKehadiranPerWaktu
	err := query.Scan(&hasil).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan kehadiran per waktu"})
		return
	}

	c.JSON(http.StatusOK, hasil)
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
