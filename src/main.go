package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/greatdanton/analytics/src/controller"
	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/model"
	"github.com/greatdanton/analytics/src/setup"
	_ "github.com/lib/pq"
)

func main() {
	// set everything up -> read config, connect to db, set up db
	config := setup.ReadConfig()
	fmt.Printf("Running server on: http://127.0.0.1:%v\n", config.Port)

	// open database connection (this should be in main file)
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("Cannot open db connection:", err)
	}
	defer db.Close()
	global.DB = db

	setup.HandleCmdFlags()

	// load all websites data into memory, replace this function with REDIS db
	global.Websites, err = model.LoadWebsitesToMemory()
	if err != nil {
		log.Fatal("Cannot load websites to memory:", err)
	}
	fmt.Println(global.Websites)

	// app handlers
	http.HandleFunc("/", controller.MainHandler)

	// start server
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
