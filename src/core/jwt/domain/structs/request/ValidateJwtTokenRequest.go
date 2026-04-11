package jwt_requests

type ValidateJwtTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
