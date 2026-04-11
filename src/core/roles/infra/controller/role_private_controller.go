package controller

import (
	"gestrym/src/core/roles/app"
	structs_roles "gestrym/src/core/roles/domain/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RolePrivateController struct {
	roleService app.IRoleService
}

func NewRolePrivateController(rs app.IRoleService) *RolePrivateController {
	return &RolePrivateController{
		roleService: rs,
	}
}

// @Summary Editar Rol
// @Description Actualiza los datos de un rol existente.
// @Tags Roles
// @Accept json
// @Produce json
// @Param id path int true "ID del rol"
// @Param request body structs_roles.UpdateRoleRequest true "Estructura de actualización del rol"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BasicAuth
// @Router /private/roles/{id} [put]
func (c *RolePrivateController) UpdateRole(ctx *gin.Context) {
	idParam := ctx.Param("id")
	roleID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de rol inválido"})
		return
	}

	var req structs_roles.UpdateRoleRequest
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
// @Security BasicAuth
// @Router /private/roles/{id} [delete]
func (c *RolePrivateController) DisableRole(ctx *gin.Context) {
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
