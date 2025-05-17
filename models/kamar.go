package models

type Kamar struct {
	ID       uint   `gorm:"column:id;primaryKey" json:"id"`
	Nama     string `gorm:"column:nama" json:"nama"`
	AsramaID uint   `gorm:"column:asrama_id" json:"asrama_id"`
	Asrama   Asrama `gorm:"foreignKey:AsramaID;references:ID" json:"asrama"`
}

func (Kamar) TableName() string {
	return "kamar"
}
