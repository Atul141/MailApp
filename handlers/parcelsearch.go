package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

func parcelSearchHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchParam := r.URL.Query().Get("q")

		if len(searchParam) < 3 {
			marshalledRes, err := json.Marshal([]*m.Parcel{})
			if err != nil {
				log.Printf("failed to marshal the success response: %s", err)
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
			fmt.Fprintf(w, string(marshalledRes))
			return
		}

		parcels, err := db.GetParcelsWith(searchParam)
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
				return
			}
			http.Error(w, string(marshalledError), http.StatusInternalServerError)
			return
		}

		for _, p := range parcels {
			dealer, err := db.GetDealerByID(*p.DealerID)
			if err != nil {
				log.Printf("Error reading dealer: %s", err)
				databaseError(w, err)
				return
			}
			p.Dealer = dealer

			owner, err := db.GetUserByID(*p.OwnerID)
			if err != nil {
				log.Printf("Error reading a owner: %s", err)
				databaseError(w, err)
				return
			}
			p.Owner = owner

			if p.RecieverID != nil {
				receiver, err := db.GetUserByID(*p.RecieverID)
				if err != nil {
					log.Printf("Error reading a receiver: %s", err)
					databaseError(w, err)
					return
				}
				p.Reciever = receiver
			}
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
