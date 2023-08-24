// Package handlers provides functionality to interact with the database and handle data operations.

package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"reflect"
	"path/filepath"
	"os"
)

// dbName is the name of the SQLite database file.
const dbName = "database.db"

var db *sql.DB

// CreateDatabaseIfNotExists checks if the database exists, and creates it if not.
// If the database already exists, it simply opens it.
func CreateDatabaseIfNotExists() (*sql.DB, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the path to the database file
	dbPath := filepath.Join(currentDir, "database", dbName)

	// Check if the database file already exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// If the database file doesn't exist, create it
		log.Println("Creating database...")
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}

		// Create the table structure
		_, err = db.Exec(`CREATE TABLE pens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			maker TEXT,
			color TEXT,
			material TEXT,
			nib_size TEXT,
			nib_color TEXT,
			filling_system TEXT,
			trims TEXT,
			year INTEGER,
			price REAL,
			misc TEXT
		)`)
		if err != nil {
			db.Close()
			return nil, err
		}

		return db, nil
	}

	// If the database file already exists, simply open it
	log.Println("Opening existing database...")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// InitDB initializes the database connection for handlers.
func InitDB(database *sql.DB) {
	db = database
}

// fetchDataFromDB fetches data from the database based on the provided query.
func fetchDataFromDB(query string) ([]string, []map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	var pens []map[string]interface{}

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

// GetColumnNames retrieves the column names of the specified table.
func GetColumnNames(tableName string) []string {
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

// SelectPens fetches all pens from the database.
func SelectPens() ([]map[string]interface{}, []string, error) {
	query := "SELECT * FROM pens"
	columns, pens, err := fetchDataFromDB(query)
	if err != nil {
		return nil, nil, err
	}
	return pens, columns, nil
}

// InsertPen inserts a new pen record into the database.
func InsertPen(values []string) error {
	columns := GetColumnNames("pens")
	columns = columns[1:] // Exclude "id"

	valuePlaceholders := make([]string, len(columns))
	for i := range columns {
    valuePlaceholders[i] = "?"
	}

	insertQuery := fmt.Sprintf("INSERT INTO pens (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(valuePlaceholders, ", "))

	_, err := db.Exec(insertQuery, convertStringSliceToInterfaceSlice(values)...)
	if err != nil {
		return err
	}
	return nil
}

// ModifyPen updates a pen record in the database.
func ModifyPen(id int64, values []string) error {
    columns := GetColumnNames("pens")
    columns = columns[1:] // Exclude "id"

    setStatements := make([]string, len(columns))
    for i, col := range columns {
        setStatements[i] = fmt.Sprintf("%s = ?", col)
    }

    updateQuery := fmt.Sprintf("UPDATE pens SET %s WHERE id = ?", strings.Join(setStatements, ", "))
    values = append(values, fmt.Sprintf("%d", id)) // Add the ID to the end of the values
    _, err := db.Exec(updateQuery, interfaceSlice(values)...)
    if err != nil {
        return err
    }
    return nil
}

// Convert []string to []interface{}
func interfaceSlice(slice []string) []interface{} {
    interfaceSlice := make([]interface{}, len(slice))
    for i, v := range slice {
        interfaceSlice[i] = v
    }
    return interfaceSlice
}

// Convert []interface{} to []string
func interfaceSliceToString(slice []interface{}) []string {
    stringSlice := make([]string, len(slice))
    for i, v := range slice {
        stringSlice[i] = fmt.Sprintf("%v", v)
    }
    return stringSlice
}

// Convert []string to []interface{}
func convertStringSliceToInterfaceSlice(slice []string) []interface{} {
    interfaceSlice := make([]interface{}, len(slice))
    for i, v := range slice {
        interfaceSlice[i] = v
    }
    return interfaceSlice
}
