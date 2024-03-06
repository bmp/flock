// Package handlers provides functionality to interact with the database and handle data operations.

package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// dbName is the name of the SQLite database file.
const dbName = "database.db"

var db *sql.DB

// CreateDatabaseIfNotExists checks if the database exists and creates it if not.
// If the database already exists, it simply opens it.
func CreateDatabaseIfNotExists() (*sql.DB, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the path to the main database file
	dbPath := filepath.Join(currentDir, "database", dbName)

	// Check if the main database file already exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		// If the main database file doesn't exist, create it
		log.Println("Creating main database...")

		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}

		// Create the users table
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			first_name TEXT NOT NULL,
			middle_name TEXT,
			last_name TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			bio TEXT
		)`)
		if err != nil {
			db.Close()
			return nil, err
		}

		return db, nil
	}

	// If the main database file already exists, simply open it
	log.Println("Opening existing main database...")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetUserDBPath returns the path to the user's pens database file.
func GetUserDBPath(userID int64) string {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct the path to the user's pens database file
	userDBPath := filepath.Join(currentDir, "database", fmt.Sprintf("%d_pens.db", userID))

	return userDBPath
}

// CreateOrUpdateUserDB creates or updates the user's pens database.
func CreateOrUpdateUserDB(userID int64) (*sql.DB, error) {
	// Get the path to the user's pens database file
	userDBPath := GetUserDBPath(userID)

	// Check if the user's pens database file already exists
	if _, err := os.Stat(userDBPath); os.IsNotExist(err) {
		// If the user's pens database file doesn't exist, create it
		log.Printf("Creating user's pens database for user with ID %d...\n", userID)

		userDB, err := sql.Open("sqlite3", userDBPath)
		if err != nil {
			return nil, err
		}

		// Create the fountainpens table
		_, err = userDB.Exec(`CREATE TABLE IF NOT EXISTS pens (
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
			userDB.Close()
			return nil, err
		}

		return userDB, nil
	}

	// If the user's pens database file already exists, simply open it
	log.Printf("Opening existing user's pens database for user with ID %d...\n", userID)
	userDB, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, err
	}

	return userDB, nil
}

// InitDB initializes the database connection for handlers.
func InitDB(database *sql.DB) {
	db = database
}

// fetchDataFromDB fetches data from the database based on the provided query.
func fetchDataFromDB(db *sql.DB, query string) ([]string, []map[string]interface{}, error) {
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

// GetColumnNames retrieves the column names of the specified table in the user's database.
func GetColumnNames(userID int64, tableName string) []string {
	db, err := sql.Open("sqlite3", GetUserDBPath(userID))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("PRAGMA table_xinfo(" + tableName + ")")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name string
		// Other unused columns can be scanned into placeholders
		err := rows.Scan(&cid, &name, new(interface{}), new(interface{}), new(interface{}), new(interface{}), new(interface{},))
		if err != nil {
			log.Fatal(err)
		}
		columns = append(columns, name)
	}

	return columns
}

// SelectPens fetches all pens from the user's pens database.
func SelectPens(userID int64) ([]map[string]interface{}, []string, error) {
	// Get the path to the user's pens database file
	userDBPath := GetUserDBPath(userID)

	// Open the user's pens database
	userDB, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, nil, err
	}
	defer userDB.Close()

	// Construct the query to select all pens from the "pens" table
	query := "SELECT * FROM pens"

	// Fetch data from the user's pens database
	columns, pens, err := fetchDataFromDB(userDB, query)
	if err != nil {
		return nil, nil, err
	}

	return pens, columns, nil
}

