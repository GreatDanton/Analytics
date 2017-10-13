package model

import (
	"fmt"

	"github.com/greatdanton/analytics/src/global"
)

// CreateDB creates new empty database and drops the old tables, when setting up
//analytics application for the first time.
func CreateDB() {
	// USERS TABLE
	_, err := global.DB.Exec(`drop table if exists users cascade`)
	if err != nil {
		fmt.Println("Could not drop table users:", err)
		return
	}
	_, err = global.DB.Exec(`CREATE TABLE users(id serial primary key,
												username varchar(25) unique NOT NULL,
												email text unique NOT NULL,
												password varchar(60) NOT NULL);`)

	// WEBSITE TABLE
	_, err = global.DB.Exec(`drop table if exists website cascade`)
	if err != nil {
		fmt.Println("Could not drop table website:", err)
		return
	}

	_, err = global.DB.Exec(`CREATE TABLE website(id serial primary key,
												  short_url varchar(8) UNIQUE,
												  owner integer references users(id) on delete cascade,
												  active boolean,
												  website_url text)`)
	if err != nil {
		fmt.Println("Could not create table website:", err)
		return
	}

	// CLICK TABLE
	_, err = global.DB.Exec(`drop table if exists click cascade`)
	if err != nil {
		fmt.Println("Could not drop table click:", err)
		return
	}
	_, err = global.DB.Exec(`CREATE TABLE click(id serial primary key,
												website_id integer references website(id) on delete cascade,
												time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
												ip varchar(16),
												url_clicked text);`)
	if err != nil {
		fmt.Println("Could not create table click:", err)
		return
	}

	// LAND TABLE
	_, err = global.DB.Exec(`drop table if exists land cascade`)
	if err != nil {
		fmt.Println("Could not drop table land:", err)
		return
	}

	_, err = global.DB.Exec(`CREATE TABLE land(id serial primary key,
											   website_id integer references website(id) on delete cascade,
											   time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
											   ip varchar(16));`)
	if err != nil {
		fmt.Println("Could not create table land:", err)
		return
	}

	fmt.Println("Database successfully set up.")
}
