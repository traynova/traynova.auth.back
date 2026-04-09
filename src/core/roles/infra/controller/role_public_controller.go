package controller

import (
	"gestrym/src/core/roles/app"
	structs_roles "gestrym/src/core/roles/domain/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RolePublicController struct {
	roleService app.IRoleService
}

func NewRolePublicController(rs app.IRoleService) *RolePublicController {
	return &RolePublicController{
		roleService: rs,
	}
}

// @Summary Obtener todos los Roles
// @Description Obtiene todos los roles disponibles.
// @Tags Roles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/roles [get]
func (c *RolePublicController) GetRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo roles"})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

// @Summary Crear Rol
// @Description Crea un nuevo rol en la base de datos (Admin, Gym, etc).
// @Tags Roles
// @Accept json
// @Produce json
// @Param request body structs_roles.CreateRoleRequest true "Estructura de creación del rol"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /public/roles [post]
func (c *RolePublicController) CreateRole(ctx *gin.Context) {
	var req structs_roles.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	role, err := c.roleService.CreateRole(ctx.Request.Context(), req.Name, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando rol"})
		return
	}

	ctx.JSON(http.StatusCreated, role)
}
