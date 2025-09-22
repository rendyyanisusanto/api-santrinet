package models

type GroupsDetail struct {
	ID        uint    `gorm:"column:id;primaryKey" json:"id"`
	GroupsID  uint    `gorm:"column:groups_id" json:"groups_id"`
	SubmenuID uint    `gorm:"column:submenu_id" json:"submenu_id"`
	C         uint8   `gorm:"column:c" json:"c"` // Create permission (1 = allowed, 0 = not allowed)
	R         uint8   `gorm:"column:r" json:"r"` // Read permission
	U         uint8   `gorm:"column:u" json:"u"` // Update permission
	D         uint8   `gorm:"column:d" json:"d"` // Delete permission
	Group     Group   `gorm:"foreignKey:GroupsID;references:ID" json:"group"`
	Submenu   Submenu `gorm:"foreignKey:SubmenuID;references:ID" json:"submenu"`
}

func (GroupsDetail) TableName() string {
	return "groups_detail"
}
