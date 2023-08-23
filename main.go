// main.go
package main

import (
	"database/sql"

	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
	"encoding/csv"
	"fmt"
	"reflect"
	//"os"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", listPens)
	http.HandleFunc("/add", addPen)
	http.HandleFunc("/export/csv", exportCSV)
	http.HandleFunc("/import/csv", importCSV)
	http.HandleFunc("/import/approve", importApprove)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/includes/", http.StripPrefix("/includes/", http.FileServer(http.Dir("includes"))))


	// Serve JavaScript files with the correct MIME type
	http.HandleFunc("/includes/sort.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript") // Set the MIME type
		http.ServeFile(w, r, "includes/sort.js")
	})


	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchDataFromDB(query string) (columns []string, pens []map[string]interface{}, err error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columns, _ = rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		scanArgs := make([]interface{}, len(columns))
		for i := range values {
			var v interface{}
			scanArgs[i] = &v
			values[i] = &v
		}
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, nil, err
		}

		retrievedPen := make(map[string]interface{})
		for i, colName := range columns {
			if colName == "id" {
				retrievedPen[colName] = reflect.ValueOf(values[i]).Elem().Interface().(int64) // Cast to int64
			} else {
				retrievedPen[colName] = reflect.ValueOf(values[i]).Elem().Interface()
			}
		}

		pens = append(pens, retrievedPen)
	}

	return columns, pens, nil
}

func listPens(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM pens"

	columns, pens, err := fetchDataFromDB(query)
	if err != nil {
		log.Fatal(err)
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

func addPen(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		//log.Printf("Form Values: %+v\n", r.Form)

		columns := getColumnNames("pens") // Fetch column names dynamically
		// Remove "id" from column names and values
		columns = columns[1:]
		columnValues := make([]interface{}, len(columns))

		for i, col := range columns {
			columnValues[i] = strings.TrimSpace(r.FormValue(col))
		}

		valuePlaceholders := make([]string, len(columns))

		for i, col := range columns {
			formValue := r.FormValue(col)
			columnValues[i] = formValue
			valuePlaceholders[i] = "?"
		}

		query := "INSERT INTO pens (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(valuePlaceholders, ", ") + ")"

		//log.Printf("Query: %s\n", query)
		//log.Printf("Values: %v\n", columnValues)

		_, err := db.Exec(query, columnValues...)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	columns := getColumnNames("pens") // Fetch column names dynamically
	data := struct {
		Columns []string
		CurrentYear int
	}{
		Columns:  columns[1:], // Exclude "id"
		CurrentYear: time.Now().Year(),

	}

	//log.Printf("Debug: CurrentYear = %d\n", data.CurrentYear)

	tmpl := template.Must(template.ParseFiles("templates/add.html"))
	tmpl.Execute(w, data)
}



func getColumnNames(tableName string) []string {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ")")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name string
		// Other unused columns can be scanned into placeholders
		err := rows.Scan(&cid, &name, new(interface{}), new(interface{}), new(interface{}), new(interface{}))
		if err != nil {
			log.Fatal(err)
		}
		columns = append(columns, name)
	}

	return columns
}

func exportCSV(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM pens"

	columns, pens, err := fetchDataFromDB(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=pens.csv")

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

func importCSV(w http.ResponseWriter, r *http.Request) {
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

		log.Printf("Imported columns: %v\n", columns)
		log.Printf("Imported rows: %v\n", rows)

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
		fmt.Println("csvDataJSON:", string(csvDataJSON))

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



func importApprove(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		csvData := r.FormValue("csvData")
		var rows [][]string
		if err := json.Unmarshal([]byte(csvData), &rows); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error unmarshaling CSV data:", err)
			return
		}

		columnsJSON := r.FormValue("columns")
		var columns []string
		if err := json.Unmarshal([]byte(columnsJSON), &columns); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error unmarshaling columns data:", err)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error beginning transaction:", err)
			return
		}

		for _, row := range rows {
			valuePlaceholders := strings.Repeat("?, ", len(row)-2) + "?"

			statement := "INSERT INTO pens (" + strings.Join(columns[1:], ", ") + ") VALUES (" + valuePlaceholders + ")"

			fmt.Println("insert query:", statement)

			stmt, err := tx.Prepare(statement)

			if err != nil {
				tx.Rollback()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error preparing statement:", err)
				return
			}
			defer stmt.Close()

			args := make([]interface{}, len(row)-1)
			for i, v := range row[1:] { // Start from index 1 to exclude the id column
				args[i] = v
			}

			_, err = stmt.Exec(args...)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println("Error executing statement:", err)
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
