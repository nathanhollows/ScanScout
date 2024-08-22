package sessions

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	gsessions "github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/nathanhollows/Rapua/internal/models"
)

var store sessions.Store

func Start() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	authStore := gsessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	authStore.Options.SameSite = http.SameSiteLaxMode
	authStore.Options.HttpOnly = true
	authStore.Options.Secure = true
	gothic.Store = authStore
	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_SECRET_ID"),
			fmt.Sprintf("%s/auth/google/callback", os.Getenv("SITE_URL")),
			"email",
			"profile",
		),
	)

}

// Get returns a session for the given request
func Get(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}

// New session for the given request and user
func New(r *http.Request, user models.User) (*sessions.Session, error) {
	session, err := store.Get(r, "admin")
	if err != nil {
		return nil, err
	}

	session.Values["user_id"] = user.ID
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode

	return session, nil
}
