package src

import (
	"fmt"
	"gestrym/src/common/config"
	"gestrym/src/common/routes"
	"gestrym/src/common/utils"

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
	if isLocalEnv {
		logger.Info("[GESTRYM_AUTH] server started in development mode")

		viper.SetConfigFile(EnvLocalConfigFile)

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("error reading local env file: %v", err)
		}
	} else {
		logger.Info("[GESTRYM_AUTH] server started in production mode")

		logger.Info("Using environment variables (Render)")
	}

	viper.AutomaticEnv()

	return nil
}

func initServer() {
	address := viper.GetString("GESTRYM_SERVER_ADDRESS")
	ginMode := viper.GetString("GIN_MODE")
	if address == "" {
		logger.Fatal("Server address env 'GESTRYM_SERVER_ADDRESS' not set")
	}
	if ginMode == "" {
		logger.Fatal("GIN mode env 'GIN_MODE' not set")
	}

	gin.SetMode(ginMode)
	serverInstance := gin.Default()
	routes.NewRoutesDefinition(serverInstance)

	logger.Info("[GESTRYM_AUTH] Start server on -> %v", address)
	if err := serverInstance.Run(address); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}
}
