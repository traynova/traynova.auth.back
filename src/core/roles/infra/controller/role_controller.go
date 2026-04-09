package controller

import (
	"gestrym/src/core/roles/app"
	"gestrym/src/core/roles/domain/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService app.IRoleService
}

func NewRoleController(rs app.IRoleService) *RoleController {
	return &RoleController{
		roleService: rs,
	}
}

// @Summary Crear Rol
// @Description Crea un nuevo rol en la base de datos (Admin, Gym, etc).
// @Tags Roles
// @Accept json
// @Produce json
// @Param request body structs.CreateRoleRequest true "Estructura de creación del rol"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req structs.CreateRoleRequest
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

// @Summary Editar Rol
// @Description Actualiza los datos de un rol existente.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param request body structs.UpdateRoleRequest true "Estructura de actualización del rol"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /roles/{id} [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	idParam := ctx.Param("id")
	roleID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	var req structs.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	if req.Name == "" && req.Description == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Al menos un campo debe ser actualizado"})
		return
	}

	role, err := c.roleService.UpdateRole(ctx.Request.Context(), uint(roleID), req.Name, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando rol"})
		return
	}

	ctx.JSON(http.StatusOK, role)
}

// @Summary Deshabilitar Rol
// @Description Marca un rol como inactivo en lugar de eliminarlo físicamente.
// @Tags Roles
// @Produce json
// @Param id path int true "ID del rol"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /roles/{id} [delete]
func (c *RoleController) DisableRole(ctx *gin.Context) {
	idParam := ctx.Param("id")
	roleID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	if err := c.roleService.DisableRole(ctx.Request.Context(), uint(roleID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deshabilitando rol"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rol deshabilitado correctamente"})
}

// @Summary Obtener todos los Roles
// @Description Obtiene todos los roles disponibles.
// @Tags Roles
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /roles [get]
func (c *RoleController) GetRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetRoles()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo roles"})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}
