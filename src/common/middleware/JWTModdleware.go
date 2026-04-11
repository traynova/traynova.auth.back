package middleware

import (
	"errors"
	"gestrym/src/common/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID        uint `json:"user_id"`
	RoleID        uint `json:"role_id"`
	AccessLevelID uint `json:"access_level_id"`
}

var (
	logger               = utils.NewLogger()
	ErrMissingAuthHeader = errors.New("falta el encabezado de autorización")
	ErrInvalidToken      = errors.New("token no válido")
	ErrInvalidSignature  = errors.New("método de firma no válido")
	ErrTokenExpired      = errors.New("token expirado")
)

func ValidateTokenMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{"mensaje": "Preflight request"})
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			logger.Error("[JWT_AUTHENTICATION_FAILED] %v", ErrMissingAuthHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuthHeader.Error()})
			return
		}

		tokenStr := authHeader
		if strings.HasPrefix(strings.ToLower(authHeader), "Bearer") {
			tokenStr = strings.TrimSpace(authHeader[7:])
		}

		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Error("[JWT_AUTHENTICATION_FAILED] %v", ErrInvalidSignature)
				return nil, ErrInvalidSignature
			}
			return jwtKey, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.Error("[JWT_AUTHENTICATION_FAILED] %v, error: %v", ErrTokenExpired, err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrTokenExpired.Error()})
			} else {
				logger.Error("[JWT_AUTHENTICATION_FAILED] %v, error: %v", ErrInvalidToken, err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			}
			return
		}

		if !token.Valid {
			logger.Error("[JWT_AUTHENTICATION_FAILED] %v", ErrInvalidToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken.Error()})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role_id", claims.RoleID)
		c.Set("access_level_id", claims.AccessLevelID)
		c.Next()
	}
}

func SetupJWTMiddleware() gin.HandlerFunc {
	jwtKey := viper.GetString("JWT_KEY")
	if jwtKey == "" {
		logger.Error("[JWT_MIDDLEWARE_CONFIG_FAILED] JWT_KEY no está configurado en el archivo de configuración")
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "JWT_KEY no está configurado"})
		}
	}
	return ValidateTokenMiddleware([]byte(jwtKey))
}
