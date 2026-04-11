package jwt_service

import (
	"errors"
	"fmt"
	"gestrym/src/common/models"
	"gestrym/src/common/utils"
	jwt_ports "gestrym/src/core/jwt/domain/ports"
	jwt_structs "gestrym/src/core/jwt/domain/structs"
	jwt_requests "gestrym/src/core/jwt/domain/structs/request"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

var (
	ErrInvalidToken = errors.New("token inválido")
	ErrMissingKey   = errors.New("la variable de entorno JWT_KEY no está configurada")
	ErrExpiration   = errors.New("la variable de entorno JWT_EXPIRATION no está configurada")
)

type IJWTService interface {
	GenerateJwtToken(request jwt_requests.GenerateJwtTokenRequest, duration *int) (string, error)
	RegisterToken(token models.UserToken) error
	ValidateJwtToken(tokenString string) (*jwt_structs.CustomClaims, error)
	ChecUserTokenUsed(token string) error
	DeleteUserToken(token string) error
}

type jwtService struct {
	logger                 utils.ILogger
	expiration             time.Duration
	jwtkey                 []byte
	refreshTokenRepository jwt_ports.IRefreshTokenRepository
	userTokenRepository    jwt_ports.IUserTokenRepository
}

var (
	jwtServiceInstance *jwtService
	jwtServiceOnce     sync.Once
)

func NewJWTService(refreshRepo jwt_ports.IRefreshTokenRepository, userTokenRepo jwt_ports.IUserTokenRepository) (IJWTService, error) {
	jwtServiceOnce.Do(func() {
		logger := utils.NewLogger()
		jwtKey := []byte(viper.GetString("JWT_KEY"))
		if len(jwtKey) == 0 {
			logger.Fatal("la variable de entorno JWT_KEY no está configurada")
		}

		expirationMinutes := viper.GetInt("JWT_EXPIRATION")
		if expirationMinutes == 0 {
			logger.Fatal("la variable de entorno JWT_EXPIRATION no está configurada")
		}

		jwtServiceInstance = &jwtService{
			logger:                 logger,
			expiration:             time.Duration(expirationMinutes) * time.Minute,
			refreshTokenRepository: refreshRepo,
			userTokenRepository:    userTokenRepo,
		}
	})
	return jwtServiceInstance, nil
}

func (j *jwtService) GenerateJwtToken(request jwt_requests.GenerateJwtTokenRequest, duration *int) (string, error) {
	if err := request.Validate(); err != nil {
		return "", fmt.Errorf("solicitud de generación de token inválida: %w", err)
	}
	var expiration time.Duration
	if duration != nil {
		expiration = time.Duration(*duration) * time.Minute
	} else {
		expiration = j.expiration
	}
	now := time.Now()
	claims := &jwt_structs.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   request.Email,                           // Identificador único del usuario.
			IssuedAt:  jwt.NewNumericDate(now),                 // Indica cuándo se emitió el token.
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)), // Indica cuándo expira el token.
			NotBefore: jwt.NewNumericDate(now),                 // Indica cuándo el token es válido.
		},
		UserID:        request.UserID,
		RoleID:        request.RoleID,
		AccessLevelID: request.AccessLevelID,
		PhoneNumber:   request.PhoneNumber,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.jwtkey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return signedToken, nil
}

func (j *jwtService) Generate24HJwtToken(request jwt_requests.GenerateJwtTokenRequest) (string, error) {
	if err := request.Validate(); err != nil {
		return "", fmt.Errorf("solicitud de generación de token inválida: %w", err)
	}

	now := time.Now()
	claims := &jwt_structs.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   request.Email,                               // Identificador único del usuario.
			IssuedAt:  jwt.NewNumericDate(now),                     // Indica cuándo se emitió el token.
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)), // Indica cuándo expira el token.
			NotBefore: jwt.NewNumericDate(now),                     // Indica cuándo el token es válido.
		},
		UserID:        request.UserID,
		RoleID:        request.RoleID,
		AccessLevelID: request.AccessLevelID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.jwtkey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return signedToken, nil
}

func (j *jwtService) ValidateJwtToken(tokenString string) (*jwt_structs.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt_structs.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return j.jwtkey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("el token ha expirado: %w", err)
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("el token aún no es válido: %w", err)
			}
		}
		return nil, fmt.Errorf("error al validar el token: %w", err)
	}

	claims, ok := token.Claims.(*jwt_structs.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("no se pudo obtener los claims del token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("el token no es válido")
	}

	return claims, nil
}

