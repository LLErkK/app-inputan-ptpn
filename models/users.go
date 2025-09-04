package models

import "time"

type User struct {
	ID        uint `gorm:"primary_key;auto_increment"`
	Username  string
	Password  string
	LastLogin time.Time
}
