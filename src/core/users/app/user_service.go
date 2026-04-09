package app

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"gestrym/src/common/models"
	"gestrym/src/core/users/domain/ports"
)

type IUserService interface {
	GetMe(ctx context.Context, userID uint) (*models.User, error)
	CreateUser(ctx context.Context, email, phone, name string, roleID uint) (*models.User, error)
}

type userService struct {
	userRepo ports.IUserRepository
}

func NewUserService(ur ports.IUserRepository) IUserService {
	return &userService{
		userRepo: ur,
	}
}

func (s *userService) GetMe(ctx context.Context, userID uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
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

func (s *userService) CreateUser(ctx context.Context, email, phone, name string, roleID uint) (*models.User, error) {
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
		go sendConfirmationEmail(user.ID, email, name, token)
	}
	return user, err
}
