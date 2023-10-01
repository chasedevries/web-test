package main // This code belongs to the 'main' package

import (
	"html/template" // injection safe html generation
	"log"           // self explanatory
	"math/rand"
	"net/http" // This provides HTTP client and server implementations for the app
	"os"       // Access operating system functionality

	"github.com/joho/godotenv" // read from a .env file for this application
)

type Joke struct {
	Noun      string `json:"Noun"`
	Verb      string `json:"Verb"`
	Adjective string `json:"Adjective"`
}

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
	nouns := []string{"horse", "eagle", "frog", "turkey", "cat"}
	verbs := []string{"gallops", "flies", "jumps", "trots", "sneaks"}
	adjectives := []string{"long", "pointy", "green", "wrinkly", "whiskery"}

	p := Joke{
		Noun:      nouns[rand.Intn(len(nouns))],
		Verb:      verbs[rand.Intn(len(verbs))],
		Adjective: adjectives[rand.Intn(len(adjectives))],
	}
	tpl, _ := template.ParseFiles("components/joke.html")
	tpl.Execute(w, p)
}

func handleRequests(mux *http.ServeMux, port string) {

	mux.HandleFunc("/", index)
	mux.HandleFunc("/generate", generate)
	mux.HandleFunc("/jokes", jokes)
	mux.HandleFunc("/contact", contact)
	mux.HandleFunc("/about", about)
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
