package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/greatdanton/analytics/src/memory"
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

	m := memory.MemWebsites{}
	err = m.LoadWebsites()
	if err != nil {
		log.Fatal(err)
	}
	memory.Memory = m
	fmt.Println(m)

	r.HandleFunc("/403", controller.Handle403)
	// app handlers
	r.HandleFunc("/", loggedInUser(controller.MainHandler))
	r.HandleFunc("/l/:ID", controller.LogTraffic)
	r.HandleFunc("/login", controller.Login)
	r.HandleFunc("/logout", controller.Logout)
	r.HandleFunc("/dashboard", loggedInUser(controller.Dashboard))
	r.HandleFunc("/dashboard/new", loggedInUser(controller.AddWebsite))
	r.HandleFunc("/website/:id", loggedInUser(controller.WebsiteTraffic))
	r.HandleFunc("/website/:id/edit", loggedInUser(controller.EditWebsite))
	r.HandleFunc("/website/:id/delete", loggedInUser(controller.DeleteWebsite))

	// server public files
	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// for now /register part is only accessible if the application
	// is ran with -env=setup flag.
	// This is just a quick temporary solution for the testing phase
	// which will soon be replaced with a more robust one
	if global.RegisterAdmin {
		r.HandleFunc("/register", controller.Register)
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

func websiteOwner(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := sessions.LoggedInUser(r)
		// check if user is logged in
		if !user.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// check if user is owner of the website
		websiteID := strings.Split(r.RequestURI, "/")[2]
		website, err := model.GetWebsiteDetail(websiteID, user.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("Logged in user is not the owner of the website")
				http.Redirect(w, r, "/403", http.StatusForbidden)
				return
			}
			fmt.Println("websiteOwner error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// everything is okay, pass website and user forward
		// user is logged in: continue with request
		fmt.Println(website)
		next.ServeHTTP(w, r)
	})
}
