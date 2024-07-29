package handlers

import (
	"database/sql"
	"net/http"
)

func IsUserLoggedInHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "User is unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Error retrieving token", http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value
	var token string
	// check if the token exists on the database
	err = db.QueryRow("SELECT token FROM sessions WHERE token = ?", sessionToken).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// user is authorized
	w.WriteHeader(http.StatusOK)
}
