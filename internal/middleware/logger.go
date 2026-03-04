package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware her isteği terminale yazar
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// İsteği asıl handler'a gönder
		next.ServeHTTP(w, r)

		// İşlem bittikten sonra logla
		log.Printf(
			"[%s] %s %s %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	}
}
