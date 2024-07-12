package flash

import (
	"encoding/gob"
	"net/http"

	"github.com/nathanhollows/Rapua/internal/sessions"
)

func init() {
	gob.Register(Message{})
}

// Message is a struct containing each flashed message
type Message struct {
	Title   string
	Message string
	Style   FlashStyle
}

type FlashStyle string

const (
	Default FlashStyle = ""
	Success FlashStyle = "success"
	Error   FlashStyle = "error"
	Warning FlashStyle = "warning"
	Info    FlashStyle = "info"
)

// New adds a new message into the cookie storage.
func New(w http.ResponseWriter, r *http.Request, title string, message string, style FlashStyle) {
	flash := Message{Title: title, Message: message, Style: style}
	session, _ := sessions.Get(r, "scanscout")
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.AddFlash(flash)
	session.Save(r, w)
}

// Save adds a new message into the cookie storage.
func (m Message) Save(w http.ResponseWriter, r *http.Request) {
	session, _ := sessions.Get(r, "scanscout")
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteStrictMode
	session.AddFlash(m)
	session.Save(r, w)
}

// Get flash messages from the cookie storage.
func Get(w http.ResponseWriter, r *http.Request) []interface{} {
	session, err := sessions.Get(r, "scanscout")
	if err == nil {
		messages := session.Flashes()
		if len(messages) > 0 {
			session.Save(r, w)
		}
		return messages
	}
	return nil
}
