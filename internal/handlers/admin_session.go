package handlers

import (
	"log/slog"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/internal/sessions"
)

// adminLoginHandler is the handler for the admin login page
func adminLoginHandler(w http.ResponseWriter, r *http.Request) {
	data := templateData(r)
	data["title"] = "Login"
	data["messages"] = flash.Get(w, r)
	render(w, data, false, "login")
}

// LoginPost handles the login form submission
func adminLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Try to authenticate the user
	user, err := services.AuthenticateUser(r.Context(), email, password)
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
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, helpers.URL("/admin"), http.StatusSeeOther)
}

// Logout destroys the user session
func adminLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r, "admin")
	if err != nil {
		log.Error("Error getting session: ", err)
		flash.Message{
			Title:   "Error",
			Message: "An error occurred while trying to log out.",
			Style:   flash.Error,
		}.Save(w, r)
		// Redirect to the login page
		http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
}

// adminRegisterHandler is the handler for the admin register page
func adminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "New user"

	render(w, data, false, "register")
}

// adminRegisterFormHandler handles the form submission for creating a new user
func adminRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "New user"

	r.ParseForm()
	var user models.User
	user.Name = r.Form.Get("name")
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")

	confirmPassword := r.Form.Get("password-confirm")

	err := services.CreateUser(r.Context(), &user, confirmPassword)
	if err != nil {
		slog.Error("Error creating user ", "err", err.Error())
		flash.Message{
			Title:   "Error",
			Message: err.Error(),
			Style:   flash.Error,
		}.Save(w, r)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	flash.Message{
		Title:   "Success",
		Message: "User created successfully. Please log in to continue.",
		Style:   flash.Success,
	}.Save(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
