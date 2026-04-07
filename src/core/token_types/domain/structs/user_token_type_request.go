package structs

type CreateUserTokenTypeRequest struct {
	Type string `json:"type" binding:"required"`
}
