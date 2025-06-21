package requestHandler

import (
	"html/template" // injection safe html generation
	jokeFactory "htmx-demo/jokes"
	requestHandler "htmx-demo/requests"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("views/index.html"))
	tpl.Execute(w, nil)
}

func jokes(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("views/jokes.html"))
	tpl.Execute(w, nil)

}

func photos(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("views/photos.html"))
	tpl.Execute(w, nil)
}

func budget(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("views/budget.html"))
	tpl.Execute(w, nil)
}

func generate(w http.ResponseWriter, r *http.Request) {
	p := jokeFactory.GetRandomJoke()
	tpl, _ := template.ParseFiles("components/joke.html")
	tpl.Execute(w, p)
}

func navBar(w http.ResponseWriter, r *http.Request) {
	n := requestHandler.GetNavbarForRequest(r)
	var tpl = template.Must(template.ParseFiles("components/navbar.html"))
	tpl.Execute(w, n)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "assets/favicon-32x32.png")
}

/**
 * as it turns out, only capitalized identifiers are exported from a module.
**/
func HandleRequests(router *http.ServeMux, port string) {
	router.HandleFunc("/favicon.ico", faviconHandler)
	router.HandleFunc("/", index)
	router.HandleFunc("/photos", photos)
	router.HandleFunc("/jokes", jokes)
	router.HandleFunc("/budget", budget)

	router.HandleFunc("/generate", generate)
	router.HandleFunc("/navbar", navBar)
	router.HandleFunc("/budgetData", requestHandler.BudgetData)
	router.HandleFunc("/auth", requestHandler.Auth)
	router.HandleFunc("/logout", requestHandler.Logout)
	// router.HandleFunc("/jokes", jokes)
	// router.HandleFunc("/contact", contact)
	// router.HandleFunc("/about", about)
	// router.HandleFunc("/comment", comment)
	// router.HandleFunc("/sam", sam)

	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	router.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	log.Println("Server is running on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
