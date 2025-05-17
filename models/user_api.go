package models

type User_api struct {
	ID       int    `gorm:"column:id;primaryKey" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

func (User_api) TableName() string {
	return "user_api"
}
