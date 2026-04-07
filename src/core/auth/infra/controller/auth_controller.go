package controller

import (
	"net/http"
	"traynova/src/core/auth/app"
	"traynova/src/core/auth/domain/structs"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type AuthController struct {
	authService app.IAuthService
}

func NewAuthController(as app.IAuthService) *AuthController {
	return &AuthController{
		authService: as,
	}
}

// @Summary Iniciar Sesión Tradicional
// @Description Realiza el login de un usuario registrado usando email y password, y devuelve un JWT.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs.LoginRequest true "Credenciales de login"
// @Success 200 {object} map[string]interface{} "access_token: xxx, refresh_token: yyy"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req structs.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	jwtKey := viper.GetString("JWT_KEY")
	accToken, refToken, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password, []byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accToken, "refresh_token": refToken})
}

// @Summary Registrar Cliente o Coach
// @Description Crea un nuevo usuario públicamente. Solo permite roles Cliente (1) o Coach (2).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs.RegisterRequest true "Datos del usuario"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req structs.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	// Restricción pública: Sólo se pueden crear clientes (1) o coach (2) asumiendo
	// IDs de roles por ahora o validando constantes. Asumiremos 1=Cliente, 2=Coach
	if req.RoleID != 1 && req.RoleID != 2 {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Esta ruta pública solo admite registro de Cliente o Coach"})
		return
	}

	err := c.authService.RegisterUser(ctx.Request.Context(), req.Email, req.Phone, req.Name, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando el usuario"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Usuario creado exitosamente"})
}

// @Summary Login con Google
// @Description Inicia sesión utilizando un id_token provisto por Google. Si el usuario no existe, lo registra como Cliente (1).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs.GoogleLoginRequest true "Token de Google"
// @Success 200 {object} map[string]interface{} "access_token: xxx, refresh_token: yyy"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/google [post]
func (c *AuthController) GoogleLogin(ctx *gin.Context) {
	var req structs.GoogleLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida: falta id_token"})
		return
	}

	jwtKey := viper.GetString("JWT_KEY")
	accToken, refToken, err := c.authService.GoogleLogin(ctx.Request.Context(), req.IDToken, []byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accToken, "refresh_token": refToken})
}

// @Summary Refrescar Token
// @Description Intercambia un refresh_token válido por un nuevo set de tokens de acceso y refresco.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs.RefreshRequest true "Refresh Token"
// @Success 200 {object} map[string]interface{} "access_token: xxx, refresh_token: yyy"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req structs.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	jwtKey := viper.GetString("JWT_KEY")
	accToken, refToken, err := c.authService.Refresh(ctx.Request.Context(), req.RefreshToken, []byte(jwtKey))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"access_token": accToken, "refresh_token": refToken})
}

// @Summary Cerrar Sesión
// @Description Invalida el token JWT actual (o provisto) en la base de datos de tokens activos.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs.LogoutRequest true "Token a revocar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var req structs.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	err := c.authService.Logout(ctx.Request.Context(), req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo cerrar la sesión"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada con éxito"})
}
