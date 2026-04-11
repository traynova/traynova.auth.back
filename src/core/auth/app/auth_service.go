package app

import (
	"bytes"
	"context"
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
	token_types_ports "gestrym/src/core/token_types/domain/ports"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IAuthService interface {
	RegisterUser(req structs_request.RegisterRequest) (*structs_response.RegisterResponse, error)
	GetAllUsers(page int, pageSize int, name string, dni string, email string, role_id uint) (shared.ResponsePaginate, error)
	GetUserByID(userID uint) (*structs_response.GetUserResponse, error)
	UpdateUser(userID uint, req structs_request.UpdateUserRequest) (*structs_response.GetUserResponse, error)
	DeleteUser(userID uint) error
	ActivateUser(token string) (*structs_response.GetUserResponse, error)
	RequestPasswordRecovery(email string) error
	ResetPassword(req structs_request.PasswordResetRequest) (*structs_response.GetUserResponse, error)
	ValidateToken(token string) (*structs_response.ValidateTokenResponse, error)
	GetClientsByUser(userID uint, roleID uint) (interface{}, error)
}

type authService struct {
	userRepo      ports_auth.IAuthRepository
	jwt_app       jwt_service.IJWTService
	tokenTypeRepo token_types_ports.IUserTokenTypeRepository
}

func NewAuthService(ur ports_auth.IAuthRepository, jwtApp jwt_service.IJWTService, tokenTypeRepo token_types_ports.IUserTokenTypeRepository) IAuthService {
	return &authService{
		userRepo:      ur,
		jwt_app:       jwtApp,
		tokenTypeRepo: tokenTypeRepo,
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

func sendPasswordRecoveryEmail(user *models.User, token string) error {
	payload := map[string]interface{}{
		"user_id":       user.ID,
		"email":         user.Email,
		"reset_token":   token,
		"dashboard_url": "http://localhost:3000",
	}
	jsonPayload, _ := json.Marshal(payload)
	http.Post("http://localhost:8443/traynova-notification/public/send-password-recovery", "application/json", bytes.NewBuffer(jsonPayload))

	return nil
}

func (s *authService) RegisterUser(req structs_request.RegisterRequest) (*structs_response.RegisterResponse, error) {
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

	jwtToken, err := s.jwt_app.GenerateJwtToken(jwtRequest, nil)
	if err != nil {
		return nil, errors.New("error generando token JWT")
	}

	activationTokenType, err := s.tokenTypeRepo.FindByType(context.Background(), models.UserTokenTypeActivation)
	if err != nil {
		return nil, errors.New("error buscando tipo de token de activación")
	}
	if activationTokenType == nil || activationTokenType.ID == 0 {
		return nil, errors.New("tipo de token de activación no encontrado")
	}

	claims, err := s.jwt_app.ValidateJwtToken(jwtToken)
	if err != nil {
		return nil, errors.New("error validando token JWT para registro")
	}

	expiresAt := time.Now()
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	userToken := models.UserToken{
		Token:           jwtToken,
		UserTokenTypeID: activationTokenType.ID,
		UserID:          user.ID,
		ExpiresAt:       expiresAt,
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
		Id:     user.ID,
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
			return s.createOrUpdateTrainerProfile(user.ID, req.SourceID, req)
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
			return s.createOrUpdateTrainerProfile(user.ID, nil, req)
		}
		if req.RoleID == middleware.RoleGym {
			return s.createOrUpdateGymProfile(user.ID, req)
		}
	}

	return nil
}

func (s *authService) createOrUpdateGymProfile(userID uint, req structs_request.RegisterRequest) error {
	if req.City == nil || req.Department == nil || req.Country == nil {
		return errors.New("city, department y country son requeridos para registrar un gimnasio")
	}

	gymProfile, err := s.userRepo.GetGymProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			gymProfile = &models.GymProfile{
				UserID:         userID,
				City:           *req.City,
				Department:     *req.Department,
				Country:        *req.Country,
				Workstation:    req.Workstation,
				PrimaryColor:   req.PrimaryColor,
				SecondaryColor: req.SecondaryColor,
				ReferralCode:   req.ReferralCode,
			}
			_, err = s.userRepo.CreateGymProfile(gymProfile)
			return err
		}
		return err
	}

	if req.City != nil {
		gymProfile.City = *req.City
	}
	if req.Department != nil {
		gymProfile.Department = *req.Department
	}
	if req.Country != nil {
		gymProfile.Country = *req.Country
	}
	gymProfile.Workstation = req.Workstation
	gymProfile.PrimaryColor = req.PrimaryColor
	gymProfile.SecondaryColor = req.SecondaryColor
	gymProfile.ReferralCode = req.ReferralCode

	_, err = s.userRepo.UpdateGymProfile(gymProfile)
	return err
}

