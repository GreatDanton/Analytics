package main

import (
	"log"
	"net/http"

	"github.com/greatdanton/analytics/src/controller"
	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/setup"
	_ "github.com/lib/pq"
)

func main() {
	// set everything up -> read config, connect to db, set up db
	setup.AppSetup()

	// app handlers
	http.HandleFunc("/", controller.MainHandler)

	// start server
	if err := http.ListenAndServe(":"+global.Config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
