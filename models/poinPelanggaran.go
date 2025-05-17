package models

// Guru defines the model for the Guru (teacher) table
type PoinPelanggaran struct {
	ID               uint   `gorm:"column:id_poin_pelanggaran;primaryKey" json:"id_poin_pelanggaran"`
	Nama_pelanggaran string `json:"nama_pelanggaran"`
	Poin             string `json:"poin"`
	Kode             string `json:"kode"`
}

func (PoinPelanggaran) TableName() string {
	return "poin_pelanggaran"
}
