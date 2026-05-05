package models

import "gorm.io/gorm"

type UserGroup struct {
	gorm.Model

	Name             string `gorm:"type:varchar(100);not null" json:"name"`
	GymProfileID     *uint  `json:"gym_profile_id"`
	TrainerProfileID *uint  `json:"trainer_profile_id"`

	GymProfile     *GymProfile       `gorm:"foreignKey:GymProfileID" json:"gym_profile"`
	TrainerProfile *TrainerProfile   `gorm:"foreignKey:TrainerProfileID" json:"trainer_profile"`
	Members        []UserGroupMember `gorm:"foreignKey:UserGroupID" json:"members"`
}
