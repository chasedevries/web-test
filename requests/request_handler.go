package requestHandler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

type NavbarProps struct {
	Referer    string `json:"referer"`
	IsLoggedIn bool   `json:"is_logged_in"`
}

type Error struct {
	ErrorMessage string `json:"error_message"`
}

type TransactionQueryType struct {
	Id                int    `json:"id"`
	CreatedAt         string `json:"created_at"`
	TransactionDate   string `json:"transaction_date"`
	TransactionAmount int    `json:"transaction_amount"`
	UserId            string `json:"user_id"`
}

type TransactionCreateType struct {
	TransactionDate   string `json:"transaction_date"`
	TransactionAmount int    `json:"transaction_amount"`
}

type TransactionResponse []TransactionQueryType

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
func GetNavbarForRequest(r *http.Request) NavbarProps {
	Referer := r.Referer()
	Host := r.Host
	_, err := getCookieFromRequest(r, "SESSION")
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
	cookie, err := getCookieFromRequest(r, "SESSION")

	if err != nil {
		log.Println("No SESSION cookie found:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	client := getSupabaseClient()
	client.UpdateAuthSession(types.Session{AccessToken: cookie.Value})
	data, _, err := client.From("transactions").Select("*", "exact", false).Execute()
	log.Println("Transaction Data:", string(data))
	transactions := TransactionResponse{}
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

// Auth handles the authentication process and then redirects to login.
// If successful, it sets a session cookie
func Auth(w http.ResponseWriter, r *http.Request) {
	client := getSupabaseClient()
	email := r.FormValue("email")
	password := r.FormValue("password")
	session, err := client.SignInWithEmailPassword(email, password)
	if err != nil {
		log.Println("Error signing in:", err)
		LoginForm(w, r)
		return
	} else {
		cookie := http.Cookie{
			Name:  "SESSION",
			Value: session.AccessToken,
		}
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
		log.Println("Session cookie set:", cookie)
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

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := getCookieFromRequest(r, "SESSION")
	if err != nil {
		log.Println("No SESSION cookie found:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	client := getSupabaseClient()
	client.UpdateAuthSession(types.Session{AccessToken: sessionCookie.Value})

	amountString := r.FormValue("transactionAmount")
	dateString := r.FormValue("transactionDate")

	log.Println("form values:", r.Form)
	log.Println("Received transaction amount:", amountString)
	log.Println("Received transaction date:", dateString)

	transactionAmount, err := strconv.Atoi(amountString)
	transaction := TransactionCreateType{
		TransactionDate:   dateString,
		TransactionAmount: transactionAmount,
	}

	response, _, err := client.From("transactions").Insert(transaction, false, "", "", "").Execute()
	if err != nil {
		log.Println("Error inserting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Println("Transaction created successfully:", response)
	http.Redirect(w, r, "/budget", http.StatusSeeOther)
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
