// Package handlers provides functionality to interact with the database and handle data operations.
package handlers

import (
	"net/http"
	"strconv"
	// "log"
)

// DeletePen handles the deletion of a pen from the database.
func DeletePen(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromSession(r)
	if userID == 0 {
		RedirectWithError(w, r, "/login", "Please login to delete a pen from your database")
		return
	}

	// Get the pen ID from the URL parameter
	penID, err := strconv.ParseInt(r.URL.Path[len("/delete/"):], 10, 64)
	if err != nil {
		RedirectWithError(w, r, "/dashboard", "Please try to delete a pen and not fight the world")
		return
	}

	// Check if the pen exists before attempting to delete
	if !PenExists(userID, penID) {
		RedirectWithError(w, r, "/dashboard", "You can only delete a pen if it exists")
		return
	}

	// Delete the pen using the DeletePenByID function from handlers
	err = DeletePenByID(userID, penID)
	if err != nil {
		RedirectWithError(w, r, "/dashboard", "Please try to delete once more")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
