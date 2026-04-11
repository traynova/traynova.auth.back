package structs_request

type PasswordRecoveryRequest struct {
	Email string `json:"email" binding:"required,email"`
}
