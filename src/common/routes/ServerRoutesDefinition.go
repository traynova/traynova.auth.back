package routes

import (
	"gestrym/docs"
	"gestrym/src/common/middleware"
	"gestrym/src/common/utils"
	"net/http"
	"sync"
	"time"

	"gestrym/src/common/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	roleApp "gestrym/src/core/roles/app"
	roleCtrl "gestrym/src/core/roles/infra/controller"
	roleRepo "gestrym/src/core/roles/infra/repository"

	authApp "gestrym/src/core/auth/app"
	authCtrl "gestrym/src/core/auth/infra/controller"
	authRepo "gestrym/src/core/auth/infra/repository"
	jwt_service "gestrym/src/core/jwt/app"
	jwtRepo "gestrym/src/core/jwt/infra"

	permApp "gestrym/src/core/permissions/app"
	permCtrl "gestrym/src/core/permissions/infra/controller"
	permRepo "gestrym/src/core/permissions/infra/repository"

	actionApp "gestrym/src/core/actions/app"
	actionCtrl "gestrym/src/core/actions/infra/controller"
	actionRepo "gestrym/src/core/actions/infra/repository"

	levelApp "gestrym/src/core/access_levels/app"
	levelCtrl "gestrym/src/core/access_levels/infra/controller"
	levelRepo "gestrym/src/core/access_levels/infra/repository"

	tokenTypeApp "gestrym/src/core/token_types/app"
	tokenTypeCtrl "gestrym/src/core/token_types/infra/controller"
	tokenTypeRepo "gestrym/src/core/token_types/infra/repository"
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
		docs.SwaggerInfo.BasePath = "/gestrym-auth"
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
	roleRepository := roleRepo.NewRoleRepository(db)
	permissionRepository := permRepo.NewPermissionRepository(db)
	actionRepository := actionRepo.NewActionRepository(db)
	accessLevelRepository := levelRepo.NewAccessLevelRepository(db)
	tokenTypeRepository := tokenTypeRepo.NewUserTokenTypeRepository(db)

	// Services
	roleService := roleApp.NewRoleService(roleRepository)
	permissionService := permApp.NewPermissionService(permissionRepository)
	actionService := actionApp.NewActionService(actionRepository)
	accessLevelService := levelApp.NewAccessLevelService(accessLevelRepository)
	tokenTypeService := tokenTypeApp.NewUserTokenTypeService(tokenTypeRepository)

	// Auth services
	authRepository := authRepo.NewAuthRepository(db)
	refreshTokenRepository := jwtRepo.NewRefreshTokenRepository(db)
	userTokenRepository := jwtRepo.NewUserTokenRepository(db)
	jwtService, err := jwt_service.NewJWTService(refreshTokenRepository, userTokenRepository)
	if err != nil {
		routesInstance.logger.Fatal("error al inicializar el servicio JWT: %v", err)
	}
	authService := authApp.NewAuthService(authRepository, jwtService, tokenTypeRepository)

	// Controllers
	rolePrivateController := roleCtrl.NewRolePrivateController(roleService)
	rolePublicController := roleCtrl.NewRolePublicController(roleService)
	permissionController := permCtrl.NewPermissionController(permissionService)
	actionController := actionCtrl.NewActionController(actionService)
	accessLevelController := levelCtrl.NewAccessLevelController(accessLevelService)
	tokenTypeController := tokenTypeCtrl.NewUserTokenTypeController(tokenTypeService)
	authPrivateController := authCtrl.NewAuthPrivateController(authService, r.logger)
	authPublicController := authCtrl.NewAuthPublicController(authService, r.logger)

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
	r.addPublicRoutes(rolePublicController, authPublicController)
	r.addPrivateRoutes(rolePrivateController, permissionController, actionController, accessLevelController, tokenTypeController, authPrivateController)
	r.addInternalRoutes()
	r.addProtectedRoutes()

}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {

	// Handle root
	serverInstance.GET("/", func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-auth OK...",
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

func (r *routesDefinition) addPublicRoutes(
	rolePublicController *roleCtrl.RolePublicController,
	authPublicController *authCtrl.AuthPublicController,
) {
	r.publicGroup.GET("/roles", rolePublicController.GetRoles)
	r.publicGroup.POST("/roles", rolePublicController.CreateRole)
	r.publicGroup.POST("/auth/register", authPublicController.Register())
	r.publicGroup.GET("/auth/confirm", authPublicController.ConfirmEmail())
	r.publicGroup.POST("/auth/password/recovery", authPublicController.RequestPasswordRecovery())
	r.publicGroup.POST("/auth/password/reset", authPublicController.ResetPassword())
	r.publicGroup.GET("/auth/validate", authPublicController.ValidateToken())
}

func (r *routesDefinition) addPrivateRoutes(
	rolePrivateController *roleCtrl.RolePrivateController,
	permissionController *permCtrl.PermissionController,
	actionController *actionCtrl.ActionController,
	accessLevelController *levelCtrl.AccessLevelController,
	tokenTypeController *tokenTypeCtrl.UserTokenTypeController,
	authPrivateController *authCtrl.AuthPrivateController,
) {

	r.privateGroup.PUT("/roles/:id", rolePrivateController.UpdateRole)
	r.privateGroup.DELETE("/roles/:id", rolePrivateController.DisableRole)
	r.privateGroup.POST("/permissions", permissionController.CreatePermission)

	// Catálogos (Solo admins, asumiendo rol 4 = Admin)
	adminAuth := middleware.RequireRoles(4)
	r.privateGroup.POST("/actions", adminAuth, actionController.CreateAction)
	r.privateGroup.GET("/actions", adminAuth, actionController.GetActions)

	r.privateGroup.POST("/access_levels", adminAuth, accessLevelController.CreateAccessLevel)
	r.privateGroup.GET("/access_levels", adminAuth, accessLevelController.GetAccessLevels)

	r.privateGroup.POST("/token_types", adminAuth, tokenTypeController.CreateUserTokenType)
	r.privateGroup.GET("/token_types", adminAuth, tokenTypeController.GetUserTokenTypes)

	r.privateGroup.GET("/auth/users", authPrivateController.GetUsers())
	r.privateGroup.GET("/auth/users/:id", authPrivateController.GetUserByID())
	r.privateGroup.PUT("/auth/users/:id", authPrivateController.UpdateUser())
	r.privateGroup.DELETE("/auth/users/:id", authPrivateController.DeleteUser())
	r.privateGroup.GET("/auth/relationships", authPrivateController.GetClientRelationships())
}

func (r *routesDefinition) addInternalRoutes() {

}

func (r *routesDefinition) addProtectedRoutes() {
}
