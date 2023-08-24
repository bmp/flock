// main.go is the entry point of Flock
package main

import (
	"database/sql"
	"path/filepath"
	//"html/template"
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

	http.HandleFunc("/", handlers.ListPens)                    // Handle listing pens
	http.HandleFunc("/add", handlers.AddPen)                   // Handle adding a pen
	http.HandleFunc("/export/csv", handlers.ExportCSV)         // Handle exporting to CSV
	http.HandleFunc("/import/csv", handlers.ImportCSV)         // Handle importing from CSV
	http.HandleFunc("/import/approve", handlers.ImportApprove) // Handle approving imported data from CSV

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
