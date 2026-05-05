package structs_response

type GroupResponse struct {
	ID      uint             `json:"id"`
	Name    string           `json:"name"`
	Members []ClientResponse `json:"members"`
}
