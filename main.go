package main

import (
	requestHandler "htmx-demo/router"
	"log"      // self explanatory
	"net/http" // This provides HTTP client and server implementations for the app
	"os"       // Access operating system functionality

	"github.com/gorilla/mux"   // router for the site
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

	router := mux.NewRouter()
	router.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/", http.FileServer(http.Dir("/Users/chasedevries/Desktop/DesktopStuff/HtmxHelloWorld/HtmxDemo/styles/"))))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("/Users/chasedevries/Desktop/DesktopStuff/HtmxHelloWorld/HtmxDemo/assets/"))))
	requestHandler.HandleRequests(router, port)
}
