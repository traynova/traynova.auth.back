package app

import (
	"errors"
	jwt_service "gestrym/src/core/jwt/app"
	jwt_requests "gestrym/src/core/jwt/domain/structs/request"
	login_ports "gestrym/src/core/login/domain/ports"
	structs_request "gestrym/src/core/login/domain/structs/request"
	structs_response "gestrym/src/core/login/domain/structs/response"

	"golang.org/x/crypto/bcrypt"
)

type ILoginService interface {
	Login(req structs_request.LoginRequest) (*structs_response.LoginResponse, error)
}

type loginService struct {
	repo   login_ports.ILoginRepository
	jwtApp jwt_service.IJWTService
}

func NewLoginService(repo login_ports.ILoginRepository, jwtApp jwt_service.IJWTService) ILoginService {
	return &loginService{repo: repo, jwtApp: jwtApp}
}

func (s *loginService) Login(req structs_request.LoginRequest) (*structs_response.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if !user.IsActive {
		return nil, errors.New("usuario inactivo")
	}

	if !user.EmailConfirmed {
		return nil, errors.New("email no confirmado")
	}

	if user.InitialLogin == false {
		err = s.repo.UpdateInitialLogin(user.Email)
		if err != nil {
			return nil, err
		}
	}

	accessToken, err := s.jwtApp.GenerateJwtToken(jwt_requests.GenerateJwtTokenRequest{
		UserID:        user.ID,
		RoleID:        user.RoleID,
		AccessLevelID: 1,
		Email:         user.Email,
		PhoneNumber:   user.Phone,
	}, nil)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtApp.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &structs_response.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RoleID:       user.RoleID,
		Email:        user.Email,
		ConfirmEmail: user.EmailConfirmed,
		InitialLogin: user.InitialLogin,
	}, nil
}
