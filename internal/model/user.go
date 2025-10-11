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
	CreatedBy     string         `gorm:"type:char(36);not null"`
	UpdatedBy     string         `gorm:"type:char(36);not null"`
	DeletedBy     *string        `gorm:"type:char(36)"`
	CreatedAt     time.Time      `gorm:"not null"`
	UpdatedAt     time.Time      `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type RefreshToken struct {
	ID        string         `gorm:"type:char(36);primary_key"`
	UserID    string         `gorm:"type:char(36);not null"`
	TokenHash string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedBy string         `gorm:"type:char(36);not null"`
	UpdatedBy string         `gorm:"type:char(36);not null"`
	DeletedBy *string        `gorm:"type:char(36)"`
	ExpiresAt time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignKey:UserID"`
}

type PasswordResetToken struct {
	ID        string         `gorm:"type:char(36);primary_key"`
	UserID    string         `gorm:"type:char(36);not null"`
	TokenHash string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	CreatedBy string         `gorm:"type:char(36);not null"`
	UpdatedBy string         `gorm:"type:char(36);not null"`
	DeletedBy *string        `gorm:"type:char(36)"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"foreignKey:UserID"`
}
