package model

type Task struct {
	TaskId    int64 `gorm:"primarykey"`
	UserId    int64
	Status    int `gorm:"default:0"`
	Title     string
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
}
