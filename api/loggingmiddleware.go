package api

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логгирует детали запроса и тело, если это POST, PUT или DELETE запрос.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			log.Printf("Method: %s, Path: %s, Body: %s", r.Method, r.URL.Path, string(body))
		}

		next.ServeHTTP(w, r)

		log.Printf("Processed request: %s %s in %v", r.Method, r.URL.Path, time.Since(startTime))
	})
}
