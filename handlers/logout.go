// handlers/logout.go

package handlers

import (
	"net/http"
)

// Logout handles user logout
func Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the user session
	session, _ := store.Get(r, "session-name")
	session.Values["userID"] = nil
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
