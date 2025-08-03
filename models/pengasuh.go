package models

type Pengasuh struct {
	ID   uint   `gorm:"column:id;primaryKey" json:"id"`
	Nama string `gorm:"column:nama" json:"nama"`
	Foto string `gorm:"column:foto" json:"foto"`
}

func (Pengasuh) TableName() string {
	return "pengasuh"
}
