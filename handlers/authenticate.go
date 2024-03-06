// handlers/authentication.go
package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

// Cookie store to manage session data
var store *sessions.CookieStore

// init initializes the cookie store with a random secret key.
func init() {
	// Generate a random secret key for the CookieStore
	secretKey, err := generateRandomKey(32) // Adjust the length based on your requirements
	if err != nil {
		log.Fatal("Error generating random key:", err)
	}

	// Create a new CookieStore with the generated secret key
	store = sessions.NewCookieStore([]byte(secretKey))
}

// SetUserIDInSession sets the user ID in the session.
func SetUserIDInSession(w http.ResponseWriter, r *http.Request, userID int64) {
    // Retrieve the session using the store
    session, _ := store.Get(r, "session-name")

    // Set the user ID in the session
    session.Values["userID"] = userID

    // Save the session to persist changes
    session.Save(r, w)
}

// GetUserIDFromSession retrieves the user ID from the session.
func GetUserIDFromSession(r *http.Request) int64 {
    // Retrieve the session using the store
    session, _ := store.Get(r, "session-name")

    // Check if the user ID is stored in the session
    if userID, ok := session.Values["userID"].(int64); ok {
        return userID
    }

    // Return 0 if user ID is not found
    return 0
}

// generateRandomKey generates a random key with the specified length.
func generateRandomKey(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert to hex string
	key := hex.EncodeToString(randomBytes)

	return key, nil
}
