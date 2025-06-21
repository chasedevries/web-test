package main

import (
	requestHandler "htmx-demo/router"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // driver for sql (necessary when using go)
	"github.com/joho/godotenv"         // read from a .env file for this application
	"github.com/labstack/echo/v4"      // web framework for Go (Golang
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	// e.GET("/")
	log.Println("Starting Golang server...")
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
