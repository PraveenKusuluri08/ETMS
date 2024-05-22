package endpoints

type CreatedResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ErrorMessage struct {
	Name string `json:"name"`
}

type BadRequestResponse struct {
	Msg    ErrorMessage `json:"message"`
	Status string       `json:"status"`
	Error  string       `json:"error"`
}

type InternalServerResponse struct {
	Msg    ErrorMessage `json:"message"`
	Status string       `json:"status"`
	Error  string       `json:"error"`
}

type ErrorResponse struct {
	BadRequestResponse
	InternalServerResponse
}
