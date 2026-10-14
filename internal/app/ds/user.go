package ds

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Login    string `gorm:"type:varchar(50);unique;not null" json:"login"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type UserAuthRequest struct {
	Login    string `json:"login" form:"login" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}