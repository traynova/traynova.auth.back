package models

import "gorm.io/gorm"

type UserGroupMember struct {
	gorm.Model

	UserGroupID uint `gorm:"not null" json:"user_group_id"`
	UserID      uint `gorm:"not null" json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
