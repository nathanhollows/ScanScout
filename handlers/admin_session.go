package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/ScanScout/flash"
	"github.com/nathanhollows/ScanScout/helpers"
	"github.com/nathanhollows/ScanScout/models"
	"github.com/nathanhollows/ScanScout/sessions"
)

// adminLoginHandler is the handler for the admin login page
func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	data["title"] = "Login"
	data["messages"] = flash.Get(w, r)
	render(w, data, false, "login")
}

// LoginPost handles the login form submission
func adminLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Try to authenticate the user
	user, err := models.AuthenticateUser(email, password)
	if err != nil {
		log.Error("Error authenticating user: ", err)
		flash.Message{
			Style:   flash.Error,
			Title:   "Invalid email or password",
			Message: "Please check your email and password and try again.",
		}.Save(w, r)
		http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
		return
	}

	// Create a session
	session, err := sessions.Get(r, "admin")
	if err != nil {
		log.Error("Error getting session: ", err)
		flash.Message{
			Title:   "Error",
			Message: "An error occurred while trying to log in.",
			Style:   flash.Error,
		}.Save(w, r)
		// Redirect to the login page
		http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
		return
	}
	session.Values["user_id"] = user.UserID
	session.Save(r, w)

	http.Redirect(w, r, helpers.URL("/admin"), http.StatusSeeOther)

}

// Logout destroys the user session
func Logout(w http.ResponseWriter, r *http.Request) error {
	return nil
}
