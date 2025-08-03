package models

type IjinPengurus struct {
	ID           uint   `gorm:"column:id;primaryKey" json:"id"`
	SantriID     string `gorm:"column:santri_id" json:"santri_id"`
	WaktuKeluar  string `gorm:"column:waktu_keluar" json:"waktu_keluar"`
	WaktuKembali string `gorm:"column:waktu_kembali" json:"waktu_kembali"`
	Status       string `gorm:"column:status" json:"status"`

	// Relasi ke Santri
	Santri Santri `json:"Santri" gorm:"foreignKey:ID;references:SantriID"`
}

func (IjinPengurus) TableName() string {
	return "ijin_pengurus"
}
