package models

type Tatib struct {
	ID              uint          `gorm:"column:id;primaryKey" json:"id"`
	Kode            string        `gorm:"column:kode" json:"kode"`
	Nama            string        `gorm:"column:nama" json:"nama"`
	Poin            int           `gorm:"column:poin" json:"poin"`
	Sanksi          string        `gorm:"column:sanksi" json:"sanksi"`
	KategoriTatibID uint          `gorm:"column:kategori_tatib_id" json:"kategori_tatib_id"`
	KategoriTatib   KategoriTatib `gorm:"foreignKey:KategoriTatibID;references:ID" json:"kategori_tatib"`
}

func (Tatib) TableName() string {
	return "tatib"
}
