package model

import "golang.org/x/crypto/bcrypt"

type User struct {
	UserID   int64  `gorm:"primarykey"`
	UserName string `gorm:"unique"`
	NickName string
	Password string
}

func (*User) TableName() string {
	return "user"
}

func (user *User) SetPwd(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 这里生成的是哈希码（乱码）
	user.Password = string(hash)
	return nil
}

func (user *User) CheckPwd(pwd string) bool {
	// 比较哈希码（乱码）是否和这个密码匹配
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	return err == nil
}
