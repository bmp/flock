// handlers/list_pens.go

package handlers

import (
	"html/template"
	"log"
	"net/http"
	//"time"
)

func ListPens(w http.ResponseWriter, r *http.Request) {
	pens, columns, err := SelectPens()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	data := struct {
		Pens    []map[string]interface{}
		Columns []string
	}{
		Pens:    pens,
		Columns: columns,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}
