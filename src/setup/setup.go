package setup

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/model"
)

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
	// parses template files and stores it into global variable
	ParseTemplates()
	// return configuration struct
	return config
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

// ParseTemplates parses template files from /template directory and
// stores them in global Templates variable
func ParseTemplates() {
	parsedTemplates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	templ := template.Must(parsedTemplates, err)
	global.Templates = templ
}
