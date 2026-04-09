package controller

import (
	"net/http"
	"gestrym/src/core/access_levels/app"
	"gestrym/src/core/access_levels/domain/structs"

	"github.com/gin-gonic/gin"
)

type AccessLevelController struct {
	accessLevelService app.IAccessLevelService
}

func NewAccessLevelController(as app.IAccessLevelService) *AccessLevelController {
	return &AccessLevelController{
		accessLevelService: as,
	}
}

// @Summary Crear un Nivel de Acceso (CRUD)
// @Description Permite crear un nuevo nivel de acceso para la jerarquía.
// @Tags Access Levels
// @Accept json
// @Produce json
// @Param request body structs.CreateAccessLevelRequest true "Datos del nivel de acceso"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /access_levels [post]
func (c *AccessLevelController) CreateAccessLevel(ctx *gin.Context) {
	var req structs.CreateAccessLevelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	level, err := c.accessLevelService.CreateAccessLevel(ctx.Request.Context(), req.Name, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando nivel de acceso"})
		return
	}

	ctx.JSON(http.StatusCreated, level)
}

// @Summary Listar Niveles de acceso (CRUD)
// @Description Obtiene un listado de niveles del sistema.
// @Tags Access Levels
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /access_levels [get]
func (c *AccessLevelController) GetAccessLevels(ctx *gin.Context) {
	actions, err := c.accessLevelService.GetAccessLevels(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo niveles"})
		return
	}
	ctx.JSON(http.StatusOK, actions)
}
