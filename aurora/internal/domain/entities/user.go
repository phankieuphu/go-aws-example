package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string    `gorm:"type:varchar(255)" json:"name"`
	Email    string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	CreateAt time.Time `gorm:"autoCreateTime" json:"create_at"`
}

func (User) TableName() string {
	return "users"
}
