package models

import (
	"gorm.io/gorm"
	"time"
)

type Invite struct {
	gorm.Model
	Email     string    `gorm:"uniqueIndex;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
