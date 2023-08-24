// handlers/import_export.go

package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strings"
	"time"
	"html/template"
)

func ExportCSV(w http.ResponseWriter, r *http.Request) {

	pens, columns, err := SelectPens()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")

	// Generate the filename based on the current date
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("flock_%s_backup.csv", timestamp)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Write CSV header
	if err := csvWriter.Write(columns); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error writing CSV header:", err)
		return
	}

	// Write CSV rows
	for _, pen := range pens {
		row := make([]string, len(columns))
		for i, col := range columns {
			row[i] = fmt.Sprintf("%v", pen[col])
		}
		if err := csvWriter.Write(row); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error writing CSV row:", err)
			return
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error flushing CSV writer:", err)
		return
	}
}

func ImportCSV(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(10 << 20) // Max memory usage for uploaded files

		file, _, err := r.FormFile("csvfile")
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error retrieving uploaded file:", err)
			return
		}
		defer file.Close()

		log.Println("CSV file retrieved successfully.")

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error reading CSV:", err)
			return
		}

		columns := rows[0] // Assume the first row contains column headers
		rows = rows[1:]    // Exclude the header row

		//log.Printf("Imported columns: %v\n", columns)
		//log.Printf("Imported rows: %v\n", rows)

		data := struct {
			Columns []string
			Rows    [][]string
		}{
			Columns: columns,
			Rows:    rows,
		}

		csvDataJSON, err := json.Marshal(data.Rows)



		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error marshaling CSV data:", err)
			return
		}

		// Print the value of csvDataJSON
		//fmt.Println("csvDataJSON:", string(csvDataJSON))

		columnsJSON, err := json.Marshal(data.Columns)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error marshaling columns data:", err)
			return
		}

		tmpl := template.Must(template.ParseFiles("templates/import_preview.html"))
		tmpl.Execute(w, struct {
			CsvData   template.JS
			Columns   template.JS
		}{
			CsvData:   template.JS(csvDataJSON),
			Columns:   template.JS(columnsJSON),
		})
		return

	}

	tmpl := template.Must(template.ParseFiles("templates/import.html"))
	tmpl.Execute(w, nil)
}

func ImportApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		csvData := r.FormValue("csvData")
		var rows [][]string
		if err := json.Unmarshal([]byte(csvData), &rows); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error unmarshaling CSV data:", err)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error beginning transaction:", err)
			return
		}

		for _, row := range rows {
			// Start from index 1 to exclude the id column
			if err := InsertPen(row[1:]); err != nil {
				tx.Rollback()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error inserting pen:", err)
				return
			}
		}

		err = tx.Commit()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error committing transaction:", err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}