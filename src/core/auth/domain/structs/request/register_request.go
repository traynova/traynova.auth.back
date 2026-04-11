package structs_request

type RegisterRequest struct {
	Email        string  `json:"email" binding:"required,email"`
	FullName     string  `json:"name" binding:"required"`
	Username     string  `json:"username" binding:"required"`
	Prefix       string  `json:"prefix" binding:"required"`
	Phone        string  `json:"phone" binding:"required"`
	RoleID       uint    `json:"role_id" binding:"required"`
	ReferralCode *string `json:"referral_code"`
	Password     string  `json:"password" binding:"required"`
	City         *string `json:"city"`
	Department   *string `json:"department"`
	Country      *string `json:"country"`
	Workstation  *string `json:"workstation"`
}
