package global

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
)

// This package is used to hold all global variables

// DB holds the connection to the database that is set up
// on application startup
var DB *sql.DB

// Configuration struct is used to hold the data parsed
// from config.json file
type Configuration struct {
	Port       string
	DbUser     string
	DbPassword string
	DbName     string
}

// ReadConfig reads configuration from config.json and returns
// Configuration file filled with data
func ReadConfig() Configuration {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Please add config.json file: ", err)
	}
	config := Configuration{}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal("Please format configuration file correctly:", err)
	}
	return config
}

// Website holds data about each website parsed from database
type Website struct {
	ID         string
	WebsiteURL string
}

// Websites var holds data for each website in database
var Websites map[string]Website
