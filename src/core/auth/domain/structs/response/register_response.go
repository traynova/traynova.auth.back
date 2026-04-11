package structs_response

type RegisterResponse struct {
	Id     uint   `json:"id"`
	Email  string `json:"email" binding:"required,email"`
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	RoleID uint   `json:"role_id" binding:"required"`
	Token  string `json:"token" binding:"required"`
}
