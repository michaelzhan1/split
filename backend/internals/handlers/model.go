package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (error *HttpError) Error() string {
	return error.Message
}

type Party struct {
	Name string `json:"name"`
}

type Member struct {
	Name string `json:"name"`
}
