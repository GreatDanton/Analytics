package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/greatdanton/analytics/controller"
	"github.com/greatdanton/analytics/global"
	"github.com/greatdanton/analytics/model"
	_ "github.com/lib/pq"
)

func main() {
	config := global.ReadConfig()
	fmt.Printf("Running server on: http://127.0.0.1:%v\n", config.Port)

	// open database connection
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("Cannot open db connection:", err)
	}
	defer db.Close()
	global.DB = db

	// create new db
	//model.CreateDB()
	model.CreateTestDB()

	// load all websites data into memory, replace this function with REDIS db
	global.Websites, err = model.LoadWebsitesToMemory()
	if err != nil {
		fmt.Println("LoadWebsitesToMemory:", err)
		return
	}
	fmt.Println(global.Websites)

	// handlers
	http.HandleFunc("/", controller.MainHandler)
	// start server
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal(err)
	}
}
