package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"git.fhict.nl/I470668/bookingsystemv2/platform/authenticator"
	"git.fhict.nl/I470668/bookingsystemv2/platform/router"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var templates *template.Template
var db *sql.DB

func errorLog() {
	logFile, err := os.OpenFile(".\\trace.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Couldn't create logfile")
	}
	log.SetOutput(logFile)
}

func init() {
	templatesPath := filepath.Join("templates", "*.html")
	templates = template.Must(template.New("").ParseGlob(templatesPath))

	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Couldn't open config file")
	}
	defer configFile.Close()

	// Decode JSON from config file
	var config struct {
		DBConnectionString string `json:"connString"`
	}
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		log.Fatal("Couldn't decode config file")
	}

	// Open a connection to the MySQL database
	db, err = sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	errorLog()
	defer db.Close()

	http.HandleFunc("/", handleLogin)
	http.HandleFunc("/dashboard/", handleDashboard)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	rtr := router.New(auth)

	fmt.Println("Server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var username, password string
		err := r.ParseForm()
		if err != nil {
			log.Println("error parsing login form:", err)
		}
		username = r.FormValue("username")
		password = r.FormValue("password")
		row := db.QueryRow("SELECT username FROM users WHERE username=? AND password=?", username, password)
		if err := row.Scan(&username); err != nil {
			// If the login was unsuccessful, render the login page with an error message
			errMsg := "username or password is incorrect"
			if err := templates.ExecuteTemplate(w, "login.html", errMsg); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				log.Println("error rendering login page:", err)
			}
			return
		}
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
		return
	}

	if err := templates.ExecuteTemplate(w, "login.html", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("error rendering login page:", err)
	}
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	if err := templates.ExecuteTemplate(w, "dashboard.html", nil); err != nil {
		log.Println("error:", err)
	}
}
