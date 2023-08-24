// Package handlers provides functionality to interact with the database and handle data operations.
package handlers

import (
	//"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// ModifyPen handles the modification of a pen in the database.
// It processes both GET and POST requests. For GET requests, it renders
// the modify pen form with pre-filled values based on the pen ID.
// For POST requests, it extracts form values, prepares column values,
// and updates the pen in the database using the ModifyPen function.
func ModifyPen(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		// Get the pen ID from the URL parameter
		penID, err := strconv.ParseInt(r.URL.Path[len("/modify/"):], 10, 64)
		if err != nil {
			http.Error(w, "Invalid pen ID", http.StatusBadRequest)
			return
		}

		columns := GetColumnNames("pens") // Fetch column names dynamically using handler function
		// Remove "id" from column names
		columns = columns[1:]

		columnValues := make([]interface{}, len(columns))
		for i, col := range columns {
			columnValues[i] = r.FormValue(col)
		}

		// Update the pen using the ModifyPen function from handlers
		err = UpdatePen(penID, convertInterfaceToStringSlice(columnValues))
		if err != nil {
			http.Error(w, "Error modifying pen", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get the pen ID from the URL parameter
	penID, err := strconv.ParseInt(r.URL.Path[len("/modify/"):], 10, 64)
	if err != nil {
		http.Error(w, "Invalid pen ID", http.StatusBadRequest)
		return
	}

	// Fetch pen details based on ID
	pen, err := GetPenByID(penID)
	if err != nil {
		http.Error(w, "Error fetching pen details", http.StatusInternalServerError)
		return
	}

	columns := GetColumnNames("pens") // Fetch column names dynamically using handler function

	data := struct {
		Columns []string
		Pen     map[string]interface{}
	}{
		Columns: columns, // Include all columns, excluding "id"
		Pen:     pen,
	}

	//fmt.Println("Data:", data)

	tmpl := template.Must(template.ParseFiles("templates/modify.html"))
	tmpl.Execute(w, data)
}
