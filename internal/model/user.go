package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:char(36);primary_key"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string         `gorm:"type:varchar(255);not null"`
	Name          string         `gorm:"type:varchar(100);not null"`
	EmailVerified bool           `gorm:"default:false;not null"`
	CreatedAt     time.Time      `gorm:"not null"`
	UpdatedAt     time.Time      `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type RefreshToken struct {
	ID        string    `gorm:"type:char(36);primary_key"`
	UserID    string    `gorm:"type:char(36);not null"`
	TokenHash string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignKey:UserID"`
}
