package jwt_requests

import "gestrym/src/common/utils"

type GenerateJwtTokenRequest struct {
	UserID        uint   `json:"user_id" validate:"required"`
	RoleID        uint   `json:"role_id" validate:"required"`
	AccessLevelID uint   `json:"access_level_id" validate:"required"`
	Email         string `json:"email" validate:"required"`
	PhoneNumber   string `json:"phone_number"`
}

func (v *GenerateJwtTokenRequest) Validate() error {
	validator := utils.GetValidator()
	structErrors := validator.New().Struct(v)
	if structErrors != nil {
		return structErrors
	}

	return nil
}
