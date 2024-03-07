// helper_functions.go

package handlers

import (
	"net/http"
	"strings"
	"html/template"
	"fmt"
	// "path/filepath"
	// "log"
)

var templates = make(map[string]*template.Template)

// Add function adds a number to the input
func Add(input, numberToAdd int) int {
	return input + numberToAdd
}


// Title capitalizes text and replaces underscores with space.
func Title(text string) string {
	text = strings.ReplaceAll(text, "_", " ")

	// Capitalize each word
	words := strings.Fields(text)
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	return strings.Join(words, " ")
}

// RegisterTemplate registers a template with a given name
func RegisterTemplate(name string, tmpl *template.Template) {
	templates[name] = tmpl
}

// GetTemplate retrieves a registered template by name
func GetTemplate(name string) (*template.Template, bool) {
	tmpl, ok := templates[name]
	return tmpl, ok
}

// renderTemplate renders an HTML template.
func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {

	// Log the data
	// log.Printf("Rendering template %s with data: %+v", templateName, data)

	// Parse the template files
	tmplFiles := fmt.Sprintf("templates/%s.html", templateName)
	tmpl, err := template.ParseFiles(tmplFiles)
	if err != nil {
		http.Error(w, "Error parsing template file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// RedirectWithError redirects to the specified URL with an error message.
func RedirectWithError(w http.ResponseWriter, r *http.Request, targetURL, errorMessage string) {
	redirectURL := fmt.Sprintf("%s?error=%s", targetURL, errorMessage)
	// log.Printf("errorMessage is %s", errorMessage)
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
