package handlers

import (
	"fmt"
	"net/http"

	"encoding/json"

	"log"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

func dealersHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dealers, err := db.GetDealers()
		if err != nil {
			log.Printf("Error fetching dealers from DB: %s", err)

			errResponse := m.Error{
				Code:    u.I32Ptr(http.StatusInternalServerError),
				Message: u.SPtr("Internal server error"),
			}
			marshalledError, _ := json.Marshal(errResponse)
			http.Error(w, string(marshalledError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		marshalledRes, _ := json.Marshal(dealers)
		fmt.Fprintf(w, string(marshalledRes))
	}
}
