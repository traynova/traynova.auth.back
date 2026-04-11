package structs_response

type GetAllUsersResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
}
