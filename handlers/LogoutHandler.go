package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

func LogoutHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value

	_, err = db.Exec("DELETE FROM sessions WHERE token = ?", sessionToken)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	log.Println("User logged out successfully")
	w.WriteHeader(http.StatusOK)
}
