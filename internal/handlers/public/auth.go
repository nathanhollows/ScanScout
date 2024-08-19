package handlers

import (
	"errors"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/helpers"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
	"github.com/nathanhollows/Rapua/internal/sessions"
	templates "github.com/nathanhollows/Rapua/internal/templates/public"
)

// LoginHandler is the handler for the admin login page
func (h *PublicHandler) Login(w http.ResponseWriter, r *http.Request) {
	c := templates.Login()
	err := templates.AuthLayout(c, "Login").Render(r.Context(), w)

	if err != nil {
		h.Logger.Error("Error rendering login page", "err", err)
	}
}

// LoginPost handles the login form submission
func (h *PublicHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Try to authenticate the user
	authService := services.NewAuthService()
	user, err := authService.AuthenticateUser(r.Context(), email, password)

	if err != nil {
		h.Logger.Error("authenticating user", "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		c := templates.LoginError("Invalid email or password.")
		c.Render(r.Context(), w)
		return
	}

	// Create a session
	session, err := sessions.Get(r, "admin")
	if err != nil {
		h.Logger.Error("getting session for login", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		c := templates.LoginError("An error occurred while trying to log in. Please try again.")
		c.Render(r.Context(), w)
		return
	}

	session.Values["user_id"] = user.ID
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.Save(r, w)

	w.Header().Add("hx-redirect", "/admin")
}

// Logout destroys the user session
func (h *PublicHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Get(r, "admin")
	if err != nil {
		h.Logger.Error("getting session for logout", "err", err)
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
func (h *PublicHandler) Register(w http.ResponseWriter, r *http.Request) {
	c := templates.Register()
	err := templates.AuthLayout(c, "Register").Render(r.Context(), w)

	if err != nil {
		h.Logger.Error("rendering register page", "err", err)
	}
}

// RegisterPostHandler handles the form submission for creating a new user
func (h *PublicHandler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var user models.User
	user.Name = r.Form.Get("name")
	user.Email = r.Form.Get("email")
	user.Password = r.Form.Get("password")

	confirmPassword := r.Form.Get("password-confirm")

	err := services.CreateUser(r.Context(), &user, confirmPassword)
	if err != nil {
		h.Logger.Error("creating user", "err", err)
		w.WriteHeader(http.StatusUnauthorized)
		if errors.Is(err, services.ErrPasswordsDoNotMatch) {
			c := templates.RegisterError("Passwords do not match.")
			c.Render(r.Context(), w)
			return
		}
		c := templates.RegisterError("Something went wrong! Please try again.")
		c.Render(r.Context(), w)
		return
	}

	c := templates.RegisterSuccess()
	c.Render(r.Context(), w)
}

// ForgotPasswordHandler is the handler for the forgot password page
func (h *PublicHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	c := templates.ForgotPassword()
	err := templates.AuthLayout(c, "Forgot Password").Render(r.Context(), w)

	if err != nil {
		h.Logger.Error("rendering forgot password page", "err", err)
	}
}

// ForgotPasswordPostHandler handles the form submission for the forgot password page
func (h *PublicHandler) ForgotPasswordPost(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	c := templates.ForgotMessage(
		*flash.NewInfo("If an account with that email exists, an email will be sent with instructions on how to reset your password."),
	)
	c.Render(r.Context(), w)
}
