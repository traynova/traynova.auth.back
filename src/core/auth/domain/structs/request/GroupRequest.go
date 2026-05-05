package structs_request

type CreateGroupRequest struct {
	Name    string `json:"name" binding:"required"`
	UserIDs []uint `json:"user_ids"`
}

type UpdateGroupRequest struct {
	Name    *string `json:"name"`
	UserIDs []uint  `json:"user_ids"`
}
