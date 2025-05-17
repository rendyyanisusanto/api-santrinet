package models

type PresensiSantri struct {
	ID               uint   `gorm:"column:id;primaryKey" json:"id"`
	Tanggal          string `gorm:"column:tanggal" json:"tanggal"`
	SantriID         uint   `gorm:"column:santri_id" json:"santri_id"`
	StatusKehadiran  string `gorm:"column:status_kehadiran" json:"status_kehadiran"`
	JadwalPresensiID uint   `gorm:"column:jadwal_presensi_id" json:"jadwal_presensi_id"`

	JadwalPresensi JadwalPresensi `gorm:"foreignKey:JadwalPresensiID;references:ID" json:"jadwal_presensi"`
}
type PresensiJoin struct {
	Tanggal         string
	NamaJadwal      string
	StatusKehadiran string
}

func (PresensiSantri) TableName() string {
	return "jadwal_presensi"
}
