package main // This code belongs to the 'main' package

import (
	"html/template" // injection safe html generation
	jokeUtil "htmx-demo/jokes"
	"log"      // self explanatory
	"net/http" // This provides HTTP client and server implementations for the app
	"os"       // Access operating system functionality

	"github.com/joho/godotenv" // read from a .env file for this application
)

func index(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("components/index.html"))
	tpl.Execute(w, nil)
}

func jokes(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("components/jokes.html"))
	tpl.Execute(w, nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("components/contact.html"))
	tpl.Execute(w, nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("components/about.html"))
	tpl.Execute(w, nil)
}

func generate(w http.ResponseWriter, r *http.Request) {
	p := jokeUtil.GetRandomJoke()
	tpl, _ := template.ParseFiles("components/joke.html")
	tpl.Execute(w, p)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/favicon-32x32.png")
}

func handleRequests(mux *http.ServeMux, port string) {
	mux.HandleFunc("/", index)
	mux.HandleFunc("/jokes", jokes)
	mux.HandleFunc("/contact", contact)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/generate", generate)
	mux.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fs := http.FileServer(http.Dir("assets"))
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	handleRequests(mux, port)
}
