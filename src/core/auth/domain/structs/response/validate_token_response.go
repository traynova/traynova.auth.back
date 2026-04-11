package structs_response

type ValidateTokenResponse struct {
	Valid         bool   `json:"valid"`
	UserID        uint   `json:"user_id"`
	RoleID        uint   `json:"role_id"`
	AccessLevelID uint   `json:"access_level_id"`
	Email         string `json:"email"`
	ExpiresAt     int64  `json:"expires_at"`
}
