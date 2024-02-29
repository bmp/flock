// Package handlers provides functionality to interact with the database and handle data operations.
package handlers

import (
	"net/http"
	"strconv"
)

// DeletePen handles the deletion of a pen from the database.
func DeletePen(w http.ResponseWriter, r *http.Request) {
	// Get the pen ID from the URL parameter
	penID, err := strconv.ParseInt(r.URL.Path[len("/delete/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid pen ID", http.StatusBadRequest)
		return
	}

	// Delete the pen using the DeletePenByID function from handlers
	err = DeletePenByID(penID)
	if err != nil {
		http.Error(w, "Error deleting pen", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
