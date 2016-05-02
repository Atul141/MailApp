package handlers

import (
	"net/http"

	m "git.mailbox.com/mailbox/models"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Router Generates a gorilla mux router with routes
func Router(db m.DB) http.Handler {

	r := mux.NewRouter()
	r.Handle("/", fileHandler("public/index.html")).Methods("GET")
	r.Handle("/users/search", userSearchHandler(db)).Methods("GET")
	r.Handle("/dealers", dealersHandler(db)).Methods("GET")
	r.Handle("/parcels/{status}", parcelSearchHandler()).Methods("GET")

	return gh.CompressHandler(r)
}
