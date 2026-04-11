package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gestrym/src/common/models"
	"gestrym/src/common/shared"
	ports_auth "gestrym/src/core/auth/domain/ports"
	structs_request "gestrym/src/core/auth/domain/structs/request"
	structs_response "gestrym/src/core/auth/domain/structs/response"
	jwt_service "gestrym/src/core/jwt/app"
	jwt_requests "gestrym/src/core/jwt/domain/structs/request"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	RegisterUser(req structs_request.RegisterRequest, userId uint) (*structs_response.RegisterResponse, error)
	GetAllUsers(page int, pageSize int, name string, dni string, email string, role_id uint) (shared.ResponsePaginate, error)
}

type authService struct {
	userRepo ports_auth.IAuthRepository
	jwt_app  jwt_service.IJWTService
}

func NewAuthService(ur ports_auth.IAuthRepository, jwtApp jwt_service.IJWTService) IAuthService {
	return &authService{
		userRepo: ur,
		jwt_app:  jwtApp,
	}
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendConfirmationEmail(user *models.User, name, token string) error {
	payload := map[string]interface{}{
		"user_id":       user.ID,
		"email":         user.Email,
		"user_name":     name,
		"confirm_token": token,
		"dashboard_url": "http://localhost:3000",
	}
	jsonPayload, _ := json.Marshal(payload)
	http.Post("http://localhost:8443/traynova-notification/public/send-confirmation", "application/json", bytes.NewBuffer(jsonPayload))

	return nil
}

func (s *authService) RegisterUser(req structs_request.RegisterRequest, userId uint) (*structs_response.RegisterResponse, error) {
	existingUser, _ := s.userRepo.ValidateEmail(req.Email)
	if existingUser.IsActive {
		return nil, errors.New("ya existe un usuario activo con ese email")
	}

	var errCreate error = nil
	var user *models.User
	if !existingUser.IsActive {
		_, errCreate = s.userRepo.UpdateUSer(existingUser)
		if errCreate != nil {
			return nil, errors.New("error reactivando usuario existente")
		}
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}

		newUser := &models.User{
			Password: string(hashedPassword),
		}

		user, errCreate = s.userRepo.CreateUser(newUser)
		if errCreate != nil {
			return nil, errors.New("error creando nuevo usuario")
		}
	}

	jwtRequest := jwt_requests.GenerateJwtTokenRequest{
		UserID:        user.ID,
		RoleID:        user.RoleID,
		AccessLevelID: 1,
		Email:         user.Email,
	}

	jwtToken, _ := s.jwt_app.GenerateJwtToken(jwtRequest, nil)

	// se registra el token con el estado activo tipo 2 "User Activation"
	userToken := models.UserToken{
		Token:           jwtToken,
		UserTokenTypeID: 2,
		UserID:          user.ID,
	}

	userTokenError := s.jwt_app.RegisterToken(userToken)
	if userTokenError != nil {
		return nil, errors.New("error registrando token de activación")
	}

	errNotifiction := sendConfirmationEmail(user, "ACTIVE_USER", jwtToken)
	if errNotifiction != nil {
		return nil, errors.New("error enviando email de confirmación")
	}

	authResponse := &structs_response.RegisterResponse{
		Email:  user.Email,
		Name:   user.FullName,
		Phone:  user.Phone,
		RoleID: user.RoleID,
		Token:  jwtToken,
	}

	return authResponse, nil

}

func (s *authService) GetAllUsers(page int, pageSize int, name string, dni string, email string, roleId uint) (shared.ResponsePaginate, error) {
	users, total, err := s.userRepo.GetAllUsers(page, pageSize, &name, &dni, &email)
	if err != nil {
		return shared.ResponsePaginate{}, err
	}

	var userList []interface{}
	for _, user := range users {
		userResponse := structs_response.GetAllUsersResponse{
			ID:       user.ID,
			Name:     user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			RoleID:   user.RoleID,
			RoleName: user.Role.Name,
		}
		userList = append(userList, userResponse)
	}

	return shared.ResponsePaginate{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		Results:  userList,
	}, nil

}
