package models

import "gorm.io/gorm"

type TrainerProfile struct {
	gorm.Model

	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`

	ReferredID *uint `json:"referred_id"`
	UserID     uint  `json:"user_id"`
	FilesID    *uint `json:"files_id"`
	GimID      *uint `json:"gim_id"`

	Referral       *Referral       `gorm:"foreignKey:ReferredID" json:"referral"`
	User           User            `gorm:"foreignKey:UserID" json:"user"`
	Files          *Files          `gorm:"foreignKey:FilesID" json:"files"`
	GymProfile     *GymProfile     `gorm:"foreignKey:GimID" json:"gym_profile"`
	TrainerClients []TrainerClient `gorm:"foreignKey:TrainerProfileID" json:"trainer_clients"`
}
