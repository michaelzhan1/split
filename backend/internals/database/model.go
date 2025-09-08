package database

type Party struct {
	Name string `db:"name"`
}

type Member struct {
	Name string `db:"name"`
}

type Payment struct {
	Description *string  `db:"description"`
	Amount      float32  `db:"amount"`
	Payer       string   `db:"payer_name"`
	Payees      []string `db:"-"` // added later
}
