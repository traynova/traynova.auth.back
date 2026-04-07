package structs

type CreateActionRequest struct {
	Name string `json:"name" binding:"required"`
}
