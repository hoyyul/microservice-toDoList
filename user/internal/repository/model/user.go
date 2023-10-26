package model

type User struct {
	UserID   int64  `gorm:"primarykey"`
	UserName string `gorm:"unique"`
	NickName string
	Password string
}

func (*User) TableName() string {
	return "user"
}
