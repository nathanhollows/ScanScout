package sessions

import (
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var store sessions.Store

func Start() {
	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
}

// Get returns a session for the given request
func Get(r *http.Request, name string) (*sessions.Session, error) {
	return store.Get(r, name)
}
