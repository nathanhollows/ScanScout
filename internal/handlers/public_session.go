package handlers

import (
	"log/slog"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/flash"
	"github.com/nathanhollows/Rapua/internal/models"
	"github.com/nathanhollows/Rapua/internal/services"
)

func adminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	setDefaultHeaders(w)
	data := templateData(r)
	data["title"] = "New user"

	render(w, data, false, "register")
}

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
