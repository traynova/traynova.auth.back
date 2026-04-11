package app

import (
	"gestrym/src/core/auth/login/domain/ports"
	structs_request "gestrym/src/core/auth/login/domain/structs/request"
	structs_response "gestrym/src/core/auth/login/domain/structs/response"
)

type ILoginService interface {
	Login(req structs_request.LoginRequest) (*structs_response.LoginResponse, error)
}

type loginService struct {
	repo ports.ILoginRepository
}

func NewLoginService(repo ports.ILoginRepository) ILoginService {
	return &loginService{repo: repo}
}

func (s *loginService) Login(req structs_request.LoginRequest) (*structs_response.LoginResponse, error) {
	return nil, nil
}
