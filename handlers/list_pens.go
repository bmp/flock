// Package handlers provides functionality to interact with the database and handle data operations.

package handlers

import (
	"html/template"
	"log"
	"net/http"
	//"time"
)

// ListPens retrieves a list of pens from the database and renders them using a template.
//
// The ListPens function handles the HTTP request to display a list of pens stored in the database.
// It queries the database for pen records, extracts the pen data and column names, and then passes
// the data to a template for rendering. If any error occurs during data retrieval or rendering, it
// returns an HTTP 500 error and logs the error message.
//
// Parameters:
//   - w (http.ResponseWriter): The HTTP response writer to write the response to.
//   - r (*http.Request): The HTTP request containing details of the request.
func ListPens(w http.ResponseWriter, r *http.Request) {
	// Fetch pens and columns from the database
	pens, columns, err := SelectPens()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	// Prepare data for template rendering
	data := struct {
		Pens    []map[string]interface{}
		Columns []string
	}{
		Pens:    pens,
		Columns: columns,
	}

	// Parse and execute the template
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl := template.Must(template.New("index.html").Funcs(template.FuncMap{"Add": Add}).ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}
