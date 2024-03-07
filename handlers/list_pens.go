// Package handlers provides functionality to interact with the database and handle data operations.

package handlers

import (
	"html/template"
	// "log"
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

	// Define data at the beginning
	var data struct {
		Pens    []map[string]interface{}
		Columns []string
		Error           string
		RedirectURL     string
	}

	// Get the user ID from the session (you need to implement this part)
	userID := GetUserIDFromSession(r)
	if userID == 0 {
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		RedirectWithError(w, r, "/login", "Please login")
		return
	}

	// Fetch pens and columns from the user's pens database
	pens, columns, err := SelectPens(userID)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// log.Println("Error fetching data:", err)
		RedirectWithError(w, r, "/login", "User tables don't exist try again")
		return
	}

	// Prepare data for template rendering
	data.Pens = pens
	data.Columns = columns

	// Check if there's any error message or redirection URL in the query parameters
	queryParams := r.URL.Query()
	if len(queryParams["error"]) > 0 {
		data.Error = queryParams["error"][0]
		// log.Printf("adding %s to data", data.Error)
	}
	if len(queryParams["redirect"]) > 0 {
		data.RedirectURL = queryParams["redirect"][0]
		// log.Printf("adding %s to data", data.RedirectURL)
	}

	// log.Printf("Calling with %+v", data)

	// Parse and execute the template
	tmpl := template.Must(template.New("dashboard.html").Funcs(template.FuncMap{"Add": Add}).ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, data)
}
