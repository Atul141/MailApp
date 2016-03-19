package handlers

import (
	"net/http"

	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Router Generates a gorilla mux router with routes
func Router() http.Handler {

	r := mux.NewRouter()
	r.Handle("/", fileHandler("public/index.html")).Methods("GET")

	return gh.CompressHandler(r)
}
