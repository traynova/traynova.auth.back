package src

import (
	"fmt"
	"traynova/src/common/config"
	"traynova/src/common/routes"
	"traynova/src/common/utils"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	EnvLocalConfigFile      = "./deployment/env_local.yaml"
	EnvProductionConfigFile = "./deployment/env.yaml"
)

var logger = utils.NewLogger()

func Run(isLocalEnv bool) {
	err := setupEnvironment(isLocalEnv)
	if err != nil {
		logger.Fatal("error setting up environment: %v", err)
	}

	conn, err := config.MigrateDB()
	if err != nil {
		logger.Fatal("error migrating database: %v", err)
	}
	defer conn.Close()

	initServer()
}

func setupEnvironment(isLocalEnv bool) error {
	viper.AutomaticEnv()

	if isLocalEnv {
		logger.Info("[TRAYNOVA_AUTH] server started in development mode")
		viper.SetConfigFile(EnvLocalConfigFile)
		err := viper.ReadInConfig()
		if err != nil {
			return fmt.Errorf("error reading env file: %v", err)
		}
	} else {
		logger.Info("[TRAYNOVA_AUTH] server started in production mode")
		// In production, rely on environment variables set in Render
	}

	return nil
}

func initServer() {
	address := viper.GetString("TRAYNOVA_AUTH_SERVER_ADDRESS")
	certFile := viper.GetString("TLS_CERT")
	certKey := viper.GetString("TLS_KEY")
	ginMode := viper.GetString("GIN_MODE")
	if address == "" {
		logger.Fatal("Server address env 'TRAYNOVA_SERVER_ADDRESS' not set")
	}
	if certFile == "" {
		logger.Fatal("TLS certificate file env 'TLS_CERT' not set")
	}
	if certKey == "" {
		logger.Fatal("TLS key file env 'TLS_KEY' not set")
	}
	if ginMode == "" {
		logger.Fatal("GIN mode env 'GIN_MODE' not set")
	}

	gin.SetMode(ginMode)
	serverInstance := gin.Default()
	routes.NewRoutesDefinition(serverInstance)

	logger.Info("[TRAYNOVA_AUTH] Start server on -> %v", address)
	if err := serverInstance.RunTLS(address, certFile, certKey); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
