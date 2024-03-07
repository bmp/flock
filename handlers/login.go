// login.go
package handlers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Extract login credentials from the form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Retrieve hashed password from the database based on the username
		hashedPassword, err := GetPasswordByUsername(username)
		if err != nil {
			RedirectWithError(w, r, "/login", "Invalid username or password")
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			RedirectWithError(w, r, "/login", "Invalid username or password")
			return
		}

		// Authentication successful, set user ID in the session
		userID, err := GetUserIDByUsername(username)
		if err != nil {
			RedirectWithError(w, r, "/login", "Error getting user ID")
			return
		}

		SetUserIDInSession(w, r, userID)

		// Redirect to the user dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Initialize error message and redirection URL
	var errorMessage string
	var redirectionURL string

	// Check if there's any error message or redirection URL in the query parameters
	queryParams := r.URL.Query()
	if len(queryParams["error"]) > 0 {
		errorMessage = queryParams["error"][0]
	}
	if len(queryParams["redirect"]) > 0 {
		redirectionURL = queryParams["redirect"][0]
	}

	// Render the login form along with error message and redirection URL
	renderTemplate(w, "login", map[string]interface{}{
		"Error":       errorMessage,
		"RedirectURL": redirectionURL,
	})
}
