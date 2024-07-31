package data

import (
	"database/sql"
	"log"
	"net/http"
)



func GetUserIDFromSession(db *sql.DB, w http.ResponseWriter, r *http.Request) int {
    cookie, err := r.Cookie("session_token")
    if err != nil {
        if err == http.ErrNoCookie {
            log.Println("No session cookie found")
            w.WriteHeader(http.StatusUnauthorized)
            return -1
        }
        log.Println("Error retrieving session cookie:", err)
        w.WriteHeader(http.StatusBadRequest)
        return -1
    }

    sessionToken := cookie.Value
    var userID int
    err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", sessionToken).Scan(&userID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("Session not found for token:", sessionToken)
            w.WriteHeader(http.StatusUnauthorized)
            return -1
        }
        log.Println("Error querying session:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return -1
    }

    return userID
}

