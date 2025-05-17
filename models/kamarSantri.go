package models

type KamarSantri struct {
	ID          uint   `gorm:"column:id;primaryKey" json:"id"`
	SantriID    uint   `gorm:"column:santri_id" json:"santri_id"`
	KamarID     uint   `gorm:"column:kamar_id" json:"kamar_id"`
	Tahun       string `gorm:"column:tahun" json:"tahun"`
	StatusAktif string `gorm:"column:status_aktif" json:"status_aktif"`
	Kamar       Kamar  `gorm:"foreignKey:KamarID;references:ID" json:"kamar"`
}

type KamarSantriResponse struct {
	ID          uint   `json:"id"`
	KamarID     uint   `json:"kamar_id"`
	AsramaID    uint   `json:"asrama_id"`
	NamaKamar   string `json:"nama_kamar"`
	NamaAsrama  string `json:"nama_asrama"`
	StatusAktif string `json:"status_aktif"`
	Tahun       string `json:"tahun"`
}

func (KamarSantri) TableName() string {
	return "kamar_santri"
}
