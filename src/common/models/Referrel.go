package models

type Referral struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `gorm:"type:varchar(50);uniqueIndex" json:"code"`
}
