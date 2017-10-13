package setup

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/model"
)

// AppSetup sets up everything
// reads configuration file, establish connection with database,
// creates database based on command line flags,
// loads website table into memory -> for faster access
func AppSetup() {
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

	// handle cmd flags -> create production or testing database based
	// on env=(test, setup) variable
	HandleCmdFlags()

	// load all websites data into memory, replace this function with REDIS db
	global.Websites, err = model.LoadWebsitesToMemory()
	if err != nil {
		fmt.Println("LoadWebsitesToMemory:", err)
		return
	}
	fmt.Println(global.Websites)
}

// HandleCmdFlags handles command line flags of the application for creating
// new:
// - production database (empty)
// - testing database (filled with some data)
func HandleCmdFlags() {
	// parsing environement flag
	wordPtr := flag.String("env", "", "use test for setting testing environment or new when setting up application")
	flag.Parse()
	if *wordPtr == "test" {
		// setUp our database (for developers) -> remove old tables and setup new ones
		model.CreateTestDB()
	} else if *wordPtr == "setup" {
		// setUp our database for the first time -> without removing old tables
		model.FirstStart()
	}
}
