package structs

type CreateUserRequest struct {
	Email  string `json:"email" binding:"required,email"`
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	RoleID uint   `json:"role_id" binding:"required"`
}
