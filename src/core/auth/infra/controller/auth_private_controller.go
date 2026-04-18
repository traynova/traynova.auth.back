package controller

import (
	"gestrym/src/common/utils"
	"gestrym/src/core/auth/app"
	structs_request "gestrym/src/core/auth/domain/structs/request"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"mime/multipart"
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
//
//	@Security		BearerAuth
//
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

// @Summary Obtener usuario por ID
// @Description Devuelve el usuario con el ID especificado
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
//
//	@Security		BearerAuth
//
// @Router /private/auth/users/{id} [get]
func (a *AuthPrivateController) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		userID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
			return
		}

		response, err := a.authService.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Actualizar información de usuario
// @Description Modifica los campos del usuario especificado
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Param request body structs_request.UpdateUserRequest true "Datos a actualizar"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
//
//	@Security		BearerAuth
//
// @Router /private/auth/users/{id} [put]
func (a *AuthPrivateController) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		userID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
			return
		}

		var req structs_request.UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, err := a.authService.UpdateUser(uint(userID), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
	}
}

// @Summary Eliminar (soft delete) usuario
// @Description Desactiva el usuario marcando is_active en falso
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
//
//	@Security		BearerAuth
//
// @Router /private/auth/users/{id} [delete]
func (a *AuthPrivateController) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		userID, err := strconv.ParseUint(idParam, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario inválido"})
			return
		}

		if err := a.authService.DeleteUser(uint(userID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
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
//
//	@Security		BearerAuth
//
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

// @Summary Actualizar branding (avatar, logo, colores)
// @Description Permite subir imagen de avatar para el usuario, logo para el gimnasio y definir colores institucionales.
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file false "Imagen de avatar"
// @Param logo formData file false "Imagen de logo (solo Gym)"
// @Param primary_color formData string false "Color primario"
// @Param secondary_color formData string false "Color secundario"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
//
//	@Security		BearerAuth
//
// @Router /private/auth/branding [post]
func (a *AuthPrivateController) UpdateBranding() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No se pudo obtener el ID del usuario del token"})
			return
		}

		userID := userIDInterface.(uint)

		var avatar *multipart.FileHeader
		var logo *multipart.FileHeader

		avatar, _ = c.FormFile("avatar")
		logo, _ = c.FormFile("logo")

		primaryColor := c.PostForm("primary_color")
		secondaryColor := c.PostForm("secondary_color")

		var pColor, sColor *string
		if primaryColor != "" {
			pColor = &primaryColor
		}
		if secondaryColor != "" {
			sColor = &secondaryColor
		}

		err := a.authService.UpdateBranding(userID, avatar, logo, pColor, sColor)
		if err != nil {
			a.logger.Error("Error actualizando branding", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Branding actualizado correctamente"})
	}
}

