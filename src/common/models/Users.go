package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255)" json:"-"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	Phone    string `gorm:"type:varchar(20)" json:"phone"`
	RoleID   uint   `gorm:"not null" json:"role_id"`
	Role     Role   `gorm:"foreignKey:RoleID" json:"role"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
