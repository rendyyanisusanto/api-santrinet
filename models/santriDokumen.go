package models

type SantriDokumen struct {
	ID       uint   `gorm:"column:id;primaryKey" json:"id"`
	SantriID string `gorm:"column:santri_id" json:"santri_id"`
	FName    string `gorm:"column:fname" json:"fname"`
	File     string `gorm:"column:file" json:"file"`

	// Relasi ke Santri
	// Santri Santri `json:"Santri" gorm:"foreignKey:ID;references:SantriID"`
}

func (SantriDokumen) TableName() string {
	return "santri_dokumen"
}
