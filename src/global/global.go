package global

import (
	"database/sql"
	"html/template"
)

// This package is used to hold all global variables

// DB holds the connection to the database that is set up
// on application startup
var DB *sql.DB

// RegisterAdmin takes care of displaying register page upon
// web application startup,
// true => display register page
// false => do not display register page
var RegisterAdmin bool

// Website holds data about each website parsed from database
type Website struct {
	ID         string
	WebsiteURL string
}

// Websites var holds data for each website in database
var Websites map[string]Website

// Templates is storing all template .html files
var Templates *template.Template

// JwtTokenPassword holds password string that is used to sign
// signature of the JWT (that way users can not forge their jwt)
var JwtTokenPassword string
