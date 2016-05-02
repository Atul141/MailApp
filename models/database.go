package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const getAllDealersQuery = "select id, name, icon from dealers;"

type DB interface {
	GetDealers() ([]*Dealer, error)
	GetUsersWith(string) ([]*User, error)
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

func (db *Database) GetUsersWith(searchParam string) ([]*User, error) {
	var users []*User
	searchUsersQuery := `select email, name, emp_id, phone_no from users where name ~* $1 or email ~* $1 or phone_no ~* $1`
	err := db.connection.Select(&users, searchUsersQuery, searchParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %s", err)
	}
	return users, err
}
