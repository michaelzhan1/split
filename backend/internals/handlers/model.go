package handlers

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (error *HttpError) Error() string {
	return error.Message
}

type Party struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Member struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Balance float32 `json:"balance"`
}

type Payment struct {
	ID          int      `json:"id"`
	Description *string  `json:"description"`
	Amount      float32  `json:"amount"`
	Payer       Member   `json:"payer"`
	Payees      []Member `json:"payees"`
}

type IOU struct {
	FromID int     `json:"from"`
	ToID   int     `json:"to"`
	Amount float32 `json:"amount"`
}
