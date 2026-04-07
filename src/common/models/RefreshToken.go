package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model

	ID          uint      `gorm:"primaryKey" json:"id"`
	UserTokenID uint      `gorm:"not null;index" json:"user_token_id"`
	UserToken   UserToken `gorm:"foreignKey:UserTokenID" json:"-"`
	Token       string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	IsRevoked   bool      `gorm:"default:false" json:"is_revoked"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
}
