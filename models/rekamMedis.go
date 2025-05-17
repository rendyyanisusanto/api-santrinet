package models

type RekamMedis struct {
	ID                 uint   `json:"id"`
	UUID               string `json:"uuid"`
	SantriID           uint   `json:"santri_id"`
	StatusRekamMedisID uint   `json:"status_rekam_medis_id"`
	Tanggal            string `json:"tanggal"`
	Foto               string `json:"foto"`
	Kode               string `json:"kode"`
	StatusAktif        bool   `json:"status_aktif"`
	PerawatID          uint   `json:"perawat_id"`
	Diagnosis          string `json:"diagnosis"`
	Catatan            string `json:"catatan"`
	LamaSakit          string `json:"lama_sakit"`

	StatusRekamMedis StatusRekamMedis `gorm:"foreignKey:StatusRekamMedisID;references:ID" json:"status_rekam_medis"`
	Perawat          Perawat          `gorm:"foreignKey:PerawatID;references:ID" json:"perawat"`
	ObatList         []Obat           `json:"obat_list" gorm:"many2many:obat_rm;joinForeignKey:RekamMedisID;JoinReferences:ObatID"`
	KeluhanList      []Keluhan        `json:"keluhan_list" gorm:"many2many:keluhan_rm;joinForeignKey:RekamMedisID;JoinReferences:KeluhanID"`
}
