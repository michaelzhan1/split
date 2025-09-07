package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Party struct {
	Name string `json:"name"`
}

type Member struct {
	Name string `json:"name"`
}
