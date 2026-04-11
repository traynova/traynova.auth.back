package structs_response

type GetUserResponse struct {
	ID             uint   `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	Prefix         string `json:"prefix"`
	RoleID         uint   `json:"role_id"`
	RoleName       string `json:"role_name"`
	IsActive       bool   `json:"is_active"`
	EmailConfirmed bool   `json:"email_confirmed"`
}
