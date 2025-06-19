package models

type PerizinanSantri struct {
	ID             uint   `gorm:"column:id;primaryKey" json:"id"`
	Kode           string `json:"kode"`
	SantriID       uint   `json:"santri_id"`
	TanggalIzin    string `json:"tanggal_izin"`
	TanggalKembali string `json:"tanggal_kembali"`
	StatusAktif    uint   `json:"status_aktif"`
	Alasan         string `json:"alasan"`
	JenisIzin      string `json:"jenis_izin"`
	StatusDokumen  string `json:"status_dokumen"`

	// Relasi detail tatib
	Santri Santri `json:"Santri" gorm:"foreignKey:ID;references:SantriID"`
}

func (PerizinanSantri) TableName() string {
	return "perizinan_santri"
}
