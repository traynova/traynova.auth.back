package jwt_requests

// CheckOneTimeJwtTokenRequest
// @Description checkea si un token de un solo uso ha sido usado
type CheckOneTimeJwtTokenRequest struct {
	// JWT Token para revisar si ha sido usado
	Token string `json:"token" binding:"required"`
}
