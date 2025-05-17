package models

type StatusRekamMedis struct {
	ID          uint   `gorm:"column:id;primaryKey" json:"id"`
	Nama        string `gorm:"column:nama" json:"nama"`
	Kode        string `gorm:"column:kode" json:"kode"`
	StatusAktif string `gorm:"column:status_aktif" json:"status_aktif"`
}

func (StatusRekamMedis) TableName() string {
	return "status_rekam_medis"
}
