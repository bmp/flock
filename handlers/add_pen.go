// Package handlers provides functionality to interact with the database and handle data operations.
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)


// AddPen handles the addition of a new pen to the database.
// It processes both GET and POST requests. For GET requests, it renders
// the add pen form with dynamic column names and the current year.
// For POST requests, it extracts form values, prepares column values,
// and inserts the pen into the database using the InsertPen function.
func AddPen(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		columns := GetColumnNames("pens") // Fetch column names dynamically using handler function
		// Remove "id" from column names and values
		columns = columns[1:]
		columnValues := make([]interface{}, len(columns))

		for i, col := range columns {
			columnValues[i] = strings.TrimSpace(r.FormValue(col))
		}

		// Insert the pen using the InsertPen function from handlers
		err := InsertPen(convertInterfaceToStringSlice(columnValues))
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	columns := GetColumnNames("pens") // Fetch column names dynamically using handler function

	data := struct {
		Columns     []string
		CurrentYear int
	}{
		Columns:     columns, // Include all columns, excluding "id"
		CurrentYear: time.Now().Year(),
	}

	tmpl := template.Must(template.ParseFiles("templates/add.html"))
	tmpl.Execute(w, data)
}

// convertInterfaceToStringSlice converts a slice of interfaces to a slice of strings.
// It filters out non-string values and returns a string slice.
func convertInterfaceToStringSlice(interfaceSlice []interface{}) []string {
    stringSlice := make([]string, len(interfaceSlice))
    for i, v := range interfaceSlice {
        if value, ok := v.(string); ok {
            stringSlice[i] = value
        }
    }
    return stringSlice
}
