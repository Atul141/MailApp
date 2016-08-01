package handlers

import (
	"encoding/json"
	"fmt"
	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
	"log"
	"net/http"
	"strings"
)

const (
	closed = "close"
	open   = "open"
)

func parcelsHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var parcels []*m.ParcelUserDetails
		var err error
		status := strings.Split(r.URL.Path, "/")[2]

		if status == closed {
			parcels, err = db.GetCloseParcels()
		} else if status == open {
			parcels, err = db.GetOpenParcels()
		}

		if err != nil {
			log.Printf("Error fetching parcels from DB: %s", err)

			errResponse := m.Error{
				Code:    u.I32Ptr(http.StatusInternalServerError),
				Message: u.SPtr("Internal server error"),
			}
			marshalledError, err := json.Marshal(errResponse)
			if err != nil {
				log.Printf("failed to marshal the error response: %s", err)
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
			http.Error(w, string(marshalledError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		marshalledRes, err := json.Marshal(parcels)
		if err != nil {
			log.Printf("failed to marshal the success response: %s", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}

		fmt.Fprintf(w, string(marshalledRes))
	}
}
