package response

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	RoleID       uint   `json:"role_id"`
	Email        string `json:"email"`
	ConfirmEmail bool   `json:"comfirm_email"`
	InitialLogin bool   `json:"initial_login"`
}
