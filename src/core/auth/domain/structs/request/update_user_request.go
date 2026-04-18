package structs_request

type UpdateUserRequest struct {
	FullName *string `json:"name"`
	Prefix   *string `json:"prefix"`
	Phone    *string `json:"phone"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	AvatarCollectionID *string `json:"avatar_collection_id"`
}

