package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CreatePartyResponse struct {
	Id int `json:"id"`
}
