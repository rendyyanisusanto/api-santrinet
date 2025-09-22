package models

type Submenu struct {
	ID      uint   `gorm:"column:id;primaryKey" json:"id"`
	MenuID  uint   `gorm:"column:menu_id" json:"menu_id"`
	Submenu string `gorm:"column:submenu" json:"submenu"`
	Menu    Menu   `gorm:"foreignKey:MenuID;references:ID" json:"menu"`
}

func (Submenu) TableName() string {
	return "submenu"
}
