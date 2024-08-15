package middlewares

import (
	"net/http"
)

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		next.ServeHTTP(w, r)
	})
}
