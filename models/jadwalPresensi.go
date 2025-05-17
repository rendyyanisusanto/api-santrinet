package models

type JadwalPresensi struct {
	ID            uint   `gorm:"column:id;primaryKey" json:"id"`
	WaktuPresensi string `gorm:"column:waktu_presensi" json:"waktu_presensi"`
}

func (JadwalPresensi) TableName() string {
	return "jadwal_presensi"
}
