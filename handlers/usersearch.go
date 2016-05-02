package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "git.mailbox.com/mailbox/models"
	u "git.mailbox.com/mailbox/utils"
)

func userSearchHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchParam := r.URL.Query().Get("q")

		if len(searchParam) < 3 {
			marshalledRes, err := json.Marshal([]*m.User{})
			if err != nil {
				log.Printf("failed to marshal the success response: %s", err)
				http.Error(w, "something went wrong", http.StatusInternalServerError)
			}
			fmt.Fprintf(w, string(marshalledRes))
			return
		}

		users, err := db.GetUsersWith(searchParam)
		if err != nil {
			log.Printf("Error fetching users from DB: %s", err)

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

		w.Header().Set("Content-Type", "application/json")

		marshalledRes, err := json.Marshal(users)
		if err != nil {
			log.Printf("failed to marshal the success response: %s", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, string(marshalledRes))
	}
}