// Refresh Token --------------------------------------------

func (j *jwtService) GenerateRefreshToken(userId uint) (string, error) {
	var key string
	var errCreateRefreshToken error
	var errUpdateRefreshToken error

	errValRefreshToken := j.ValidateRefreshTokenByUserId(userId)
	if errValRefreshToken != nil {
		if strings.Contains(errValRefreshToken.Error(), "record not found") {
			key, errCreateRefreshToken = j.CreateRefreshToken(userId)
			if errCreateRefreshToken != nil {
				return "", fmt.Errorf("error al crear el refresh token: %w", errCreateRefreshToken)
			}
		} else {
			// Si hay un error diferente al no encontrar el token
			return "", fmt.Errorf("error al validar el token de actualización: %w", errValRefreshToken)
		}
	} else {
		key, errUpdateRefreshToken = j.UpdateRefreshToken(userId)
		if errUpdateRefreshToken != nil {
			return "", fmt.Errorf("error al actualizar el refresh token: %w", errUpdateRefreshToken)
		}
	}
	return key, nil
}

func (j *jwtService) ValidateRefreshTokenByUserId(userId uint) error {
	_, err := j.refreshTokenRepository.GetRefreshTokenByUserId(userId)
	if err != nil {
		return err
	}

	return nil
}

func (j *jwtService) CreateRefreshToken(userId uint) (string, error) {

	refreshToken := models.RefreshToken{
		Key:        utils.GenerateUuid(),
		UserID:     userId,
		ExpiryDate: time.Now().Add(j.expiration),
	}

	err := j.refreshTokenRepository.InsertRefreshToken(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("error al insertar el token de actualización: %w", err)
	}

	return refreshToken.Key, nil
}

func (j *jwtService) UpdateRefreshToken(userId uint) (string, error) {
	refreshToken, err := j.refreshTokenRepository.GetRefreshTokenByUserId(userId)
	if err != nil {
		return "", fmt.Errorf("error al obtener el token de actualización: %w", err)
	}

	refreshToken.Key = utils.GenerateUuid()
	refreshToken.ExpiryDate = time.Now().Add(time.Hour * 24 * 7)

	err = j.refreshTokenRepository.UpdateRefreshToken(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("error al actualizar el token de actualización: %w", err)
	}

	return refreshToken.Key, nil
}

func (j *jwtService) GetRefreshTokenWithUserAndRole(key string) (models.RefreshToken, error) {
	refreshToken, err := j.refreshTokenRepository.GetRefreshTokenWithUserAndRole(key)
	if err != nil {
		return refreshToken, fmt.Errorf("error al obtener el token de actualización: %w", err)
	}

	return refreshToken, nil
}

// User Token --------------------------------------------
func (j *jwtService) RegisterToken(token models.UserToken) error {
	err := j.userTokenRepository.InsertUserToken(token)
	if err != nil {
		return fmt.Errorf("error al insertar el token de usuario: %w", err)
	}

	return nil
}

func (j *jwtService) DeleteUserToken(token string) error {
	err := j.userTokenRepository.DeleteUserToken(token)
	if err != nil {
		return fmt.Errorf("error al eliminar el token de usuario: %w", err)
	}

	return nil
}

func (j *jwtService) ChecUserTokenUsed(token string) error {
	_, err := j.userTokenRepository.GetUserToken(token)
	if err != nil {
		return err
	}

	return nil
}

func (j *jwtService) CheckLastActivationUserTokenExpired(userId uint) (bool, error) {
	activationToken, err := j.userTokenRepository.GetLastActivationUserToken(userId)
	if err != nil {
		return true, fmt.Errorf("error al obtener el token de activacion: %w", err)
	}
	token, err := jwt.ParseWithClaims(activationToken, &jwt_structs.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar el método de firma
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return j.jwtkey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return true, nil
			}
		}
		return true, fmt.Errorf("error al validar el token: %w", err)
	}

	tokenClaims, ok := token.Claims.(*jwt_structs.CustomClaims)
	if !ok {
		return true, fmt.Errorf("no se pudo obtener los claims del token")
	}

	if tokenClaims.ExpiresAt.Time.Before(time.Now()) {
		return true, nil
	}
	j.logger.Success(fmt.Sprintf("%s >= %s", tokenClaims.ExpiresAt, time.Now()))
	return false, nil
}
