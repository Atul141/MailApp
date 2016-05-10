package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
	"time"
)

/*Parcel parcel

swagger:model Parcel
*/
type Parcel struct {

	/* dealer

	Required: true
	*/
	DealerID *string `db:"dealer_id" json:"-"`

	/* dealer

	Required: true
	*/
	Dealer *Dealer `json:"dealer"`

	/* owner

	Required: true
	*/
	OwnerID *string `db:"owner_id" json:"-"`
	/* owner

	Required: true
	*/
	Owner *User `json:"owner"`

	/* pickup date
	 */
	PickupDate time.Time `json:"pickup_date,omitempty"`

	/* recieved date

	Required: true
	*/
	RecievedDate time.Time `db:"received_date" json:"recieved_date"`

	/* reciever
	 */
	RecieverID *string `db:"receiver_id" json:"-"`
	/* reciever
	 */
	Reciever *User `json:"reciever,omitempty"`

	/* registration no

	Required: true
	*/
	ID string `db:"id" json:"id"`

	/* status

	Required: true
	*/
	Status bool `db:"status" json:"status"`

	CreatedOn time.Time `db:"created_on" json:"-"`
}

// Validate validates this parcel
func (m *Parcel) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDealer(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateOwner(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Parcel) validateDealer(formats strfmt.Registry) error {

	if m.Dealer != nil {

		if err := m.Dealer.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}

func (m *Parcel) validateID(formats strfmt.Registry) error {

	if err := validate.RequiredString("id", "body", string(m.ID)); err != nil {
		return err
	}

	return nil
}

func (m *Parcel) validateOwner(formats strfmt.Registry) error {

	if m.Owner != nil {

		if err := m.Owner.Validate(formats); err != nil {
			return err
		}
	}

	return nil
}

func (m *Parcel) validateStatus(formats strfmt.Registry) error {

	if err := validate.Required("status", "body", bool(m.Status)); err != nil {
		return err
	}

	return nil
}
