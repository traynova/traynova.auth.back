package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"traynova/src/common/middleware"
	"traynova/src/common/models"
	authPorts "traynova/src/core/auth/domain/ports"
	userPorts "traynova/src/core/users/domain/ports"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

type IAuthService interface {
	RegisterUser(ctx context.Context, email, phone, name string, roleID uint) error
	Login(ctx context.Context, email, password string, jwtKey []byte) (string, string, error)
	GoogleLogin(ctx context.Context, token string, jwtKey []byte) (string, string, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, refreshToken string, jwtKey []byte) (string, string, error)
}

type authService struct {
	userRepo  userPorts.IUserRepository
	tokenRepo authPorts.ITokenRepository
}

func NewAuthService(ur userPorts.IUserRepository, tr authPorts.ITokenRepository) IAuthService {
	return &authService{
		userRepo:  ur,
		tokenRepo: tr,
	}
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendConfirmationEmail(userID uint, email, name, token string) {
	payload := map[string]interface{}{
		"user_id":       userID,
		"email":         email,
		"user_name":     name,
		"confirm_token": token,
		"dashboard_url": "http://localhost:3000",
	}
	jsonPayload, _ := json.Marshal(payload)
	http.Post("http://localhost:8443/traynova-notification/public/send-confirmation", "application/json", bytes.NewBuffer(jsonPayload))
}

func (s *authService) RegisterUser(ctx context.Context, email, phone, name string, roleID uint) error {
	user := &models.User{
		Email:    email,
		Phone:    phone,
		Password: "",
		Name:     name,
		RoleID:   roleID,
		IsActive: true,
	}
	err := s.userRepo.Create(ctx, user)
	if err == nil {
		token := generateToken()
		
		ut := &models.UserToken{
			UserID:          user.ID,
			UserTokenTypeID: 2, 
			Token:           token,
			ExpiresAt:       time.Now().Add(24 * time.Hour), 
		}
		s.tokenRepo.CreateUserToken(ctx, ut)

		go sendConfirmationEmail(user.ID, email, name, token)
	}
	return err
}

func (s *authService) Login(ctx context.Context, email, password string, jwtKey []byte) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("credenciales inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("credenciales inválidas")
	}

	if !user.IsActive {
		return "", "", errors.New("usuario inactivo")
	}

	return s.generateAndSaveTokens(ctx, user, jwtKey)
}

func (s *authService) GoogleLogin(ctx context.Context, token string, jwtKey []byte) (string, string, error) {
	clientID := viper.GetString("GOOGLE_CLIENT_ID")
	payload, err := idtoken.Validate(ctx, token, clientID)
	if err != nil {
		return "", "", errors.New("token de google inválido")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return "", "", errors.New("el token no contiene un email")
	}
	name, _ := payload.Claims["name"].(string)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		user = &models.User{
			Email:    email,
			Password: "", 
			Name:     name,
			RoleID:   1, // 1 = Cliente
			IsActive: true,
		}
		if createErr := s.userRepo.Create(ctx, user); createErr != nil {
			return "", "", errors.New("error creando usuario automáticamente")
		}
	}

	if !user.IsActive {
		return "", "", errors.New("usuario inactivo")
	}

	return s.generateAndSaveTokens(ctx, user, jwtKey)
}

func (s *authService) generateAndSaveTokens(ctx context.Context, user *models.User, jwtKey []byte) (string, string, error) {
	accExp := time.Now().Add(24 * time.Hour)
	refExp := time.Now().Add(7 * 24 * time.Hour)

	// Crear Access Token
	claims := &middleware.CustomClaims{
		UserID: user.ID,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accExp),
		},
	}
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accTokenString, err := accToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	// Persistir base token (Asumiendo TypeToken ID 1 = AccessToken localmente por ahora)
	userToken := &models.UserToken{
		UserID:          user.ID,
		UserTokenTypeID: 1, // Access Token 
		Token:           accTokenString,
		ExpiresAt:       accExp,
	}
	s.tokenRepo.CreateUserToken(ctx, userToken)

	// Crear Refresh Token
	refClaims := jwt.RegisteredClaims{
		Subject:   user.Email,
		ExpiresAt: jwt.NewNumericDate(refExp),
	}
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	refTokenString, err := refToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	if userToken.ID != 0 {
		refreshToken := &models.RefreshToken{
			UserTokenID: userToken.ID,
			Token:       refTokenString,
			ExpiresAt:   refExp,
		}
		s.tokenRepo.CreateRefreshToken(ctx, refreshToken)
	}

	return accTokenString, refTokenString, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.tokenRepo.RevokeUserToken(ctx, token)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string, jwtKey []byte) (string, string, error) {
	// 1. Check refresh token in db
	rt, err := s.tokenRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil || rt.IsRevoked || time.Now().After(rt.ExpiresAt) {
		return "", "", errors.New("refresh token inválido o expirado")
	}

	// 2. Extact user from ut
	ut, err := s.tokenRepo.FindUserToken(ctx, rt.UserToken.Token)
	if err != nil { // Fetch user token properly
		return "", "", errors.New("sesión corrupta")
	}
	user, err := s.userRepo.FindByID(ctx, ut.UserID)
	if err != nil || !user.IsActive {
		return "", "", errors.New("usuario no válido")
	}

	// 3. Mark old tokens as revoked
	s.tokenRepo.RevokeRefreshToken(ctx, refreshToken)
	s.tokenRepo.RevokeUserToken(ctx, ut.Token)

	// 4. Issue new pair
	return s.generateAndSaveTokens(ctx, user, jwtKey)
}
