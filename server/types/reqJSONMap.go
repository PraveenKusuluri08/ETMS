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

type SuccessResponse struct {
	Message string `json:"message"`
}

type InviteGroupMembersResponse struct {
	Message                 string   `json:"message"`
	Non_existing_users      []string `json:"non_existing_users"`
	Total_no_existing_users int      `json:"total_no_existing_users"`
}

type AcceptInvitationResponse struct {
	Message string `json:"message"`
}

type GetUsersResponse struct {
	Users []map[string]string `json:"users"`
}
