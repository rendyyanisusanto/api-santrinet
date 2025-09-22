package models

type Menu struct {
	ID       uint      `gorm:"column:id;primaryKey" json:"id"`
	Menu     string    `gorm:"column:menu" json:"menu"`
	Submenus []Submenu `gorm:"foreignKey:MenuID;references:ID" json:"submenus"`
}

func (Menu) TableName() string {
	return "menu"
}
