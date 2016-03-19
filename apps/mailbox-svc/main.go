package main

import (
	"log"

	h "git.mailbox.com/mailbox/handlers"

	"github.com/codegangsta/negroni"
)

func main() {

	router := h.Router()

	n := negroni.New()
	n.UseHandler(router)

	appPort := ":8080"
	log.Printf("Starting service on http://localhost%s", appPort)
	n.Run(appPort)
}
