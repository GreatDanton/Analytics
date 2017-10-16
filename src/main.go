package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/greatdanton/analytics/src/sessions"

	"github.com/go-zoo/bone"
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
	r := bone.New()

	// open database connection
	connection := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=disable", config.DbUser, config.DbPassword, config.DbName)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Fatal("Cannot open db connection:", err)
	}
	defer db.Close()
	global.DB = db

	// handle command line flags -env=[setup, test]
	setup.HandleCmdFlags()

	// load all websites data into memory, TODO: replace this function with REDIS db
	global.Websites, err = model.LoadWebsitesToMemory()
	if err != nil {
		log.Fatal("Cannot load websites to memory:", err)
	}
	fmt.Println("Websites in memory:", global.Websites)

	r.HandleFunc("/403", controller.Handle403)
	// app handlers
	r.HandleFunc("/", loggedInUser(controller.MainHandler))
	r.HandleFunc("/l/:ID", controller.LogTraffic)
	r.HandleFunc("/login", controller.Login)
	r.HandleFunc("/logout", controller.Logout)
	r.HandleFunc("/dashboard", loggedInUser(controller.Dashboard))
	r.HandleFunc("/dashboard/new", loggedInUser(controller.AddWebsite))
	r.HandleFunc("/website/:name", loggedInUser(controller.WebsiteTraffic))
	r.HandleFunc("/website/:name/edit", loggedInUser(controller.EditWebsite))
	r.HandleFunc("/website/:name/delete", loggedInUser(controller.DeleteWebsite))

	// server public files
	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// for now /register part is only accessible if the application
	// is ran with -env=setup flag.
	// This is just a quick temporary solution for the testing phase
	// which will soon be replaced with a more robust one
	if global.RegisterAdmin {
		http.HandleFunc("/register/", controller.Register)
	}

	// start server
	if err := http.ListenAndServe(":"+config.Port, r); err != nil {
		log.Fatal(err)
	}
}

// loggedInUser is middleware that checks if the user is logged in
// and redirects him to the /login page if he is not
func loggedInUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := sessions.LoggedInUser(r)
		// check if user is logged in
		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// user is logged in: continue with request
		next.ServeHTTP(w, r)
	})
}
