package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// https://developpaper.com/gorm-relation-one-to-one-one-to-many-many-to-many-query/

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey;"`
	Username  string `gorm:"unique"`
	Password  string
	Joined    time.Time `gorm:"index:joined;autoCreateTime"`
	LastLogin time.Time `gorm:"index:last_login;autoUpdateTime"`
	// ChatID    uint
}

type Chat struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;"`
	Outbound User
	Inbound  User
	// Messages []Message `gorm:"many2many:chat_message"`
	Date time.Time `gorm:"index:joined;autoCreateTime"`
}

type Message struct {
	gorm.Model
	ID   uint `gorm:"primaryKey;"`
	Chat Chat `gorm:"foreignKey:ChatRefer"`
	// Outbound string
	// Inbound  string
	Message string    `gorm:"index:message"`
	Date    time.Time `gorm:"index:joined;autoCreateTime"`
	// Chats    []Chat    `gorm:"many2many:chat_message"`
}

type Claims struct {
	jwt.StandardClaims
	ID uint `gorm:"primaryKey"`
}
