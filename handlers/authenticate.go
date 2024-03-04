// handlers/authentication.go
package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

var store *sessions.CookieStore


func init() {
	// Generate a random secret key for the Cookiestore
	secretKey, err := generateRandomKey(32) // Adjust the length based on your requirements
	if err != nil {
		log.Fatal("Error generating random key:", err)
	}

	// Create a new CookieStore with the generated secret key
	store = sessions.NewCookieStore([]byte(secretKey))
}

// SetUserIDInSession sets the user ID in the session.
func SetUserIDInSession(w http.ResponseWriter, r *http.Request, userID int64) {
    session, _ := store.Get(r, "session-name")
    session.Values["userID"] = userID
    session.Save(r, w)
}

// GetUserFromSession retrieves the user ID from the session.
func GetUserIDFromSession(r *http.Request) int64 {
    session, _ := store.Get(r, "session-name")
    if userID, ok := session.Values["userID"].(int64); ok {
        return userID
    }
    return 0
}

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
