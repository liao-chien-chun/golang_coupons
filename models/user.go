package models

type User struct {
	ID   uint   `json:"id" gorm:"primaryKey;comment:使用者ID"`
	Name string `json:"name" gorm:"type:varchar(100);not null;comment:使用者名稱"`
}
