package dto

type HttpSuccess struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HttpUserSuccess struct {
	Status string      `json:"status"`
	User   UserViewDto `json:"user"`
}

type HttpUsersSuccess struct {
	Status string         `json:"status"`
	Users  []*UserViewDto `json:"users"`
}

type HttpFailure struct {
	Status string `json:"status"`
	Error  error  `json:"error"`
}

type HttpFailures struct {
	Status string  `json:"status"`
	Errors []error `json:"errors"`
}

type HttpAccessTokenSuccess struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}
