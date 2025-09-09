package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:100;not null;unique" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"password"`
	LastLogin time.Time `json:"last_login"`
}
