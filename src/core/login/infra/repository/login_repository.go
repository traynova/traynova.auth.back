package repository

import (
	"gestrym/src/common/models"
	login_ports "gestrym/src/core/login/domain/ports"

	"gorm.io/gorm"
)

type loginRepository struct {
	db *gorm.DB
}

func NewLoginRepository(db *gorm.DB) login_ports.ILoginRepository {
	return &loginRepository{db: db}
}

func (r *loginRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *loginRepository) UpdateInitialLogin(email string) error {
	var user models.User

	if err := r.db.Where("email = ?", email).First(&user).UpdateColumn("initial_login = ?", true).Error; err != nil {
		return err
	}
	return nil
}
