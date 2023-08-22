// main.go
package main

import (
	"database/sql"

	"html/template"
	"log"
	"net/http"
	"strings"
	"time"	

	"reflect"	

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
		columns := getColumnNames("pens") // Fetch column names dynamically
		// Remove "id" from column names and values
		columns = columns[1:]
		columnValues := make([]interface{}, len(columns))
		for i, col := range columns {
			columnValues[i] = r.FormValue(col)
		}

		query := "INSERT INTO pens (" + strings.Join(columns, ", ") + ") VALUES (" +
			strings.Join(strings.Split(strings.Repeat("?", len(columns)), ""), ", ") + ")"

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
		Columns: columns[1:], // Exclude "id"
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
		displayName := strings.Title(strings.ReplaceAll(name, "_", " "))
		columns = append(columns, displayName)
	}

	return columns
}


