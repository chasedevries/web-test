package main

import (
	requestHandler "htmx-demo/router"
	"log" // self explanatory
	"net/http"
	"time"

	// This provides HTTP client and server implementations for the app
	"os" // Access operating system functionality

	"database/sql" // interact with a sql database

	_ "github.com/go-sql-driver/mysql" // driver for sql (necessary when using go)
	"github.com/joho/godotenv"         // read from a .env file for this application
)

func main() {
	log.Println("Hey! I'm over here! The logs are over here!")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	const (
		DB_HOST = "tcp(127.0.0.1:3306)"
		DB_NAME = "demo-db"
		DB_USER = "root"
		DB_PASS = "secret"
	)

	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"

	for {

		log.Print(".")

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Error: " + err.Error())
			return
		}

		err = db.Ping()
		if err != nil {
			log.Println("Error: " + err.Error())
			time.Sleep(2 * time.Second)
			continue
		} else {
			log.Println("Success!")
			break
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := http.NewServeMux()
	requestHandler.HandleRequests(router, port)
}
