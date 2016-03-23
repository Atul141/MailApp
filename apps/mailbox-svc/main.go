package main

import (
	"flag"
	"log"

	"git.mailbox.com/mailbox/config"
	h "git.mailbox.com/mailbox/handlers"

	"github.com/codegangsta/negroni"
)

var (
	configFilePath = flag.String("config", "default.conf", "Full path to application config file")
	appName        = "mainbox-service"
)

func main() {
	flag.Parse()
	appconfig, err := config.ReadApplicationConfig(*configFilePath)
	if err != nil {
		log.Fatalf("Error occured while reading file: %v", err)
	}

	if err := appconfig.Validate(); err != nil {
		log.Fatalf("Config validation failure: %v", err)
	}
	router := h.Router()

	n := negroni.New()
	n.UseHandler(router)

	appPort := appconfig.GetServerPort()
	log.Printf("Starting service on http://localhost%s", appPort)
	n.Run(appPort)
}
