package structs_response

type ClientResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type TrainerWithClientsResponse struct {
	TrainerID   uint             `json:"trainer_id"`
	TrainerName string           `json:"trainer_name"`
	TrainerEmail string          `json:"trainer_email"`
	TrainerPhone string          `json:"trainer_phone"`
	Clients     []ClientResponse `json:"clients"`
}

type TrainerClientGroupsResponse struct {
	IndependentClients []ClientResponse             `json:"independent_clients"`
	GymClients         []TrainerWithClientsResponse `json:"gym_clients"`
}

type GymTrainersResponse struct {
	Trainers []TrainerWithClientsResponse `json:"trainers"`
}
