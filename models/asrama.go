package models

type Asrama struct {
	ID    uint    `gorm:"column:id;primaryKey" json:"id"`
	Nama  string  `gorm:"column:nama" json:"nama"`
	Kamar []Kamar `gorm:"foreignKey:AsramaID" json:"kamar,omitempty"`
}

func (Asrama) TableName() string {
	return "asrama"
}
