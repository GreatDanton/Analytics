package setup

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/greatdanton/analytics/src/global"
	"github.com/greatdanton/analytics/src/model"
)

// Configuration struct is used to hold the data parsed
// from config.json file
type Configuration struct {
	Port             string
	DbUser           string
	DbPassword       string
	DbName           string
	JwtTokenPassword string
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
	global.JwtTokenPassword = config.JwtTokenPassword
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
		global.RegisterAdmin = true
	}
}

// ParseTemplates parses all template files from /templates directory and
// stores them in global Templates variable
func ParseTemplates() {
	templ := template.New("")
	err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err = templ.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}

		return err
	})

	if err != nil {
		log.Fatal("ParseTemplates:", err)
	}

	global.Templates = templ
}
