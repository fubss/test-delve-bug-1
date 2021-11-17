package main

import (
	db "github.com/fubss/test-delve-bug-1/db"
)

func main() {
	dbPath := "./db"

	DB := db.CreateDB(&db.Conf{
		DBPath:         dbPath,
		ExpectedFormat: db.CurrentFormat,
	})

	DB.Open()

}
