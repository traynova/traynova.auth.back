package models

import "gorm.io/gorm"

type UserTokenType struct {
	gorm.Model

	ID   uint   `gorm:"primaryKey" json:"id"`
	Type string `gorm:"type:varchar(50);uniqueIndex;not null" json:"type"`
}
