package endpoints

type CreatedResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type BadRequestResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}

type InternalServerResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}
