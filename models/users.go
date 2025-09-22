package models

type Users struct {
	ID        uint    `gorm:"column:id;primaryKey" json:"id"`
	Username  string  `gorm:"column:username" json:"username"`
	Password  string  `gorm:"column:password" json:"-"`
	FirstName string  `gorm:"column:first_name" json:"first_name"`
	LastName  string  `gorm:"column:last_name" json:"last_name"`
	Company   string  `gorm:"column:company" json:"company"`
	Phone     string  `gorm:"column:phone" json:"phone"`
	Foto      string  `gorm:"column:foto" json:"foto"`
	AnggotaID uint    `gorm:"column:anggota_id" json:"anggota_id"`
	Table     string  `gorm:"column:table" json:"table"`
	Groups    []Group `gorm:"many2many:users_groups;joinForeignKey:UserID;joinReferences:GroupID" json:"groups"`
}

type Group struct {
	ID           uint           `gorm:"column:id;primaryKey" json:"id"`
	Name         string         `gorm:"column:name" json:"name"`
	GroupsDetail []GroupsDetail `gorm:"foreignKey:GroupsID;references:ID" json:"groups_detail"`
}
