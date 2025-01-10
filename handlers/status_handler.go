package handlers

import (
	"database/sql"
	"net/http"
)

func ReadinessHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the database is reachable
		err := db.Ping()
		if err != nil {
			http.Error(w, "Database not connected", http.StatusInternalServerError)
			return
		}

		// Write response and check for error
		if _, err := w.Write([]byte("ready")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
	}
}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	// Write response and check for error
	if _, err := w.Write([]byte("alive")); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
