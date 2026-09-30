package ds

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Login    string `gorm:"type:varchar(50);unique;not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

func (User) TableName() string {
	return "users"
}