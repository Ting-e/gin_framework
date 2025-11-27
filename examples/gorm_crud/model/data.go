package model

import "time"

type User struct {
	User_id     string    `gorm:"type:varchar(50);primaryKey"`
	Name        string    `gorm:"type:varchar(20)"`
	Age         int       `gorm:"type:int(11)"`
	Create_time time.Time `gorm:"type:datetime"`
}
