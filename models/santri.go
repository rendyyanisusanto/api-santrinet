package models

type Santri struct {
	ID           uint   `gorm:"column:id;primaryKey" json:"id"`
	Nis          string `gorm:"column:nis" json:"nis"`
	Nama         string `gorm:"column:nama" json:"nama"`
	JenisKelamin string `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	TempatLahir  string `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir string `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	Alamat       string `gorm:"column:alamat" json:"alamat"`
	NamaAyah     string `gorm:"column:nama_ayah" json:"nama_ayah"`
	NamaIbu      string `gorm:"column:nama_ibu" json:"nama_ibu"`
	NoHPAyah     string `gorm:"column:no_hp_ayah" json:"no_hp_ayah"`
	NoHPIbu      string `gorm:"column:no_hp_ibu" json:"no_hp_ibu"`
	StatusAktif  string `gorm:"column:status_aktif" json:"status_aktif"`
	StatusSantri string `gorm:"column:status_santri" json:"status_santri"`
	NamaWali     string `gorm:"column:nama_wali" json:"nama_wali"`
	NoHPWali     string `gorm:"column:no_hp_wali" json:"no_hp_wali"`
	Foto         string `gorm:"column:foto" json:"foto"`
	Nik          string `gorm:"column:nik" json:"nik"`
}

func (Santri) TableName() string {
	return "santri"
}
