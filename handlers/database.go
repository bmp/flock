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

	"golang.org/x/crypto/bcrypt"
)

// dbName is the name of the SQLite database file.
const dbName = "database.db"

var db *sql.DB

// InitDB initializes the database connection for handlers.
func InitDB(database *sql.DB) {
	db = database
}

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

		// Insert demo user
		demoUsername := "demo"
		demoPassword := "demo123"
		demoFirstName := "Demo"
		demoLastName := "User"
		demoEmail := "demo@example.com"
		demoBio := "This is a demo user."

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
		if err != nil {
			db.Close()
			return nil, err
		}

		_, err = db.Exec(`INSERT INTO users (username, first_name, last_name, email, password, bio)
			VALUES (?, ?, ?, ?, ?, ?)`, demoUsername, demoFirstName, demoLastName, demoEmail, hashedPassword, demoBio)
		if err != nil {
			db.Close()
			return nil, err
		}

		// Close the database connection after creating the demo user
		db.Close()
	}

	// Open the main database file
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// If the demo user was just created, also create the pens table for the demo user
	demoUserID, err := GetUserIDByUsername("demo")
	if err != nil {
		return nil, err
	}

	_, err = CreateOrUpdateUserDB(demoUserID)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Insert demo pens
	demoPens := [][]string{
		{"LAMY Safari", "LAMY", "Charcoal", "Plastic", "M", "Black", "Converter", "Silver", "2001-01-11", "30.00", "Smooth writer"},
		{"Pilot Metropolitan", "Pilot", "Silver", "Metal", "F", "Silver", "Cartridge", "Black", "2002-05-23", "18.00", "Classic design"},
		{"Pelikan Souverän M800", "Pelikan", "Green", "Resin", "F", "Gold", "Piston", "Gold", "2005-08-17", "600.00", "Timeless design"},
		{"Sailor Pro Gear", "Sailor", "Black", "Resin", "M", "Gold", "Converter", "Gold", "2008-11-30", "250.00", "Japanese craftsmanship"},
		{"Parker Duofold Centennial", "Parker", "Black", "Resin", "F", "Gold", "Converter", "Gold", "2010-03-02", "500.00", "Classic elegance"},
		{"Faber-Castell E-Motion", "Faber-Castell", "Pearwood", "Wood", "M", "Steel", "Converter", "Chrome", "2012-07-14", "150.00", "Unique wooden design"},
		{"Platinum 3776 Century", "Platinum", "Bourgogne", "Resin", "M", "Gold", "Converter", "Gold", "2014-10-05", "200.00", "Japanese precision"},
		{"Sheaffer Prelude", "Sheaffer", "Gunmetal", "Metal", "F", "Steel", "Converter", "Chrome", "2016-02-18", "80.00", "Sleek and modern"},
		{"Kaweco Sport", "Kaweco", "Classic Sport", "Plastic", "F", "Steel", "Cartridge", "Gold", "2018-04-21", "25.00", "Compact pocket pen"},
		{"Ranga Model 4", "Ranga", "Ebonite", "Ebonite", "B", "Steel", "Eyedropper", "Gold", "2020-09-10", "50.00", "Handmade Indian pen"},
		{"Deccan Advocate", "Deccan", "Red", "Acrylic", "M", "Steel", "Converter", "Chrome", "2021-12-03", "70.00", "Indian craftsmanship"},
		{"Guider Acrylic", "Guider", "Blue", "Acrylic", "F", "Steel", "Eyedropper", "Silver", "2022-06-14", "60.00", "Handmade Indian pen"},
		{"Ratnam Supreme", "Ratnam", "Green", "Ebonite", "UEF", "Gold", "Eyedropper", "Gold", "2003-09-27", "120.00", "Vintage Indian pen"},
		{"Bhramam Mystique", "Bhramam", "Purple", "Acrylic", "BBB", "Steel", "Converter", "Chrome", "2007-12-19", "90.00", "Artisan Indian pen"},
		{"Nakaya Piccolo Cigar", "Nakaya", "Kuro-Tamenuri", "Urushi", "EF", "Gold", "Converter", "Gold", "2011-04-07", "800.00", "Japanese Urushi masterpiece"},
		{"Hakase Fountain Pen", "Hakase", "Brown", "Ebonite", "Music", "Gold", "Piston", "Gold", "2015-08-29", "2000.00", "Custom handmade Japanese pen"},
		{"Conid Bulkfiller Regular", "Conid", "Black", "Resin", "Architect", "Gold", "Bulkfiller", "Gold", "2019-11-12", "700.00", "Innovative filling mechanism"},
		{"BCHR Waterman Ideal", "Waterman", "Black", "Hard Rubber", "Italic", "Gold", "Eyedropper", "Gold", "2023-01-15", "250.00", "Vintage BCHR pen"},
		{"Fosfor Islander", "Fosfor", "Blue", "Ebonite", "F", "Gold", "Vacuum", "Gold", "2004-06-09", "250.00", "Custom handmade pen with Vacuum system"},
	}

	for _, pen := range demoPens {
		err := insertDemoPen(demoUserID, pen)
		if err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

// insertDemoPen inserts a demo pen record into the user's database.
func insertDemoPen(userID int64, values []string) error {
	db, err := CreateOrUpdateUserDB(userID)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert demo pen into the "pens" table
	err = InsertPen(userID, values)
	if err != nil {
		return err
	}

	return nil
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
	// log.Printf("Opening existing user's pens database for user with ID %d...\n", userID)
	userDB, err := sql.Open("sqlite3", userDBPath)
	if err != nil {
		return nil, err
	}

	return userDB, nil
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

// PenExists checks if a pen with the given ID exists in the database for the given user.
func PenExists(userID, penID int64) bool {
	// Open the user's pens database
	userDB, err := CreateOrUpdateUserDB(userID)
	if err != nil {
		// log.Printf("Error opening user database: %s", err)
		return false
	}
	defer userDB.Close()

	// Query to check if the pen with the given ID exists for the user
	query := "SELECT COUNT(*) FROM pens WHERE id = ? AND user_id = ?"
	var count int
	err = userDB.QueryRow(query, penID, userID).Scan(&count)
	if err != nil {
		// log.Printf("Error checking pen existence: %s", err)
		return false
	}

	// If count is greater than 0, the pen exists; otherwise, it doesn't
	return count > 0
}

// DeletePenByID deletes a pen from the database by its ID.
func DeletePenByID(userID int64, id int64) error {
	// Open the user's pens database
	userDB, err := CreateOrUpdateUserDB(userID)
	if err != nil {
		return err
	}
	defer userDB.Close()

	// Construct the delete query
	deleteQuery := "DELETE FROM pens WHERE id = ?"

	// Execute the delete query
	_, err = userDB.Exec(deleteQuery, id)
	if err != nil {
		log.Printf("Error deleting pen with ID %d: %s", id, err)
		return err
	}

	return nil
}

// InsertUser inserts a new user record into the database.
func InsertUser(username, firstName, middleName, lastName, email string, hashedPassword []byte, bio string) (int64, error) {
	// Ensure the db variable is not nil
	if db == nil {
		return 0, errors.New("database connection not initialized")
	}

	// Construct the INSERT query for users
	insertUserQuery := `INSERT INTO users (username, first_name, middle_name, last_name, email, password, bio)
                        VALUES (?, ?, ?, ?, ?, ?, ?)`

	// log.Printf("Insert called with %s", insertUserQuery)

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
