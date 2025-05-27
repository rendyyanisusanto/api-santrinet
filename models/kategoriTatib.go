package models

type KategoriTatib struct {
	ID   uint   `gorm:"column:id;primaryKey" json:"id"`
	Kode string `gorm:"column:kode" json:"kode"`
	Nama string `gorm:"column:nama" json:"nama"`
}

func (KategoriTatib) TableName() string {
	return "kategori_tatib"
}
