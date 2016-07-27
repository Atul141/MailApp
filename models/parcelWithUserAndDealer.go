package models

import (
	"github.com/go-openapi/strfmt"
	"time"
)

type ParcelWithUserAndDealer struct {
	/* registration no

	Required: true
	*/
	ID string `db:"id" json:"id"`

	/* dealer

	Required: true
	*/
	DealerID *string `db:"dealer_id" json:"-"`

	/* recieved date

	Required: true
	*/
	RecievedDate time.Time `db:"received_date" json:"recieved_date"`

	/* status

	Required: true
	*/
	Status bool `db:"status" json:"status"`

	/* owner

	Required: true
	*/
	OwnerID *string `db:"owner_id" json:"owner_id"`

	/* reciever
	 */
	RecieverID *string `db:"receiver_id" json:"receiver_id"`

	/* User's personal/official email address

	Required: true
	*/
	UserEmail strfmt.Email `db:"user_email" json:"user_email"`

	/* name
	 */
	UserName *string `db:"user_name" json:"user_name,omitempty"`

	/* emp id

	Required: true
	*/
	UserEmpID string `db:"user_emp_id" json:"user_emp_id"`

	/* User's personal/official phone number
	 */
	UserPhoneNo *string `db:"user_phone_no" json:"user_phone_no,omitempty"`

	/* name
	 */
	DealerName *string `db:"dealer_name" json:"dealer_name,omitempty"`

	/* icon
	 */
	DealerIcon *strfmt.URI `db:"dealer_icon" json:"dealer_icon,omitempty"`
}
