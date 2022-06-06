package models

import (
	"time"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint64 	`gorm:"primaryKey;"`
	Username  string    `gorm:"index:username,unique"`
	Password  string    `gorm:"index:password"`
	Joined    time.Time `gorm:"index:joined;autoCreateTime"`
	LastLogin time.Time `gorm:"index:last_login;autoUpdateTime"`
}

type Chat struct {
	gorm.Model
	ID          uint64 	  `gorm:"primaryKey;"`
	Messages 	[]Message `gorm:"foreignKey:ID;references:ID"` 
	Date		time.Time `gorm:"index:joined;autoCreateTime"`
}

type Message struct {
	gorm.Model
	ID        	uint64 	`gorm:"primaryKey;"`
	From 		User
	To   		User
	Message     string
	Date		time.Time `gorm:"index:joined;autoCreateTime"`
}
