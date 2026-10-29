package entities

type Role struct {
	Id   string
	Name string
}

type UserStore struct {
	Id      string `json:"id" gorm:"type:varchar(36);primaryKey"`
	StoreId string `json:"store_id" gorm:"column:store_id;type:varchar(36);not null"`
	UserId  string `json:"user_id" gorm:"column:user_id;type:varchar(36);not null"`
	RoleId  string `json:"role_id" gorm:"column:role_id;type:varchar(50);not null"`
	Role    Role   `gorm:"foreignKey:RoleId;references:Id"`
}
