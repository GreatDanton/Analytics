package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/greatdanton/analytics/src/global"
)

// CreateDB creates new empty database and drops the old tables, when setting up
//analytics application for the first time.
func CreateDB() {
	// USERS TABLE
	err := dropUsers()
	if err != nil {
		fmt.Println("Could not drop table users:", err)
		return
	}

	err = createUsers()
	if err != nil {
		fmt.Println("Could not create users table:", err)
		return
	}
	// WEBSITE TABLE
	err = dropWebsite()
	if err != nil {
		fmt.Println("Could not drop table website:", err)
		return
	}

	err = createWebsite()
	if err != nil {
		fmt.Println("Could not create website:", err)
		return
	}

	// BROWSER table
	err = dropBrowser()
	if err != nil {
		fmt.Println("Could not drop table browser")
		return
	}
	err = createBrowser()
	if err != nil {
		fmt.Println("Could not create browser table")
		return
	}

	// CLICK TABLE
	err = dropClick()
	if err != nil {
		fmt.Println("Could not drop table click:", err)
		return
	}

	err = createClick()
	if err != nil {
		fmt.Println("Could not create table click:", err)
		return
	}

	// LAND TABLE
	err = dropLand()
	if err != nil {
		fmt.Println("Could not drop table land:", err)
		return
	}
	err = createLand()
	if err != nil {
		fmt.Println("Could not create table land:", err)
		return
	}

	fmt.Println("Database successfully set up.")
}

// drop users table
func dropUsers() error {
	_, err := global.DB.Exec(`drop table if exists users cascade`)
	if err != nil {
		return err
	}
	return nil
}

// create users table
func createUsers() error {
	_, err := global.DB.Exec(`CREATE TABLE users(id serial primary key,
												username varchar(25) unique NOT NULL,
												email text unique NOT NULL,
												password varchar(60) NOT NULL,
												active bool);`)
	if err != nil {
		return err
	}
	return nil
}

// drop website table
func dropWebsite() error {
	_, err := global.DB.Exec(`drop table if exists website cascade`)
	if err != nil {
		return err
	}
	return nil
}

// create website table
func createWebsite() error {
	_, err := global.DB.Exec(`CREATE TABLE website(id serial primary key,
												  short_url varchar(8) UNIQUE,
												  owner integer references users(id) on delete cascade,
												  active boolean,
												  website_url text)`)
	if err != nil {
		return err
	}
	return nil
}

// drop click table
func dropClick() error {
	// CLICK TABLE
	_, err := global.DB.Exec(`drop table if exists click cascade`)
	if err != nil {
		return err
	}
	return nil
}

// create click table
func createClick() error {
	_, err := global.DB.Exec(`CREATE TABLE click(id serial primary key,
												website_id integer references website(id) on delete cascade,
												time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
												browser_id integer references browser(id),
												ip varchar(16),
												url_clicked text);`)
	if err != nil {
		return err
	}
	return nil
}

// drop land table
func dropLand() error {
	_, err := global.DB.Exec(`drop table if exists land cascade`)
	if err != nil {
		return err
	}
	return nil
}

// create land table
func createLand() error {
	_, err := global.DB.Exec(`CREATE TABLE land(id serial primary key,
											   website_id integer references website(id) on delete cascade,
											   time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
											   browser_id integer references browser(id),
											   ip varchar(16));`)
	if err != nil {
		return err
	}
	return nil
}

// drop browser table
func dropBrowser() error {
	_, err := global.DB.Exec(`drop table if exists browser cascade`)
	if err != nil {
		return err
	}
	return nil
}

// create browser table
func createBrowser() error {
	_, err := global.DB.Exec(`CREATE TABLE browser(id serial primary key,
												   version text UNIQUE)`)
	if err != nil {
		return err
	}
	return nil
}

// FirstStart creates a database and drops all old tables
// use this when running application for the first time
func FirstStart() {
	CreateDB()
}

// CreateTestDB drops all old tables, creates new and inserts initial data
func CreateTestDB() {
	CreateDB() // drop old database and create new tables
	passHash, err := bcrypt.GenerateFromPassword([]byte("1234"), 10)
	if err != nil {
		fmt.Println("CreateTestDB: bcrypt:", err)
		return
	}
	_, err = global.DB.Exec(`INSERT into users(username, email, password, active)
							  values('user1', 'some@email.com', $1, TRUE);`, passHash)
	if err != nil {
		fmt.Println("Problem inserting data into users:", err)
		return
	}
	_, err = global.DB.Exec(`INSERT into website(short_url, owner, active, website_url)
							  values('12345678', 1, TRUE, 'http://www.somewebsite.com');`)
	if err != nil {
		fmt.Println("Problem inserting data into website:", err)
		return
	}
}
