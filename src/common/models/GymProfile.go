package models

type GymProfile struct {
	City           string  `gorm:"type:varchar(100);not null" json:"city"`
	Department     string  `gorm:"type:varchar(100);not null" json:"department"`
	Country        string  `gorm:"type:varchar(100);not null" json:"country"`
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	Workstation    *string `json:"workstation"`

	ReferredID *uint `json:"referred_id"`
	UserID     uint  `gorm:"uniqueIndex" json:"user_id"`
	FilesID    *uint `json:"files_id"`

	Referral *Referral `gorm:"foreignKey:ReferredID" json:"referral"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
	Files    *Files    `gorm:"foreignKey:FilesID" json:"files"`
}
