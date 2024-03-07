// handlers/add_pen.go

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
	var err error // Declare err here so that it's accessible in the entire function

	// Get the user ID from the session
	userID := GetUserIDFromSession(r)
	if userID == 0 {
		// Redirect to the login page with an error message if the user is not logged in
		RedirectWithError(w, r, "/login", "Please login")
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		// Fetch dynamic column names based on user and table
		columns := GetColumnNames(userID, "pens")

		// Remove "id" from column names and values
		columns = columns[1:]
		columnValues := make([]interface{}, len(columns))

		// Extract form values and trim spaces
		for i, col := range columns {
			columnValues[i] = strings.TrimSpace(r.FormValue(col))
		}

		// Insert the pen using the InsertPen function from handlers
		err := InsertPen(userID, convertInterfaceToStringSlice(columnValues))
		if err != nil {
			log.Fatal(err)
		}

		// Redirect to the dashboard after successful insertion
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// For GET requests, render the form
	// Fetch dynamic column names based on user and table
	columns := GetColumnNames(userID, "pens")

	// Prepare data for template rendering
	data := struct {
		Columns     []string
		CurrentYear int
		Title       func(string) string // Function to capitalize and replace underscores
	}{
		Columns:     columns, // Include all columns, excluding "id"
		CurrentYear: time.Now().Year(),
		Title:       Title, // Pass the Title function to the template
	}

	// log.Printf("Data for adding pen is %=v", data)

	// Parse and execute the template
	tmpl := template.Must(template.New("add.html").Funcs(template.FuncMap{"Title": Title}).ParseFiles("templates/add.html"))
	if err != nil {
		log.Fatal("Error parsing add.html template:", err)
	}
	tmpl.Execute(w, data)
}
