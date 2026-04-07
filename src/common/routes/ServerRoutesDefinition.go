package routes

import (
	"net/http"
	"sync"
	"time"
	"traynova/docs"
	"traynova/src/common/middleware"
	"traynova/src/common/utils"

	"traynova/src/common/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	authApp "traynova/src/core/auth/app"
	authCtrl "traynova/src/core/auth/infra/controller"
	authRepo "traynova/src/core/auth/infra/repository"

	roleApp "traynova/src/core/roles/app"
	roleCtrl "traynova/src/core/roles/infra/controller"
	roleRepo "traynova/src/core/roles/infra/repository"

	permApp "traynova/src/core/permissions/app"
	permCtrl "traynova/src/core/permissions/infra/controller"
	permRepo "traynova/src/core/permissions/infra/repository"

	userApp "traynova/src/core/users/app"
	userCtrl "traynova/src/core/users/infra/controller"
	userRepo "traynova/src/core/users/infra/repository"

	actionApp "traynova/src/core/actions/app"
	actionCtrl "traynova/src/core/actions/infra/controller"
	actionRepo "traynova/src/core/actions/infra/repository"

	levelApp "traynova/src/core/access_levels/app"
	levelCtrl "traynova/src/core/access_levels/infra/controller"
	levelRepo "traynova/src/core/access_levels/infra/repository"

	tokenTypeApp "traynova/src/core/token_types/app"
	tokenTypeCtrl "traynova/src/core/token_types/infra/controller"
	tokenTypeRepo "traynova/src/core/token_types/infra/repository"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.BasePath = "/traynova-auth"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	// Aplica el middleware CORS
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	// Add default routes
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB
	db := config.NewPostgresConnection().GetDB()

	// Repositories
	userRepository := userRepo.NewUserRepository(db)
	roleRepository := roleRepo.NewRoleRepository(db)
	permissionRepository := permRepo.NewPermissionRepository(db)
	tokenRepository := authRepo.NewTokenRepository(db)
	actionRepository := actionRepo.NewActionRepository(db)
	accessLevelRepository := levelRepo.NewAccessLevelRepository(db)
	tokenTypeRepository := tokenTypeRepo.NewUserTokenTypeRepository(db)

	// Services
	authService := authApp.NewAuthService(userRepository, tokenRepository)
	userService := userApp.NewUserService(userRepository)
	roleService := roleApp.NewRoleService(roleRepository)
	permissionService := permApp.NewPermissionService(permissionRepository)
	actionService := actionApp.NewActionService(actionRepository)
	accessLevelService := levelApp.NewAccessLevelService(accessLevelRepository)
	tokenTypeService := tokenTypeApp.NewUserTokenTypeService(tokenTypeRepository)

	// Controllers
	authController := authCtrl.NewAuthController(authService)
	userController := userCtrl.NewUserController(userService)
	roleController := roleCtrl.NewRoleController(roleService)
	permissionController := permCtrl.NewPermissionController(permissionService)
	actionController := actionCtrl.NewActionController(actionService)
	accessLevelController := levelCtrl.NewAccessLevelController(accessLevelService)
	tokenTypeController := tokenTypeCtrl.NewUserTokenTypeController(tokenTypeService)

	// Add server group
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add groups
	r.publicGroup = r.serverGroup.Group("/public")
	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Add middleware to private group
	r.privateGroup.Use(middleware.SetupJWTMiddleware())

	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	// Add routes to groups
	r.addPublicRoutes(authController)
	r.addPrivateRoutes(userController, roleController, permissionController, actionController, accessLevelController, tokenTypeController, authController)
	r.addInternalRoutes()
	r.addProtectedRoutes()

}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {

	// Handle root
	serverInstance.GET("/", func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "OK",
			"message": "traynova-auth OK...",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusOK, response)
	})

	// Handle 404
	serverInstance.NoRoute(func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusNotFound, response)
	})
}

func (r *routesDefinition) addPublicRoutes(authController *authCtrl.AuthController) {
	r.publicGroup.POST("/auth/login", authController.Login)
	r.publicGroup.POST("/auth/google", authController.GoogleLogin)
	r.publicGroup.POST("/auth/register", authController.Register)
}

func (r *routesDefinition) addPrivateRoutes(
	userController *userCtrl.UserController,
	roleController *roleCtrl.RoleController,
	permissionController *permCtrl.PermissionController,
	actionController *actionCtrl.ActionController,
	accessLevelController *levelCtrl.AccessLevelController,
	tokenTypeController *tokenTypeCtrl.UserTokenTypeController,
	authController *authCtrl.AuthController,
) {
	r.privateGroup.GET("/auth/validate", userController.Validate)

	// Refresh / Logout endpoints
	r.privateGroup.POST("/auth/refresh", authController.Refresh)
	r.privateGroup.POST("/auth/logout", authController.Logout)

	// El registro privado (dashboard) solo está permitido a Gym, Coach o Admin (roles 2, 3, 4)
	r.privateGroup.POST("/users/register", middleware.RequireRoles(2, 3, 4), userController.CreateUser)

	r.privateGroup.POST("/roles", roleController.CreateRole)
	r.privateGroup.PUT("/roles/:id", roleController.UpdateRole)
	r.privateGroup.DELETE("/roles/:id", roleController.DisableRole)
	r.privateGroup.POST("/permissions", permissionController.CreatePermission)

	// Catálogos (Solo admins, asumiendo rol 4 = Admin)
	adminAuth := middleware.RequireRoles(4)
	r.privateGroup.POST("/actions", adminAuth, actionController.CreateAction)
	r.privateGroup.GET("/actions", adminAuth, actionController.GetActions)

	r.privateGroup.POST("/access_levels", adminAuth, accessLevelController.CreateAccessLevel)
	r.privateGroup.GET("/access_levels", adminAuth, accessLevelController.GetAccessLevels)

	r.privateGroup.POST("/token_types", adminAuth, tokenTypeController.CreateUserTokenType)
	r.privateGroup.GET("/token_types", adminAuth, tokenTypeController.GetUserTokenTypes)
}

func (r *routesDefinition) addInternalRoutes() {

}

func (r *routesDefinition) addProtectedRoutes() {
}
