package requestHandler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type Navbar struct {
	Referer string `json:"referer"`
}

type Transaction struct {
	Id                int    `json:"id"`
	CreatedAt         string `json:"created_at"`
	TransactionDate   string `json:"transaction_date"`
	TransactionAmount int    `json:"transaction_amount"`
}

type TransactionResponse []Transaction

// Creates a supabase client using the anonymous key
// TODO: move this to env or config file
func getSupabaseClient() *supabase.Client {
	client, err := supabase.NewClient("https://oyucpcsumehntdxplvgl.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im95dWNwY3N1bWVobnRkeHBsdmdsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTA0Njc0NjAsImV4cCI6MjA2NjA0MzQ2MH0.N9RIgBcDY6t4cpvohSz7h0XE9cctp6Dl1hhepxFsOrc", &supabase.ClientOptions{})
	if err != nil {
		log.Println("Error creating Supabase client:", err)
		return nil
	}
	return client
}

func getCookieFromRequest(r *http.Request, cookieName string) (*http.Cookie, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

// GetNavbarForRequest extracts the referer from the request and returns a Navbar struct.
// TODO: dynamically trim the host part of the referer URL
func GetNavbarForRequest(r *http.Request) Navbar {
	Referer := r.Referer()

	return Navbar{
		Referer: strings.TrimPrefix(Referer, "http://localhost:6969"),
	}
}

// This function renders the budget if auth is found, otherwise it renders the login page.
func BudgetData(w http.ResponseWriter, r *http.Request) {
	cookie, err := getCookieFromRequest(r, "SESSION")
	if err != nil {
		log.Println("No SESSION cookie found:", err)
	} else {
		log.Println("Found SESSION cookie:", cookie.Value)
		var tpl = template.Must(template.ParseFiles("components/budget.html"))
		tpl.Execute(w, nil)
		return
	}

	var tpl = template.Must(template.ParseFiles("components/login.html"))
	tpl.Execute(w, nil)
}

// Auth handles the authentication process and then redirects to login.
// If successful, it sets a session cookie
func Auth(w http.ResponseWriter, r *http.Request) {
	client := getSupabaseClient()
	username := r.FormValue("username")
	password := r.FormValue("password")
	session, err := client.SignInWithEmailPassword(username, password)
	if err != nil {
		log.Println("Error signing in:", err)
	} else {
		cookie := http.Cookie{
			Name:  "SESSION",
			Value: session.AccessToken,
		}
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
		log.Println("Session cookie set:", cookie)
	}

	BudgetData(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "SESSION",
		Value:  "",
		MaxAge: -1, // This will delete the cookie
	}
	http.SetCookie(w, &cookie)
	log.Println("Session cookie cleared")
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func TransactionQuery(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := getCookieFromRequest(r, "SESSION")

	session := types.Session{
		AccessToken: sessionCookie.Value,
	}

	client := getSupabaseClient()
	client.UpdateAuthSession(session)

	if err != nil {
		log.Println("Error creating Supabase client:", err)
	} else {
		data, count, err := client.From("transactions").Select("*", "exact", false).Execute()
		if err != nil {
			log.Println("Error fetching data:", err)
		}

		var d TransactionResponse
		err = json.Unmarshal(data, &d)
		log.Println("Data:", data)
		log.Println("d:", d)
		log.Println("d[0]:", d[0])
		log.Println("d[0].TransactionAmount:", d[0].TransactionAmount)
		log.Println("Count:", count)
	}
}
