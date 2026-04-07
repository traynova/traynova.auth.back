package controller

import (
	"net/http"
	"traynova/src/core/token_types/app"
	"traynova/src/core/token_types/domain/structs"

	"github.com/gin-gonic/gin"
)

type UserTokenTypeController struct {
	tokenTypeService app.IUserTokenTypeService
}

func NewUserTokenTypeController(ts app.IUserTokenTypeService) *UserTokenTypeController {
	return &UserTokenTypeController{
		tokenTypeService: ts,
	}
}

// @Summary Crear Tipo de Token (CRUD)
// @Description Permite crear un nuevo tipo de token.
// @Tags Token Types
// @Accept json
// @Produce json
// @Param request body structs.CreateUserTokenTypeRequest true "Datos del tipo de token"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /token_types [post]
func (c *UserTokenTypeController) CreateUserTokenType(ctx *gin.Context) {
	var req structs.CreateUserTokenTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Petición inválida"})
		return
	}

	result, err := c.tokenTypeService.CreateTokenType(ctx.Request.Context(), req.Type)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando el tipo de token"})
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

// @Summary Listar Tipos de Token (CRUD)
// @Description Obtiene un listado de los tipos.
// @Tags Token Types
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /token_types [get]
func (c *UserTokenTypeController) GetUserTokenTypes(ctx *gin.Context) {
	results, err := c.tokenTypeService.GetTokenTypes(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo los tipos"})
		return
	}
	ctx.JSON(http.StatusOK, results)
}
