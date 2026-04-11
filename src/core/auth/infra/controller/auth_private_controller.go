package controller

import (
	"gestrym/src/common/utils"
	"gestrym/src/core/auth/app"
	structs_request "gestrym/src/core/auth/domain/structs/request"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthPrivateController struct {
	authService app.IAuthService
	logger      utils.ILogger
}

func NewAuthPrivateController(as app.IAuthService, logger utils.ILogger) *AuthPrivateController {
	return &AuthPrivateController{
		authService: as,
		logger:      logger,
	}
}

// @Summary Registrar Coach o Gimnasio o Cliente
// @Description Crea un nuevo usuario públicamente. Solo permite roles Cliente (1) o Coach (2) o Gym (3).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body structs_request.RegisterRequest true "Datos del usuario"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /private/auth/register [post]
func (a *AuthPrivateController) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		createUserRequest := &structs_request.RegisterRequest{}

		if err := c.ShouldBindJSON(&createUserRequest); err != nil {
			a.logger.Error("Error while binding JSON", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el ID del usuario"})
			return
		}

		userId, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del usuario en formato incorrecto"})
			return
		}

		response, err := a.authService.RegisterUser(*createUserRequest, userId)
		if err != nil {
			a.logger.Error("Error while creating user", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"response": response})
	}
}

// @Summary Consultar todos los usuarios registrados
// @Description Obtiene una lista de todos los usuarios registrados en el sistema
// @Tags Auth
// @Accept json
// @Produce json
// @Param			page		query		int		false	"Número de página"
// @Param			page_size	query		int		false	"Tamaño de la página"
// @Param			name		query		string	false	"Nombre del usuario"
// @Param			dni			query		string	false	"DNI del usuario"
// @Param			email		query		string	false	"email del usuario"
// @Param			role_id		query		int		false	"ID del rol del usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /private/auth/users [get]
func (a *AuthPrivateController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		name := c.Query("name")
		dni := c.Query("dni")
		email := c.Query("email")
		roleIDStr := c.Query("role_id")

		var roleID uint
		if roleIDStr != "" {
			parsedRoleID, err := strconv.ParseUint(roleIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role_id format"})
				return
			}
			roleID = uint(parsedRoleID)
		}

		response, err := a.authService.GetAllUsers(page, pageSize, name, dni, email, roleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Obtener relaciones de clientes del entrenador o entrenadores del gimnasio
// @Description Si el usuario es coach, devuelve clientes independientes y clientes de gimnasio separados. Si el usuario es gym, devuelve entrenadores y sus clientes de gimnasio.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /private/auth/relationships [get]
func (a *AuthPrivateController) GetClientRelationships() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el ID del usuario"})
			return
		}

		userRoleInterface, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo obtener el rol del usuario"})
			return
		}

		userId, ok := userIDInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del usuario en formato incorrecto"})
			return
		}

		roleId, ok := userRoleInterface.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ID del rol en formato incorrecto"})
			return
		}

		response, err := a.authService.GetClientsByUser(userId, roleId)
		if err != nil {
			a.logger.Error("Error while fetching client relationships", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}