// InsertPen inserts a new pen record into the database.
func InsertPen(userID int64, values []string) error {
	// Check if values have the necessary number of elements
	if len(values) < 1 {
		return errors.New("insufficient values for InsertPen")
	}

	// Get column names excluding "id"
	columns := GetColumnNames(userID, "pens")
	columns = columns[1:]

	// Check if the number of values matches the number of columns
	if len(columns) != len(values) {
		return errors.New("mismatched number of values for InsertPen")
	}

	valuePlaceholders := make([]string, len(columns))
	for i := range columns {
		valuePlaceholders[i] = "?"
	}

	insertQuery := fmt.Sprintf("INSERT INTO pens (%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(valuePlaceholders, ", "))

	// Open the user's pens database
	userDB, err := CreateOrUpdateUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	_, err = userDB.Exec(insertQuery, convertStringSliceToInterfaceSlice(values)...)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePen updates a pen record in the database.
func UpdatePen(userID int64, id int64, values []string) error {
	columns := GetColumnNames(userID, "pens")
	columns = columns[1:] // Exclude "id"

	setStatements := make([]string, len(columns))
	for i, col := range columns {
		setStatements[i] = fmt.Sprintf("%s = ?", col)
	}

	updateQuery := fmt.Sprintf("UPDATE pens SET %s WHERE id = ?", strings.Join(setStatements, ", "))

	// Open the user's pens database
	userDB, err := CreateOrUpdateUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	values = append(values, fmt.Sprintf("%d", id)) // Add the ID to the end of the values
	_, err = userDB.Exec(updateQuery, convertStringSliceToInterfaceSlice(values)...)
	if err != nil {
		return err
	}
	return nil
}

// GetPenByID retrieves a pen's data by its ID for a specific user.
func GetPenByID(userID int64, penID int64) (map[string]interface{}, error) {
	// Get the path to the user's pens database file
	userDBPath := GetUserDBPath(userID)

	// Open the user's pens database
	userDB, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, err
	}
	defer userDB.Close()

	query := fmt.Sprintf("SELECT * FROM pens WHERE id = %d", penID)
	_, pens, err := fetchDataFromDB(userDB, query)
	if err != nil {
		return nil, err
	}
	if len(pens) == 0 {
		return nil, fmt.Errorf("pen not found")
	}
	return pens[0], nil
}

// interfaceSlice converts []string to []interface{}.
func interfaceSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

// interfaceSliceToString converts []interface{} to []string.
func interfaceSliceToString(slice []interface{}) []string {
	stringSlice := make([]string, len(slice))
	for i, v := range slice {
		stringSlice[i] = fmt.Sprintf("%v", v)
	}
	return stringSlice
}

// convertStringSliceToInterfaceSlice converts []string to []interface{}.
func convertStringSliceToInterfaceSlice(slice []string) []interface{} {
	interfaceSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

// DeletePenByID deletes a pen from the database by its ID.
func DeletePenByID(id int64) error {
	// Construct the delete query
	deleteQuery := fmt.Sprintf("DELETE FROM pens WHERE id = ?")

	// Execute the delete query
	_, err := db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// InsertUser inserts a new user record into the database.
func InsertUser(username, firstName, middleName, lastName, email string, hashedPassword []byte, bio string) (int64, error) {
	// Construct the INSERT query for users
	insertUserQuery := `INSERT INTO users (username, first_name, middle_name, last_name, email, password, bio)
                        VALUES (?, ?, ?, ?, ?, ?, ?)`

	// Execute the INSERT query
	result, err := db.Exec(insertUserQuery, username, firstName, middleName, lastName, email, hashedPassword, bio)
	if err != nil {
		return 0, err
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// GetPasswordByUsername retrieves the hashed password from the database based on the username.
func GetPasswordByUsername(username string) ([]byte, error) {
	// Construct the SELECT query for users
	selectUserQuery := `SELECT password FROM users WHERE username = ?`

	// Execute the SELECT query
	row := db.QueryRow(selectUserQuery, username)

	// Initialize a variable to store the retrieved password
	var hashedPassword []byte

	// Scan the hashed password from the query result
	err := row.Scan(&hashedPassword)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// convertInterfaceToStringSlice converts a slice of interfaces to a slice of strings.
func convertInterfaceToStringSlice(slice []interface{}) []string {
	stringSlice := make([]string, len(slice))
	for i, v := range slice {
		stringSlice[i] = fmt.Sprintf("%v", v)
	}
	return stringSlice
}

// GetUserIDByUsername retrieves the user ID from the database based on the username.
func GetUserIDByUsername(username string) (int64, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return 0, err
	}

	// Construct the path to the main database file
	dbPath := filepath.Join(currentDir, "database", dbName)
	// Open the main database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Query the database to get the user ID based on the username
	query := "SELECT id FROM users WHERE username = ?"
	var userID int64
	err = db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
