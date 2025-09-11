package database

type Party struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Member struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Balance float32 `db:"balance"`
}

type Payment struct {
	ID            int       `db:"id"`
	Description   *string   `db:"description"`
	Amount        float32   `db:"amount"`
	PayerID       int       `db:"payer_id"`
	PayerName     string    `db:"payer_name"`
	PayerBalance  float32   `db:"payer_balance"`
	PayeeIDs      []int     `db:"payee_ids"`
	PayeeNames    []string  `db:"payee_names"`
	PayeeBalances []float32 `db:"payee_balances"`
}
