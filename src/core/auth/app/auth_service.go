package app

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gestrym/src/common/middleware"
	"gestrym/src/common/models"
	"gestrym/src/common/shared"
	ports_auth "gestrym/src/core/auth/domain/ports"
	structs_request "gestrym/src/core/auth/domain/structs/request"
	structs_response "gestrym/src/core/auth/domain/structs/response"
	jwt_service "gestrym/src/core/jwt/app"
	jwt_requests "gestrym/src/core/jwt/domain/structs/request"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	RegisterUser(req structs_request.RegisterRequest, userId uint) (*structs_response.RegisterResponse, error)
	GetAllUsers(page int, pageSize int, name string, dni string, email string, role_id uint) (shared.ResponsePaginate, error)
	GetClientsByUser(userID uint, roleID uint) (interface{}, error)
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
	existingUser, err := s.userRepo.ValidateEmail(req.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("error validating email")
	}

	var user *models.User
	if existingUser != nil && existingUser.ID != 0 {
		if existingUser.RoleID != req.RoleID {
			return nil, errors.New("ya existe un usuario activo con ese email y rol diferente")
		}

		if !existingUser.IsActive {
			existingUser.IsActive = true
		}

		if err := s.updateUserFields(existingUser, req); err != nil {
			return nil, err
		}

		user, err = s.userRepo.UpdateUSer(existingUser)
		if err != nil {
			return nil, errors.New("error actualizando usuario existente")
		}
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}

		newUser := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
			FullName: req.FullName,
			Prefix:   req.Prefix,
			Phone:    req.Phone,
			RoleID:   req.RoleID,
			IsActive: true,
		}

		user, err = s.userRepo.CreateUser(newUser)
		if err != nil {
			return nil, errors.New("error creando nuevo usuario")
		}
	}

	if err := s.attachRelationship(req, user); err != nil {
		return nil, err
	}

	jwtRequest := jwt_requests.GenerateJwtTokenRequest{
		UserID:        user.ID,
		RoleID:        user.RoleID,
		AccessLevelID: 1,
		Email:         user.Email,
	}

	jwtToken, _ := s.jwt_app.GenerateJwtToken(jwtRequest, nil)

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

func (s *authService) updateUserFields(user *models.User, req structs_request.RegisterRequest) error {
	user.FullName = req.FullName
	user.Prefix = req.Prefix
	user.Phone = req.Phone
	user.RoleID = req.RoleID

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("error al hashear la contraseña")
		}
		user.Password = string(hashedPassword)
	}

	return nil
}

func (s *authService) attachRelationship(req structs_request.RegisterRequest, user *models.User) error {
	switch req.RegistrationSource {
	case structs_request.RegistrationSourceGym:
		if req.RoleID == middleware.RoleCoach {
			if req.SourceID == nil {
				return errors.New("source_id es requerido para registrar un entrenador desde un gimnasio")
			}
			return s.createOrUpdateTrainerProfile(user.ID, req.SourceID)
		}

		if req.RoleID == middleware.RoleCliente {
			if req.SourceID == nil {
				return errors.New("source_id es requerido para registrar un cliente desde un gimnasio")
			}
			return s.createGymClient(user.ID, *req.SourceID)
		}

	case structs_request.RegistrationSourceTrainer:
		if req.RoleID == middleware.RoleCliente {
			if req.SourceID == nil {
				return errors.New("source_id es requerido para registrar un cliente desde un entrenador")
			}
			return s.createTrainerClient(user.ID, *req.SourceID)
		}

	case structs_request.RegistrationSourceSelf:
		if req.RoleID == middleware.RoleCoach {
			return s.createOrUpdateTrainerProfile(user.ID, nil)
		}
	}

	return nil
}

func (s *authService) createOrUpdateTrainerProfile(userID uint, gymUserID *uint) error {
	trainerProfile, err := s.userRepo.GetTrainerProfileByUserIDAndGymID(userID, gymUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			trainerProfile = &models.TrainerProfile{
				UserID: userID,
				GimID:  gymUserID,
			}
			_, err = s.userRepo.CreateTrainerProfile(trainerProfile)
			return err
		}
		return err
	}

	trainerProfile.GimID = gymUserID
	_, err = s.userRepo.UpdateTrainerProfile(trainerProfile)
	return err
}

