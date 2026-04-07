package controller

import (
	"net/http"
	"traynova/src/core/users/app"
	"traynova/src/core/users/domain/structs"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService app.IUserService
}

func NewUserController(us app.IUserService) *UserController {
	return &UserController{
		userService: us,
	}
}

// @Summary Validar Autenticación
// @Description Permite testear la validez de un access token y extraer la metadata de sesión real.
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "id, email, name, role"
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /auth/validate [get]
func (c *UserController) Validate(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}

	user, err := c.userService.GetMe(ctx.Request.Context(), userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

// @Summary Crear Usuario Interno
// @Description Permite a perfiles elevados (Admin, Gym, Coach) registrar usuarios como clientes o sub-perfiles sin necesidad de auto-registro.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body structs.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} map[string]interface{} "message, id"
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /users/register [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req structs.CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	requesterRoleIDAny, exists := ctx.Get("role_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autorizado"})
		return
	}
	requesterRoleID := requesterRoleIDAny.(uint)

	// Reglas de negocio sobre quién puede crear qué rol:
	// Admin (4) puede crear a todos.
	// Gym (3) puede crear Coach (2) o Cliente (1).
	// Coach (2) puede crear Cliente (1).

	if requesterRoleID == 3 { // Gym
		if req.RoleID != 1 && req.RoleID != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Un Gym solo puede crear Coaches o Clientes"})
			return
		}
	} else if requesterRoleID == 2 { // Coach
		if req.RoleID != 1 {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Un Coach solo puede registrar Clientes"})
			return
		}
	} else if requesterRoleID != 4 { // Si no es admin y tampoco coach/gym, deniega
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No tienes permisos de creación de usuarios"})
		return
	}

	// Proceder con la creación
	user, err := c.userService.CreateUser(ctx.Request.Context(), req.Email, req.Phone, req.Name, req.RoleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno creando usuario"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado con éxito",
		"id":      user.ID,
	})
}
