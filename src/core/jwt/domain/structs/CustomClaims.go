package jwt_structs

import "github.com/golang-jwt/jwt/v4"

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID        uint   `json:"user_id"`
	RoleID        uint   `json:"role_id"`
	AccessLevelID uint   `json:"access_level_id"`
	PhoneNumber   string `json:"phone_number"`
}
