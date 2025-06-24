package supabaseHandler

import (
	"encoding/csv"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"

	util "htmx-demo/util"
)

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
func GetSupabaseClient() *supabase.Client {
	client, err := supabase.NewClient("https://oyucpcsumehntdxplvgl.supabase.co", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6Im95dWNwY3N1bWVobnRkeHBsdmdsIiwicm9sZSI6ImFub24iLCJpYXQiOjE3NTA0Njc0NjAsImV4cCI6MjA2NjA0MzQ2MH0.N9RIgBcDY6t4cpvohSz7h0XE9cctp6Dl1hhepxFsOrc", &supabase.ClientOptions{})
	if err != nil {
		log.Println("Error creating Supabase client:", err)
		return nil
	}
	return client
}

func Authenticate(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, *http.Request, error) {
	client := GetSupabaseClient()
	email := r.FormValue("email")
	password := r.FormValue("password")
	session, err := client.SignInWithEmailPassword(email, password)
	if err != nil {
		log.Println("Error signing in:", err)
	} else {
		cookie := http.Cookie{
			Name:  "SESSION",
			Value: session.AccessToken,
		}
		http.SetCookie(w, &cookie)
		r.AddCookie(&cookie)
	}

	return w, r, err
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := util.GetCookieFromRequest(r, "SESSION")
	if err != nil {
		log.Println("No SESSION cookie found:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	client := GetSupabaseClient()
	client.UpdateAuthSession(types.Session{AccessToken: sessionCookie.Value})

	amountString := r.FormValue("transactionAmount")
	dateString := r.FormValue("transactionDate")

	transactionAmount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		log.Println("Error parsing transaction amount:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	transaction := TransactionCreateType{
		TransactionDate:   dateString,
		TransactionAmount: int(transactionAmount * 100), // Convert to cents
	}

	_, _, err = client.From("transactions").Insert(transaction, false, "", "", "").Execute()
	if err != nil {
		log.Println("Error inserting transaction:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/budget", http.StatusSeeOther)
}

func getTransactionFromRecord(record []string, headerToIndex map[string]int, isDebit bool) (TransactionCreateType, error) {
	if isDebit {
		index := headerToIndex["Amount"]
		transactionType := record[headerToIndex["Type"]]

		stringAmount := record[index]
		amount, _ := strconv.ParseFloat(stringAmount, 64)
		log.Println("String Amount", stringAmount)
		log.Println("index", index)
		log.Println("record:", record)

		if transactionType == "CREDIT" {
			amount = -amount // Convert credit to negative amount
		}

		return TransactionCreateType{
			TransactionDate:   record[headerToIndex["Date"]],
			TransactionAmount: int(amount * 100),
		}, nil
	} else {
		creditAmountIndex := headerToIndex["Credit"]
		creditAmount := record[creditAmountIndex]
		debitAmountIndex := headerToIndex["Debit"]
		debitAmount := record[debitAmountIndex]
		// amount := record[amountIndex]
		if creditAmount == "" && debitAmount == "" {
			return TransactionCreateType{}, nil // No transaction to create
		}

		log.Println("credit/debit Amount", creditAmount, debitAmount)
		log.Println("credit/debit Index", creditAmountIndex, debitAmountIndex)
		log.Println("headerToIndex:", headerToIndex)
		log.Println("record:", record)
		if creditAmount != "" {
			amount, err := strconv.ParseFloat(record[creditAmountIndex], 64)
			if err != nil {
				log.Println("Error converting credit amount:", err)
				return TransactionCreateType{}, err
			}

			return TransactionCreateType{
				TransactionDate:   record[headerToIndex["Transaction Date"]],
				TransactionAmount: int(-amount * 100), // Credit amount is negative
			}, nil

		} else {
			amount, err := strconv.ParseFloat(record[debitAmountIndex], 64)
			if err != nil {
				log.Println("Error converting debit amount:", err)
				return TransactionCreateType{}, err
			}
			// Debit amount is positive
			return TransactionCreateType{
				TransactionDate:   record[headerToIndex["Transaction Date"]],
				TransactionAmount: int(amount * 100),
			}, nil
		}
	}
}

func BulkUpload(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := util.GetCookieFromRequest(r, "SESSION")
	if err != nil {
		log.Println("No SESSION cookie found:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	client := GetSupabaseClient()
	client.UpdateAuthSession(types.Session{AccessToken: sessionCookie.Value})

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	log.Println("Received file from form:", file)
	if err != nil {
		log.Println("Error retrieving file from form:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// var transactions []TransactionCreateType
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Error reading CSV file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	debitHeaders := []string{"Date", "Amount", "Type", "Description"}
	creditHeaders := []string{"Transaction Date", "Description", "Debit", "Credit"}

	headers := records[0]
	debitHeaderToRecordIndex := make(map[string]int)
	creditHeaderToRecordIndex := make(map[string]int)
	headerToRecordIndex := make(map[string]int)
	isDebit := true
	isCredit := true

	for _, header := range debitHeaders {
		index := slices.Index(headers, header)
		if index == -1 {
			isDebit = false
		} else {
			debitHeaderToRecordIndex[header] = index
		}
	}

	for _, header := range creditHeaders {
		index := slices.Index(headers, header)
		if index == -1 {
			isCredit = false
		} else {
			creditHeaderToRecordIndex[header] = index
		}
	}

	if isDebit {
		headerToRecordIndex = debitHeaderToRecordIndex
	} else if isCredit {
		headerToRecordIndex = creditHeaderToRecordIndex
	} else {
		log.Println("CSV file does not match expected headers for debit or credit transactions.")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	log.Println("Is Debit:", isDebit)

	transactions := []TransactionCreateType{}
	for i, record := range records {
		if i == 0 {
			log.Println("header row:", record)
		} else {
			transaction, err := getTransactionFromRecord(record, headerToRecordIndex, isDebit)

			if err != nil {
				log.Println("Error getting transaction from record:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if transaction.TransactionAmount != 0 {
					transactions = append(transactions, transaction)
					log.Println("Transaction created:", transaction)
				} else {
					log.Println("No transaction created for record:", record)
				}
			}
		}
	}

	// client.From("transactions").BulkUpload(transactions, false, "", "", "").Execute()
	_, _, err = client.From("transactions").Insert(transactions, false, "", "", "").Execute()
	if err != nil {
		log.Println("Error inserting transactions:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// log.Println("Bulk upload successful:", response)
	http.Redirect(w, r, "/budget", http.StatusSeeOther)
}