func (s *authService) createOrUpdateTrainerProfile(userID uint, gymUserID *uint, req structs_request.RegisterRequest) error {
	trainerProfile, err := s.userRepo.GetTrainerProfileByUserIDAndGymID(userID, gymUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			trainerProfile = &models.TrainerProfile{
				UserID:         userID,
				GimID:          gymUserID,
				PrimaryColor:   req.PrimaryColor,
				SecondaryColor: req.SecondaryColor,
				FilesID:        req.AvatarFileID,
				ReferralCode:   req.ReferralCode,
			}
			_, err = s.userRepo.CreateTrainerProfile(trainerProfile)
			return err
		}
		return err
	}

	trainerProfile.GimID = gymUserID
	if req.PrimaryColor != nil {
		trainerProfile.PrimaryColor = req.PrimaryColor
	}
	if req.SecondaryColor != nil {
		trainerProfile.SecondaryColor = req.SecondaryColor
	}
	if req.AvatarFileID != nil {
		trainerProfile.FilesID = req.AvatarFileID
	}
	trainerProfile.ReferralCode = req.ReferralCode

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

func (s *authService) GetUserByID(userID uint) (*structs_response.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return buildUserResponse(user), nil
}

func (s *authService) UpdateUser(userID uint, req structs_request.UpdateUserRequest) (*structs_response.GetUserResponse, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	if req.Email != nil && *req.Email != user.Email {
		existing, err := s.userRepo.ValidateEmail(*req.Email)
		if err == nil && existing != nil && existing.ID != 0 && existing.ID != userID {
			return nil, errors.New("ya existe un usuario con ese email")
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Prefix != nil {
		user.Prefix = *req.Prefix
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Password != nil && *req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("error al hashear la contraseña")
		}
		user.Password = string(hashedPassword)
	}

	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return buildUserResponse(updatedUser), nil
}

func (s *authService) DeleteUser(userID uint) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	_, err = s.userRepo.UpdateUser(user)
	return err
}

func (s *authService) ActivateUser(token string) (*structs_response.GetUserResponse, error) {
	if err := s.jwt_app.ChecUserTokenUsed(token); err != nil {
		return nil, errors.New("token de activación inválido o no registrado")
	}

	claims, err := s.jwt_app.ValidateJwtToken(token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	user.IsActive = true
	user.EmailConfirmed = true
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	if err := s.jwt_app.DeleteUserToken(token); err != nil {
		return nil, err
	}

	return buildUserResponse(updatedUser), nil
}

func (s *authService) RequestPasswordRecovery(email string) error {
	user, err := s.userRepo.ValidateEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	token, err := s.jwt_app.GenerateJwtToken(jwt_requests.GenerateJwtTokenRequest{
		UserID:        user.ID,
		RoleID:        user.RoleID,
		AccessLevelID: 1,
		Email:         user.Email,
		PhoneNumber:   user.Phone,
	}, nil)
	if err != nil {
		return err
	}

	passwordRecoveryTokenType, err := s.tokenTypeRepo.FindByType(context.Background(), models.UserTokenTypePasswordRecovery)
	if err != nil {
		return errors.New("error buscando tipo de token de recuperación")
	}
	if passwordRecoveryTokenType == nil || passwordRecoveryTokenType.ID == 0 {
		return errors.New("tipo de token de recuperación no encontrado")
	}

	claims, err := s.jwt_app.ValidateJwtToken(token)
	if err != nil {
		return err
	}

	expiresAt := time.Now()
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}

	userToken := models.UserToken{
		Token:           token,
		UserTokenTypeID: passwordRecoveryTokenType.ID,
		UserID:          user.ID,
		ExpiresAt:       expiresAt,
	}

	if err := s.jwt_app.RegisterToken(userToken); err != nil {
		return err
	}

	if err := sendPasswordRecoveryEmail(user, token); err != nil {
		return errors.New("error enviando email de recuperación")
	}

	return nil
}

func (s *authService) ResetPassword(req structs_request.PasswordResetRequest) (*structs_response.GetUserResponse, error) {
	if err := s.jwt_app.ChecUserTokenUsed(req.Token); err != nil {
		return nil, errors.New("token de recuperación inválido o no registrado")
	}

	claims, err := s.jwt_app.ValidateJwtToken(req.Token)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al hashear la contraseña")
	}

	user.Password = string(hashedPassword)
	updatedUser, err := s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	if err := s.jwt_app.DeleteUserToken(req.Token); err != nil {
		return nil, err
	}

	return buildUserResponse(updatedUser), nil
}

func (s *authService) ValidateToken(token string) (*structs_response.ValidateTokenResponse, error) {
	claims, err := s.jwt_app.ValidateJwtToken(token)
	if err != nil {
		return nil, err
	}

	var expiresAt int64
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time.Unix()
	}

	return &structs_response.ValidateTokenResponse{
		Valid:         true,
		UserID:        claims.UserID,
		RoleID:        claims.RoleID,
		AccessLevelID: claims.AccessLevelID,
		Email:         claims.Subject,
		ExpiresAt:     expiresAt,
	}, nil
}

func buildUserResponse(user *models.User) *structs_response.GetUserResponse {
	return &structs_response.GetUserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.FullName,
		Phone:          user.Phone,
		Prefix:         user.Prefix,
		RoleID:         user.RoleID,
		RoleName:       user.Role.Name,
		IsActive:       user.IsActive,
		EmailConfirmed: user.EmailConfirmed,
	}
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
