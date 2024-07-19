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

// LoginHandler is the handler for the admin login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := TemplateData(r)
	data["title"] = "Login"
	data["messages"] = flash.Get(w, r)
	Render(w, data, false, "login")
}

// LoginPost handles the login form submission
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Try to authenticate the user
	user, err := services.AuthenticateUser(r.Context(), email, password)
	if err != nil {
		log.Error("Error authenticating user", "err", err)
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
		log.Error("getting session for login: ", err)
		flash.NewError("An error occurred while trying to log in. Please try again.").Save(w, r)
		// Redirect to the login page
		http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
		return
	}
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, helpers.URL("/admin"), http.StatusSeeOther)
}

// Logout destroys the user session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r, "admin")
	if err != nil {
		log.Error("getting session for logout: ", err)
		flash.NewError("An error occurred while trying to log out. Please try again.").Save(w, r)
		// Redirect to the login page
		http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
		return
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, helpers.URL("/login"), http.StatusSeeOther)
}

// RegisterHandler is the handler for the admin register page
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	data := TemplateData(r)
	data["title"] = "New user"

	Render(w, data, false, "register")
}

// RegisterPostHandler handles the form submission for creating a new user
func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	SetDefaultHeaders(w)
	data := TemplateData(r)
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
		flash.NewError("Error creating user. Please try again.").Save(w, r)
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	flash.NewSuccess("User created successfully. Please log in to continue.").Save(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
