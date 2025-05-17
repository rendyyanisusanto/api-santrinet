package models

type Tatib struct {
	ID   uint   `gorm:"column:id;primaryKey" json:"id"`
	Kode string `gorm:"column:kode" json:"kode"`
	Nama string `gorm:"column:nama" json:"nama"`
	Poin int    `gorm:"column:poin" json:"poin"`
}

func (Tatib) TableName() string {
	return "tatib"
}
