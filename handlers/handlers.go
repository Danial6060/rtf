package handlers

import (
	"database/sql"
	"net/http"
)

func MakeHandler(db *sql.DB, handler func(db *sql.DB, w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(db, w, r)
	}
}
