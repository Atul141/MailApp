package main

import (
	"flag"
	"log"

	"git.mailbox.com/mailbox/config"
	h "git.mailbox.com/mailbox/handlers"

	"github.com/codegangsta/negroni"
)

var configFilePath *string

func init() {
	configFilePath = flag.String("config", "default.conf", "Full path to application config file")
}

func main() {
	flag.Parse()

	config, err := config.ReadConfig(*configFilePath)
	if err != nil {
		log.Fatalf("error occured while reading file: %s", err)
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("error in validating config file: %s", err)
	}

	router := h.Router()
	
	n := negroni.New()
	n.UseHandler(router)

	appPort := config.GetServerPort()
	log.Printf("Starting service on http://localhost%s", appPort)
	n.Run(appPort)
}
