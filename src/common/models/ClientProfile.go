package models

import "time"

type ClientProfile struct {
	BirthOfDate time.Time `json:"birth_of_date"`
	Group       string    `json:"group"`
	UserID      uint      `json:"user_id"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}
