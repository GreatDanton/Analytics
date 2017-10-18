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

// Templates is storing all template .html files
var Templates *template.Template

// JwtTokenPassword holds password string that is used to sign
// signature of the JWT (that way users can not forge their jwt)
var JwtTokenPassword string
