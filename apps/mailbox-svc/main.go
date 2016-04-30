package main

import (
	"flag"
	"log"

	"git.mailbox.com/mailbox/config"
	h "git.mailbox.com/mailbox/handlers"
	m "git.mailbox.com/mailbox/models"

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

	db, err := m.NewDatabase(config.DbConnString)
	if err != nil {
		log.Fatalf("fail to connect to database: %s", err)
	}
	defer db.Close()

	router := h.Router(db)
	n := negroni.New()
	n.UseHandler(router)

	appPort := config.GetServerPort()
	log.Printf("Starting service on http://localhost%s", appPort)
	n.Run(appPort)
}
