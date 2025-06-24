package requestHandler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	supabaseHandler "htmx-demo/supabase"
	util "htmx-demo/util"

	"github.com/supabase-community/gotrue-go/types"
)

type NavbarProps struct {
	Referer    string `json:"referer"`
	IsLoggedIn bool   `json:"is_logged_in"`
}

type Error struct {
	ErrorMessage string `json:"error_message"`
}

// GetNavbarForRequest extracts the referer from the request and returns a Navbar struct.
func GetNavbarForRequest(r *http.Request) NavbarProps {
	Referer := r.Referer()
	Host := r.Host
	_, err := util.GetCookieFromRequest(r, "SESSION")
	IsLoggedIn := err == nil

	path := strings.Split(Referer, Host)
	return NavbarProps{
		IsLoggedIn: IsLoggedIn,
		Referer:    path[len(path)-1],
	}
}

func Navbar(w http.ResponseWriter, r *http.Request) {
	n := GetNavbarForRequest(r)
	tpl := template.Must(template.ParseFiles("components/navbar.html"))
	err := tpl.Execute(w, n)
	if err != nil {
		log.Println("Error executing navbar template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// This function renders the budget if auth is found, otherwise it renders the login page.
func BudgetData(w http.ResponseWriter, r *http.Request) {
	cookie, err := util.GetCookieFromRequest(r, "SESSION")

	if err != nil {
		log.Println("No SESSION cookie found:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	client := supabaseHandler.GetSupabaseClient()
	client.UpdateAuthSession(types.Session{AccessToken: cookie.Value})
	data, _, err := client.From("transactions").Select("*", "exact", false).Execute()
	transactions := supabaseHandler.TransactionResponse{}
	err = json.Unmarshal(data, &transactions)
	if err != nil {
		log.Println("Error parsing supabase data:", err)
	}

	tpl := template.Must(template.ParseFiles("components/budget.html"))
	tpl.Execute(w, transactions)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("components/login.html"))
	err := tpl.Execute(w, nil)
	if err != nil {
		log.Println("Error executing login template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	w, r, err := supabaseHandler.Authenticate(w, r)
	if err != nil {
		LoginForm(w, r)
		return
	}

	http.Redirect(w, r, "/budget", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "SESSION",
		Value:  "",
		MaxAge: -1, // This will delete the cookie
	}
	http.SetCookie(w, &cookie)
	log.Println("Session cookie cleared")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
