package handlers

import (
	"encoding/json"
	m "git.mailbox.com/mailbox/models"
	"log"
	"net/http"
	"io/ioutil"
	"strings"
)

type updateParcelRequest struct {
	Status string `json:"status"`
}

func updateParcelHandler(db m.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			log.Printf("request body should not be empty")
			http.Error(w, "request body is empty", http.StatusBadRequest)
			return
		}

		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Printf("failed to read the request body: %s", err)
			http.Error(w, "request body error", http.StatusBadRequest)
			return
		}
		upr := &updateParcelRequest{}
		err = json.Unmarshal(reqBody, upr)
		if err != nil {
			log.Printf("failed to unmarshal the request body: %s", err)
			http.Error(w, "request body parsing failed", http.StatusBadRequest)
			return
		}

		parcelId := strings.Split(r.URL.Path, "/")[2]

		err = db.UpdateParcelStatusById(parcelId,upr.Status)

		if err != nil {
			log.Printf("Failed to update parcel: %s", err)
			databaseError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

