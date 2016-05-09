package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	me "github.com/rshetty/multierror"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

const UUIDRegex = "^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$"

type createParcelRequest struct {
	DealerID string `json:"dealerId"`
	OwnerID  string `json:"ownerId"`
}

func (cpr *createParcelRequest) validate() *me.MultiError {
	cprError := &me.MultiError{}

	if len(cpr.DealerID) == 0 {
		cprError.Push("dealerId should not be empty")
	}

	if !cpr.validUUIDV4(cpr.DealerID) {
		cprError.Push("dealerId should be a valid UUID V4 string")
	}

	if len(cpr.OwnerID) == 0 {
		cprError.Push("dealerId should not be empty")
	}

	if !cpr.validUUIDV4(cpr.OwnerID) {
		cprError.Push("dealerId should be a valid UUID V4 string")
	}

	return cprError.HasError()
}

func (cpr *createParcelRequest) validUUIDV4(text string) bool {
	r := regexp.MustCompile(UUIDRegex)
	return r.MatchString(text)
}

type createParcelResponse struct {
	ID string `json:"id"`
}

//parcelCreateHandler creates new parcel object
func parcelCreateHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			log.Printf("request body should not be empty")
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		respBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to read the request body: %s", err)
			http.Error(w, "request body error", http.StatusBadRequest)
			return
		}

		cpr := &createParcelRequest{}
		err = json.Unmarshal(respBody, cpr)
		if err != nil {
			log.Printf("failed to unmarshal the request body: %s", err)
			http.Error(w, "request body parsing failed", http.StatusBadRequest)
			return
		}

		if err := cpr.validate(); err != nil {
			log.Printf("Error fetching dealer from DB: %#v", err)
			badRequestError(w, err.Error())
			return
		}

		// Find the dealer by ID
		dealer, err := db.GetDealerByID(cpr.DealerID)
		if err != nil {
			log.Printf("Error fetching dealer from DB: %s", err)
			databaseError(w, err)
			return
		}

		// Find the owner by ID
		owner, err := db.GetUserByID(cpr.OwnerID)
		if err != nil {
			log.Printf("Error fetching owner from DB: %s", err)
			databaseError(w, err)
			return
		}

		// Create a parcel with owner and dealer
		parcel, err := db.CreateParcel(cpr.DealerID, cpr.OwnerID)
		if err != nil {
			log.Printf("Error creating a parcel: %s", err)
			databaseError(w, err)
			return
		}

		parcel.Dealer = dealer
		parcel.Owner = owner

		w.Header().Set("Content-Type", "application/json")

		createParcelResponse := &createParcelResponse{
			ID: parcel.ID,
		}

		marshalledRes, err := json.Marshal(createParcelResponse)
		if err != nil {
			log.Printf("failed to marshal the success response: %s", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, string(marshalledRes))
	}
}

func badRequestError(w http.ResponseWriter, msg string) {
	errResponse := m.Error{
		Code:    u.I32Ptr(http.StatusBadRequest),
		Message: u.SPtr(msg),
	}

	marshalledError, err := json.Marshal(errResponse)
	if err != nil {
		unexpectedError(w, err)
		return
	}

	http.Error(w, string(marshalledError), http.StatusBadRequest)
	return
}

func databaseError(w http.ResponseWriter, err error) {
	errResponse := m.Error{
		Code:    u.I32Ptr(http.StatusInternalServerError),
		Message: u.SPtr("Internal server error"),
	}
	marshalledError, err := json.Marshal(errResponse)
	if err != nil {
		unexpectedError(w, err)
		return
	}
	http.Error(w, string(marshalledError), http.StatusInternalServerError)
	return
}

func unexpectedError(w http.ResponseWriter, err error) {
	log.Printf("failed to marshal the error response: %s", err)
	http.Error(w, "something went wrong", http.StatusInternalServerError)
}
