package models

type LembagaPengurus struct {
	ID   uint   `gorm:"column:id;primaryKey" json:"id"`
	Kode string `gorm:"column:kode" json:"kode"`
	Nama string `gorm:"column:nama" json:"nama"`
}

func (LembagaPengurus) TableName() string {
	return "lembaga_pengurus"
}
