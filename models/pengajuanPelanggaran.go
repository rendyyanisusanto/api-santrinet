package models

type PengajuanPelanggaran struct {
	ID              uint   `gorm:"column:id;primaryKey" json:"id"`
	Kode            string `json:"kode"`
	TatibID         uint   `json:"tatib_id"`
	PelaporID       uint   `json:"pelapor_id"`
	Tanggal         string `json:"tanggal"`
	StatusAktif     uint   `json:"status_aktif"`
	Kronologi       string `json:"kronologi"`
	Foto            string `json:"foto"`
	StatusPengajuan string `json:"status_pengajuan"`

	// Relasi detail tatib
	Tatib Tatib `json:"Tatib" gorm:"foreignKey:ID;references:TatibID"`
}

func (PengajuanPelanggaran) TableName() string {
	return "pengajuan_pelanggaran"
}

type PengajuanInput struct {
	Kode            string `json:"kode"`
	TatibID         uint   `json:"tatib_id"`
	PelaporID       uint   `json:"pelapor_id"`
	Tanggal         string `json:"tanggal"`
	StatusAktif     uint   `json:"status_aktif"`
	Kronologi       string `json:"kronologi"`
	Foto            string `json:"foto"`
	StatusPengajuan string `json:"status_pengajuan"`
	SantriIDs       []uint `json:"santri_ids"`
}
