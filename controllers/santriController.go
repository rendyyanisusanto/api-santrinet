package controllers

import (
	"net/http"
	"santrinet-api/database"
	"santrinet-api/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllSantri(c *gin.Context) {
	db := database.DB

	var santri []models.Santri
	if err := db.Find(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data"})
		return
	}

	c.JSON(http.StatusOK, santri)
}

// GetSantriWithPagination - Endpoint dengan pagination dan filter untuk santri
func GetSantriWithPagination(c *gin.Context) {
	db := database.DB

	// Parameter pagination
	pageParam := c.DefaultQuery("page", "1")
	limitParam := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 10
	}

	// Parameter filter
	search := c.Query("search")              // Pencarian berdasarkan nama atau NIS
	jenisKelamin := c.Query("jenis_kelamin") // Filter jenis kelamin
	statusAktif := c.Query("status_aktif")   // Filter status aktif
	statusSantri := c.Query("status_santri") // Filter status santri

	// Query builder
	query := db.Model(&models.Santri{})

	// Apply filters
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(nama) LIKE ? OR LOWER(nis) LIKE ?", searchPattern, searchPattern)
	}

	if jenisKelamin != "" {
		query = query.Where("jenis_kelamin = ?", jenisKelamin)
	}

	if statusAktif != "" {
		query = query.Where("status_aktif = ?", statusAktif)
	}

	if statusSantri != "" {
		query = query.Where("status_santri = ?", statusSantri)
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghitung total data",
		})
		return
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get paginated data
	var santri []models.Santri
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data santri",
		})
		return
	}

	// Calculate pagination info
	totalPages := (int(total) + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	// Response
	response := gin.H{
		"data": santri,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total":        total,
			"total_pages":  totalPages,
			"has_next":     hasNext,
			"has_previous": hasPrev,
		},
		"filters": gin.H{
			"search":        search,
			"jenis_kelamin": jenisKelamin,
			"status_aktif":  statusAktif,
			"status_santri": statusSantri,
		},
	}

	c.JSON(http.StatusOK, response)
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
	// ambil data dokumen santri
	var dokumenSantri []models.SantriDokumen
	if err := db.Where("santri_id = ?", id).Find(&dokumenSantri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil dokumen santri",
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
		"santri":         santri,
		"kamar_santri":   response,
		"dokumen_santri": dokumenSantri,
	})
}

// CreateSantri - Endpoint untuk insert data santri
func CreateSantri(c *gin.Context) {
	db := database.DB

	var santri models.Santri
	if err := c.ShouldBindJSON(&santri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	// Validasi sederhana: NIS dan Nama wajib diisi
	if santri.Nis == "" || santri.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIS dan Nama wajib diisi"})
		return
	}

	// Insert ke database
	if err := db.Create(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Data santri berhasil ditambahkan", "data": santri})
}

// SoftDeleteSantri - Ubah status_aktif santri menjadi nonaktif ("0")
func SoftDeleteSantri(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	var santri models.Santri
	if err := db.First(&santri, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Santri tidak ditemukan"})
		return
	}

	// Update status_aktif
	santri.StatusAktif = "0"
	if err := db.Save(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengubah status santri"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Santri berhasil dinonaktifkan", "data": santri})
}

// UpdateSantri - Update data santri berdasarkan ID
func UpdateSantri(c *gin.Context) {
	db := database.DB
	id := c.Param("id")

	var santri models.Santri
	if err := db.First(&santri, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Santri tidak ditemukan"})
		return
	}

	var input models.Santri
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid", "detail": err.Error()})
		return
	}

	// Update semua field kecuali ID
	santri.Nis = input.Nis
	santri.Nama = input.Nama
	santri.JenisKelamin = input.JenisKelamin
	santri.TempatLahir = input.TempatLahir
	santri.TanggalLahir = input.TanggalLahir
	santri.Alamat = input.Alamat
	santri.NamaAyah = input.NamaAyah
	santri.NamaIbu = input.NamaIbu
	santri.NoHPAyah = input.NoHPAyah
	santri.NoHPIbu = input.NoHPIbu
	santri.StatusAktif = input.StatusAktif
	santri.StatusSantri = input.StatusSantri
	santri.NamaWali = input.NamaWali
	santri.NoHPWali = input.NoHPWali
	santri.Foto = input.Foto
	santri.Nik = input.Nik

	if err := db.Save(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update data santri"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data santri berhasil diupdate", "data": santri})
}

// GetSantriFiltered - Get semua santri dengan filter (tanpa pagination)
func GetSantriFiltered(c *gin.Context) {
	db := database.DB

	search := c.Query("search") // Pencarian nama/NIS
	jenisKelamin := c.Query("jenis_kelamin")
	statusAktif := c.Query("status_aktif")
	statusSantri := c.Query("status_santri")

	query := db.Model(&models.Santri{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("nama LIKE ? OR nis LIKE ?", searchPattern, searchPattern)
	}
	if jenisKelamin != "" {
		query = query.Where("jenis_kelamin = ?", jenisKelamin)
	}
	if statusAktif != "" {
		query = query.Where("status_aktif = ?", statusAktif)
	}
	if statusSantri != "" {
		query = query.Where("status_santri = ?", statusSantri)
	}

	var santri []models.Santri
	if err := query.Order("id DESC").Find(&santri).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data santri"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": santri})
}
