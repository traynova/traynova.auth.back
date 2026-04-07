package structs

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email  string `json:"email" binding:"required,email"`
	Name   string `json:"name" binding:"required"`
	Phone  string `json:"phone" binding:"required"`
	RoleID uint   `json:"role_id" binding:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

