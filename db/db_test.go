package db

import (
	"testing"
)

func TestOpening(t *testing.T) {
	dbPath := "./db"

	db := CreateDB(&Conf{
		DBPath:         dbPath,
		ExpectedFormat: CurrentFormat,
	})

	db.Open()
}
