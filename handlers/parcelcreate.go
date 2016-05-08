package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"fmt"

	"strings"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

//parcelCreateHandler creates new parcel object
func parcelCreateHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dealerID := r.PostFormValue("dealer_id")
		if !validateRequestParam(w, dealerID, "dealer_id") {
			return
		}

		ownerID := r.PostFormValue("owner_id")
		if !validateRequestParam(w, ownerID, "owner_id") {
			return
		}

		dealer, err := db.GetDealerByID(dealerID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				log.Printf("Dealer with id `%s` not found.", dealerID)
				notFoundError(w, "dealer_id")
				return
			}
			log.Printf("Error fetching dealer from DB: %s", err)
			databaseError(w, err)
			return
		}

		owner, err := db.GetUserByID(ownerID)
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				log.Printf("Owner with id `%s` not found.", dealerID)
				notFoundError(w, "owner_id")
				return
			}
			log.Printf("Error fetching owner from DB: %s", err)
			databaseError(w, err)
			return
		}

		parcel, err := db.CreateParcel(dealerID, ownerID, r.PostFormValue("comments"))
		if err != nil {
			log.Printf("Error creating parcel: %s", err)
			databaseError(w, err)
			return
		}

		parcel.Dealer = dealer
		parcel.Owner = owner

		w.Header().Set("Content-Type", "application/json")

		marshalledRes, err := json.Marshal(parcel)
		if err != nil {
			log.Printf("failed to marshal the success response: %s", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, string(marshalledRes))
	}
}

func validateRequestParam(w http.ResponseWriter, value string, field string) bool {
	if len(value) == 0 {
		log.Printf("Bad Request. `%s` request param is empty.", field)
		badRequestError(w, field, fmt.Sprintf("Bad Request. `%s` field can not be empty.", field))
		return false
	}

	if !validateUuidV4(value) {
		log.Printf("Bad Request. `%s` request param is invalid uuid v4.", field)
		badRequestError(w, field, fmt.Sprintf("Bad Request. `%s` field is invalid.", field))
		return false
	}
	return true
}

func badRequestError(w http.ResponseWriter, field string, msg string) {
	errResponse := m.Error{
		Code:    u.I32Ptr(http.StatusBadRequest),
		Message: u.SPtr(msg),
		Fields:  u.SPtr(field),
	}
	marshalledError, err := json.Marshal(errResponse)
	if err != nil {
		unexpectedError(w, err)
		return
	}
	http.Error(w, string(marshalledError), http.StatusBadRequest)
	return
}

func notFoundError(w http.ResponseWriter, field string) {
	errResponse := m.Error{
		Code:    u.I32Ptr(http.StatusNotFound),
		Message: u.SPtr(fmt.Sprintf("%s not found", field)),
		Fields:  u.SPtr(field),
	}
	marshalledError, err := json.Marshal(errResponse)
	if err != nil {
		unexpectedError(w, err)
		return
	}
	http.Error(w, string(marshalledError), http.StatusNotFound)
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

func validateUuidV4(text string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	return r.MatchString(text)
}
