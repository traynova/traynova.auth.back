package models

type TrainerClient struct {
	TrainerProfileID *uint `json:"trainer_profile_id"`
	ClientID         *uint `json:"client_id"`

	TrainerProfile *TrainerProfile `gorm:"foreignKey:TrainerProfileID" json:"trainer_profile"`
	Client         *User           `gorm:"foreignKey:ClientID" json:"client"`
}
