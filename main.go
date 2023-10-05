package main

import (
	requestHandler "htmx-demo/router"
	"log" // self explanatory
	"net/http"

	// This provides HTTP client and server implementations for the app
	"os" // Access operating system functionality

	"github.com/joho/godotenv" // read from a .env file for this application
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router := http.NewServeMux()
	requestHandler.HandleRequests(router, port)
}
