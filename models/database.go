package models

import (
	"fmt"

	u "git.mailbox.com/mailbox/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"time"
)

const (
	getAllDealersQuery   = `select id, name, icon from dealers;`
	getDealerByIDQuery   = `select id, name, icon from dealers WHERE id=$1;`
	getCloseParcelsQuery = `select p.id, p.dealer_id, p.received_date, p.status, p.owner_id, p.receiver_id,
 			     u.email as user_email, u.name as user_name, u.emp_id as user_emp_id, u.phone_no as user_phone_no,
 			     d.name as dealer_name, d.icon as dealer_icon
 			     from parcels as p, dealers as d,users as u where p.dealer_id = d.id AND p.owner_id = u.id AND p.status = false;`
	getOpenParcelsQuery = `select p.id, p.dealer_id, p.received_date, p.status, p.owner_id, p.receiver_id,
 			     u.email as user_email, u.name as user_name, u.emp_id as user_emp_id, u.phone_no as user_phone_no,
 			     d.name as dealer_name, d.icon as dealer_icon
 			     from parcels as p, dealers as d,users as u where p.dealer_id = d.id AND p.owner_id = u.id AND p.status = true;`
)

var parcelStatus map[string]bool = map[string]bool{
	"closed": false,
	"open":   true,
}

type DB interface {
	GetDealers() ([]*Dealer, error)
	GetUsersWith(string) ([]*User, error)
	GetUserByID(string) (*User, error)
	GetDealerByID(id string) (*Dealer, error)
	GetParcelByID(id string) (*Parcel, error)
	GetParcelsWith(searchParam string) ([]*Parcel, error)
	CreateParcel(dealerID string, ownerID string) (*Parcel, error)
	GetCloseParcels() ([]*ParcelUserDetails, error)
	GetOpenParcels() ([]*ParcelUserDetails, error)
	UpdateParcelStatusById(parcelId string, status string) error
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
	query := `SELECT id, dealer_id, received_date, status, owner_id, receiver_id FROM parcels WHERE id=$1`
	err := db.connection.Get(&parcel, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parcel: %s", err)
	}
	return &parcel, err
}

func (db *Database) GetParcelsWith(searchParam string) ([]*Parcel, error) {
	var parcels []*Parcel
	query := `SELECT p.id, p.dealer_id, p.received_date, p.status, p.owner_id, p.receiver_id FROM parcels p
				inner join dealers d on d.id = p.dealer_id
				inner join users o on o.id = p.owner_id
				where d.name ~*$1 or
				o.name ~*$1 or o.email ~*$1 or o.phone_no ~* $1 or o.emp_id ~*$1;`
	err := db.connection.Select(&parcels, query, searchParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parcels: %s", err)
	}
	return parcels, err
}

func (db *Database) CreateParcel(dealerID string, ownerID string) (*Parcel, error) {
	id := uuid.NewV4().String()
	parcel := &Parcel{
		ID:           id,
		DealerID:     u.SPtr(dealerID),
		OwnerID:      u.SPtr(ownerID),
		Status:       true,
		RecievedDate: time.Now().UTC(),
		CreatedOn:    time.Now().UTC(),
	}
	query := "INSERT INTO parcels (id, dealer_id, owner_id, status, received_date, created_on) VALUES (:id, :dealer_id, :owner_id, :status, :received_date, :created_on)"

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

func (db *Database) GetCloseParcels() ([]*ParcelUserDetails, error) {
	var parcelDetails []*ParcelUserDetails
	err := db.connection.Select(&parcelDetails, getCloseParcelsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dealers: %s", err)
	}
	return parcelDetails, err
}

func (db *Database) GetOpenParcels() ([]*ParcelUserDetails, error) {
	var parcelDetails []*ParcelUserDetails
	err := db.connection.Select(&parcelDetails, getOpenParcelsQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dealers: %s", err)
	}
	return parcelDetails, err
}

func (db *Database) UpdateParcelStatusById(parcelId string, status string) error {
	query := "UPDATE parcels SET status=$1 WHERE id = $2;"

	_, err := db.connection.Exec(query, parcelStatus[status], parcelId)
	return err
}
