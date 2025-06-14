package models

type Pelanggaran struct {
	ID                       uint   `gorm:"column:id;primaryKey" json:"id"`
	Nama                     string `gorm:"column:nama" json:"nama"`
	Kode                     string `gorm:"column:kode" json:"kode"`
	SantriID                 uint   `gorm:"column:santri_id" json:"santri_id"`
	TatibID                  uint   `gorm:"column:tatib_id" json:"tatib_id"`
	PelaporID                uint   `gorm:"column:pelapor_id" json:"pelapor_id"`
	PengurusID               uint   `gorm:"column:pengurus_id" json:"pengurus_id"`
	StatusDokumenPelanggaran string `gorm:"column:status_dokumen_pelanggaran" json:"status_dokumen_pelanggaran"`
	Tanggal                  string `gorm:"column:tanggal" json:"tanggal"`
	StatusAktif              int    `gorm:"column:status_aktif" json:"status_aktif"`
	Kronologi                string `gorm:"column:kronologi" json:"kronologi"`
	Foto                     string `gorm:"column:foto" json:"foto"`
	StatusTakzir             string `gorm:"column:status_takzir" json:"status_takzir"`
	StatusPengajuan          string `gorm:"column:status_pengajuan" json:"status_pengajuan"`

	Santri   Santri   `gorm:"foreignKey:SantriID;references:ID" json:"santri"`
	Tatib    Tatib    `gorm:"foreignKey:ID;references:TatibID" json:"tatib"`
	Pengurus Pengurus `gorm:"foreignKey:ID;references:PengurusID" json:"pengurus"`
	Pelapor  Pengurus `gorm:"foreignKey:ID;references:PelaporID" json:"pelapor"`
}

func (Pelanggaran) TableName() string {
	return "pelanggaran"
}
