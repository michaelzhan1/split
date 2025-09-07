package database

type Party struct {
	Name string `db:"name"`
}

type Member struct {
	Name string `db:"name"`
}
