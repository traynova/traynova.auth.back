package structs

type CreatePermissionRequest struct {
	RoleID     uint `json:"role_id" binding:"required"`
	ActionID   uint `json:"action_id" binding:"required"`
	ResourceID uint `json:"resource_id" binding:"required"`
}
