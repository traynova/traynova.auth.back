package config

import (
	"fmt"
	"gestrym/src/common/models"
	"gestrym/src/common/utils"
)

var logger = utils.NewLogger()

func MigrateDB() (IDatabaseConnection, error) {
	connection := NewPostgresConnection()
	db := connection.GetDB()

	//Se agregan los modelos de base de datos
	err := db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.Action{},
		&models.AccessLevel{},
		&models.UserTokenType{},
		&models.UserToken{},
		&models.RefreshToken{},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al migrar las entidades: %s", err.Error()))
		return nil, err
	}

	// Agrega un indice único a la tabla 'permissions'
	// haciendo que la combinación de 'role', 'action' y 'resource' sea única
	// para evitar convinaciones duplicadas
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_role_action_resource ON permissions (role_id, action_id, resource_id)").Error
	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al agregar el índice único a la tabla de permisos: %s", err.Error()))
		return nil, err
	}

	logger.Info("[OK] Todas las migraciones completadas exitosamente")
	return connection, nil
}
