// handlers/register.go

package handlers

import (
	// "fmt"
	// "html/template"
  // "log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HandleRegister handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		// Extract user details from the form
		username := r.FormValue("username")
		firstName := r.FormValue("firstName")
		middleName := r.FormValue("middleName")
		lastName := r.FormValue("lastName")
		email := r.FormValue("email")
		password := r.FormValue("password")
		bio := r.FormValue("bio")

		// Hash the password before storing it
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		userID, err := InsertUser(username, firstName, middleName, lastName, email, hashedPassword, bio)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				// log.Println("Username or email already exists")
				http.Error(w, "Username or email already exists", http.StatusBadRequest)
				return
			}
			// log.Println("Error inserting user:", err)
			http.Error(w, "Error inserting user into the database", http.StatusInternalServerError)
			return
		}

		// Create or update the user's pens database
		_, err = CreateOrUpdateUserDB(userID)
		if err != nil {
			http.Error(w, "Error creating or updating user's pens database", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page or dashboard
		SetUserIDInSession(w, r, userID) // Set the user session
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
				return

	}

	// Generate a random CAPTCHA question
	question, answer := generateRandomCaptcha()

	// Render the registration form with CAPTCHA question
	data := struct {
		CaptchaQuestion string
		CaptchaAnswer string
	}{
		CaptchaQuestion: question,
		CaptchaAnswer: answer,
	}

	// Use renderTemplate when parsing templates
	renderTemplate(w, "register", data)
}

// generateRandomCaptcha generates a random CAPTCHA question and answer.
func generateRandomCaptcha() (question, answer string) {
	// List of 15 sample questions with corresponding answers
	questionsAndAnswers := map[string]string{
		"What is the purpose of the breather hole on a nib?":                                                           "regulation",
		"What is a common material for vintage pen bodies?":                                                            "ebonite",
		"What is the common name for the liquid used by fountain pen?":                                                 "ink",
		"Are demonstrators transparent or opaque?":                                                                     "opaque",
		"What is the common and cheap nib tipping material?":                                                           "iridium",
		"What is the name for a filling system which has an empty barrel?":                                             "eyedropper",
		"Which is more water-resistant ink: pigment or dye?":                                                           "pigment",
		"Why do you need a feed?":                                                                                      "regulation",
		"Which brand of pens has 'Safari' and 'Al-Star' pens?":                                                         "lamy",
		"What is the difference between flex and stub nib?":                                                            "flexibility",
		"Which amongst gold and steel nibs flexes more?":                                                               "gold",
		"What is phenomena when the nib flexes and only provides two parallel lines instead of a filled stroke?":       "railroad",
		"Which country has brands/companies such as Pilot, Platinum, etc.?":                                            "japan",
		"What is it called when you fix the cap of the pen to the barrel when you write?":                              "post",
		"What is the name of the part holding ink and connecting to the nib?":                                          "reservoir",
	}
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Select a random question
	randomIndex := rand.Intn(len(questionsAndAnswers))
	i := 0
	for q, a := range questionsAndAnswers {
		if i == randomIndex {
			question = q
			answer = a
			break
		}
		i++
	}

	// Use defer to ensure that the response is closed before leaving the function
	defer func() {
		if err := recover(); err != nil {
			// Handle the panic and log the error
			// log.Println("Panic in generateRandomCaptcha:", err)
		}
	}()

	return question, answer
}
