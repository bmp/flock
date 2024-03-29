// handlers/index.go

package handlers

import (
	// "html/template"
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
func Index(w http.ResponseWriter, r *http.Request) {

	// // Parse and execute the template
	// tmpl := template.Must(template.ParseFiles("templates/index.html"))
	// // tmpl := template.Must(template.New("dashboard.html").Funcs(template.FuncMap{"Add": Add}).ParseFiles("templates/dashboard.html"))
	// tmpl.Execute(w, nil)

	// Check if the user is already authenticated
	userID := GetUserIDFromSession(r)
	if userID != 0 {
		// If the user is logged in, redirect to the dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Render the index page for non-authenticated users
	renderTemplate(w, "index", nil)

}
