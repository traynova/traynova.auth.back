package controller

import (
	"net/http"
	"gestrym/src/core/permissions/app"
	"gestrym/src/core/permissions/domain/structs"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permissionService app.IPermissionService
}

func NewPermissionController(ps app.IPermissionService) *PermissionController {
	return &PermissionController{
		permissionService: ps,
	}
}

// @Summary Crear Permiso
// @Description Crea un nuevo permiso en el sistema.
// @Tags Permissions
// @Accept json
// @Produce json
// @Param request body structs.CreatePermissionRequest true "Datos del permiso"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /permissions [post]
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req structs.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	permission, err := c.permissionService.CreatePermission(ctx.Request.Context(), req.RoleID, req.ActionID, req.ResourceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando permiso"})
		return
	}

	ctx.JSON(http.StatusCreated, permission)
}
