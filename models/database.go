package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const getAllDealersQuery = "select id, name, icon from dealers;"

type DB interface {
	GetDealers() ([]*Dealer, error)
}

type Database struct {
	connection *sqlx.DB
}

func NewDatabase(connString string) (*Database, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %s", err)
	}

	return &Database{connection: db}, nil
}

func (db *Database) Close() {
	db.connection.Close()
}

func (db *Database) GetDealers() ([]*Dealer, error) {
	var dealers []*Dealer
	err := db.connection.Select(&dealers, getAllDealersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dealers: %s", err)
	}
	return dealers, err
}
