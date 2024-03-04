// main.go is the entry point of Flock
package main

import (
	"database/sql"
	"path/filepath"
	//"html/template"
	//"text/template"
	"log"
	"net/http"
	//"strings"
	//"time"
	//"encoding/csv"
	//"fmt"
	//"reflect"
	//"os"
	//"encoding/json"

	_ "github.com/mattn/go-sqlite3"

	// "golang.org/x/crypto/bcrypt"

	"flock/handlers"
)

var db *sql.DB

func main() {

	// Check if the database exists, and create or open it
	db, err := handlers.CreateDatabaseIfNotExists()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize the database for handlers
	handlers.InitDB(db)

	log.Println("Database connection established")

	http.HandleFunc("/", handlers.Index)                       // Handler listing pens
	http.HandleFunc("/register", handlers.Register)            // Handler for registering user
	http.HandleFunc("/login", handlers.Login)                  // Handler for login
	http.HandleFunc("/dashboard", handlers.ListPens)           // Handler listing pens
	http.HandleFunc("/add", handlers.AddPen)                   // Handler adding a pen
	http.HandleFunc("/export/csv", handlers.ExportCSV)         // Handler exporting to CSV
	http.HandleFunc("/import/csv", handlers.ImportCSV)         // Handler importing from CSV
	http.HandleFunc("/import/approve", handlers.ImportApprove) // Handler approving imported data from CSV
	http.HandleFunc("/modify/", handlers.ModifyPen)            // Handler to modify details for a pen
	http.HandleFunc("/delete/", handlers.DeletePen)            // Handler to delete a pen
	http.HandleFunc("/logout", handlers.Logout)                // Handler for logout

	// Serve static assets
	http.HandleFunc("/includes/", func(w http.ResponseWriter, r *http.Request) {
		// Get the requested file path
		filePath := "." + r.URL.Path

		// Determine the file extension
		fileExt := filepath.Ext(filePath)

		// Set the Content-Type header based on the file extension
		switch fileExt {
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		}

		// Serve the file using http.FileServer
		http.FileServer(http.Dir(".")).ServeHTTP(w, r)
	})

	port := ":8000"
	log.Println("Starting server on", port)
	log.Fatal(http.ListenAndServe(port, nil))

}
