package models

type Perawat struct {
	ID          uint   `gorm:"column:id;primaryKey" json:"id"`
	Nama        string `gorm:"column:nama" json:"nama"`
	Kode        string `gorm:"column:kode" json:"kode"`
	StatusAktif string `gorm:"column:status_aktif" json:"status_aktif"`
}

func (Perawat) TableName() string {
	return "perawat"
}
