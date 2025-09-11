package database

type Party struct {
	ID int `db:"id"`
	Name string `db:"name"`
}

type Member struct {
	ID		int `db:"id"`
	Name    string `db:"name"`
	Balance string `db:"balance"`
}

type Payment struct {
	ID          int      `db:"id"`
	Description *string  `db:"description"`
	Amount      float32  `db:"amount"`
	Payer       string   `db:"payer_name"`
	Payees      []string `db:"payees"`
}