func (s *authService) createTrainerClient(clientID, coachUserID uint) error {
	trainerProfile, err := s.userRepo.GetTrainerProfileByUserIDAndGymID(coachUserID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			trainerProfile = &models.TrainerProfile{
				UserID: coachUserID,
				GimID:  nil,
			}
			trainerProfile, err = s.userRepo.CreateTrainerProfile(trainerProfile)
			if err != nil {
				return err
			}
		} else {
			return errors.New("no se encontró el perfil del entrenador independiente")
		}
	}

	_, err = s.userRepo.GetTrainerClientByProfileAndClient(trainerProfile.ID, clientID)
	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	trainerClient := &models.TrainerClient{
		TrainerProfileID: trainerProfile.ID,
		ClientID:         clientID,
	}

	_, err = s.userRepo.CreateTrainerClient(trainerClient)
	return err
}

func (s *authService) createGymClient(clientID, gymUserID uint) error {
	_, err := s.userRepo.GetGymClientByGymAndClient(gymUserID, clientID)
	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	gymClient := &models.GymClient{
		GymUserID: gymUserID,
		ClientID:  clientID,
	}

	_, err = s.userRepo.CreateGymClient(gymClient)
	return err
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

func (s *authService) GetClientsByUser(userID uint, roleID uint) (interface{}, error) {
	switch roleID {
	case middleware.RoleCoach:
		return s.getClientsForTrainer(userID)
	case middleware.RoleGym:
		return s.getTrainersForGym(userID)
	default:
		return nil, errors.New("acceso no permitido para este rol")
	}
}

func (s *authService) getClientsForTrainer(userID uint) (interface{}, error) {
	independentClients := []structs_response.ClientResponse{}
	gymClients := []structs_response.TrainerWithClientsResponse{}

	profiles, err := s.userRepo.GetTrainerProfilesByUserID(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	for _, profile := range profiles {
		clients, err := s.userRepo.GetTrainerClientsByProfileID(profile.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var clientResponses []structs_response.ClientResponse
		for _, client := range clients {
			clientResponses = append(clientResponses, structs_response.ClientResponse{
				ID:    client.ID,
				Email: client.Email,
				Name:  client.FullName,
				Phone: client.Phone,
			})
		}

		if profile.GimID == nil {
			independentClients = append(independentClients, clientResponses...)
			continue
		}

		gymClients = append(gymClients, structs_response.TrainerWithClientsResponse{
			TrainerID:    profile.UserID,
			TrainerName:  profile.User.FullName,
			TrainerEmail: profile.User.Email,
			TrainerPhone: profile.User.Phone,
			Clients:      clientResponses,
		})
	}

	return structs_response.TrainerClientGroupsResponse{
		IndependentClients: independentClients,
		GymClients:         gymClients,
	}, nil
}

func (s *authService) getTrainersForGym(gymUserID uint) (interface{}, error) {
	profiles, err := s.userRepo.GetGymTrainersByGymUserID(gymUserID)
	if err != nil {
		return nil, err
	}

	trainers := []structs_response.TrainerWithClientsResponse{}
	for _, profile := range profiles {
		clients, err := s.userRepo.GetTrainerClientsByProfileID(profile.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		var clientResponses []structs_response.ClientResponse
		for _, client := range clients {
			clientResponses = append(clientResponses, structs_response.ClientResponse{
				ID:    client.ID,
				Email: client.Email,
				Name:  client.FullName,
				Phone: client.Phone,
			})
		}

		trainers = append(trainers, structs_response.TrainerWithClientsResponse{
			TrainerID:    profile.UserID,
			TrainerName:  profile.User.FullName,
			TrainerEmail: profile.User.Email,
			TrainerPhone: profile.User.Phone,
			Clients:      clientResponses,
		})
	}

	return structs_response.GymTrainersResponse{Trainers: trainers}, nil
}
