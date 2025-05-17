package models

type Pengurus struct {
	ID                uint            `gorm:"column:id;primaryKey" json:"id"`
	SantriID          string          `gorm:"column:santri_id" json:"santri_id"`
	LembagaPengurusID uint            `gorm:"column:lembaga_pengurus_id" json:"lembaga_pengurus_id"`
	Santri            Santri          `gorm:"foreignKey:SantriID;references:ID" json:"santri"`
	LembagaPengurus   LembagaPengurus `gorm:"foreignKey:LembagaPengurusID;references:ID" json:"lembaga_pengurus"`
}

func (Pengurus) TableName() string {
	return "pengurus"
}
