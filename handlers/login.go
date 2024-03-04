// handlers/login.go
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
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Compare the hashed password with the provided password
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Authentication successful, set user ID in the session
		userID, err := GetUserIDByUsername(username) // Implement this function to get the user ID based on the username
		if err != nil {
			http.Error(w, "Error getting user ID", http.StatusInternalServerError)
			return
		}

		SetUserIDInSession(w, r, userID)

		// Redirect to the user dashboard
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Render the login form for GET requests
	renderTemplate(w, "login", nil)
}
