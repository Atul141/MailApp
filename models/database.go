package models

import (
	"fmt"

	u "git.mailbox.com/mailbox/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"time"
)

const getAllDealersQuery = `select id, name, icon from dealers;`
const getDealerByIDQuery = `select id, name, icon from dealers WHERE id=$1;`

type DB interface {
	GetDealers() ([]*Dealer, error)
	GetUsersWith(string) ([]*User, error)
	GetUserByID(string) (*User, error)
	GetDealerByID(id string) (*Dealer, error)
	GetParcelByID(id string) (*Parcel, error)
	CreateParcel(dealerID string, ownerID string, comments string) (*Parcel, error)
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

func (db *Database) GetDealerByID(id string) (*Dealer, error) {
	dealer := Dealer{}
	err := db.connection.Get(&dealer, getDealerByIDQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Dealer: %s", err)
	}
	return &dealer, err
}

func (db *Database) GetUsersWith(searchParam string) ([]*User, error) {
	var users []*User
	searchUsersQuery := `select id, email, name, emp_id, phone_no from users where name ~* $1 or email ~* $1 or phone_no ~* $1`
	err := db.connection.Select(&users, searchUsersQuery, searchParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %s", err)
	}
	return users, err
}

func (db *Database) GetUserByID(id string) (*User, error) {
	user := User{}
	query := `SELECT id, email, name, emp_id, phone_no FROM users WHERE id=$1`
	err := db.connection.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %s", err)
	}
	return &user, err
}

func (db *Database) GetParcelByID(id string) (*Parcel, error) {
	parcel := Parcel{}
	query := `SELECT id, dealer_id, received_date, status,owner_id,receiver_id  FROM parcels WHERE id=$1`
	err := db.connection.Get(&parcel, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parcel: %s", err)
	}
	return &parcel, err
}
func (db *Database) CreateParcel(dealerID string, ownerID string, comments string) (*Parcel, error) {
	id := uuid.NewV4().String()
	parcel := &Parcel{
		ID:           id,
		DealerID:     u.SPtr(dealerID),
		OwnerID:      u.SPtr(ownerID),
		Status:       true,
		RecievedDate: time.Now().UTC(),
		CreatedOn:    time.Now().UTC(),
	}
	query := "INSERT INTO parcels (id,dealer_id,owner_id,status,received_date,created_on) VALUES (:id,:dealer_id,:owner_id,:status,:received_date,:created_on)"

	tx := db.connection.MustBegin()
	_, err := tx.NamedExec(query, &parcel)
	if err != nil {
		return nil, fmt.Errorf("Error inserting record: %s", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Error inserting record: %s", err)
	}
	return db.GetParcelByID(id)
}
