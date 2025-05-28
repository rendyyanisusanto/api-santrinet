// models/detail_pengajuan_pelanggaran.go
package models

type DetailPengajuanPelanggaran struct {
	ID                     uint `json:"id" gorm:"primaryKey"`
	PengajuanPelanggaranID uint `json:"pengajuan_pelanggaran_id"`
	SantriID               uint `json:"santri_id"`
}

func (DetailPengajuanPelanggaran) TableName() string {
	return "detail_pengajuan_pelanggaran"
}
