package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"uniqueIndex; not null" json:"username"`
	Password      string `gorm:"not null" json:"password"`
	RefreshTokens string `json:"refresh_tokens"`
}
