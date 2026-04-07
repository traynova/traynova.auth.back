package controller

import (
	"net/http"
	"traynova/src/core/actions/app"
	"traynova/src/core/actions/domain/structs"

	"github.com/gin-gonic/gin"
)

type ActionController struct {
	actionService app.IActionService
}

func NewActionController(as app.IActionService) *ActionController {
	return &ActionController{
		actionService: as,
	}
}

// @Summary Crear una Acción (CRUD)
// @Description Permite crear una nueva acción para el catálogo de permisos.
// @Tags Actions
// @Accept json
// @Produce json
// @Param request body structs.CreateActionRequest true "Datos de la acción"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /actions [post]
func (c *ActionController) CreateAction(ctx *gin.Context) {
	var req structs.CreateActionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	action, err := c.actionService.CreateAction(ctx.Request.Context(), req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando acción"})
		return
	}

	ctx.JSON(http.StatusCreated, action)
}

// @Summary Listar Acciones (CRUD)
// @Description Obtiene un listado de todas las acciones del sistema.
// @Tags Actions
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /actions [get]
func (c *ActionController) GetActions(ctx *gin.Context) {
	actions, err := c.actionService.GetActions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo acciones"})
		return
	}
	ctx.JSON(http.StatusOK, actions)
}
