package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetPartyResponse struct {
	Name string `json:"name"`
}

type CreatePartyResponse struct {
	Id int `json:"id"`
}
