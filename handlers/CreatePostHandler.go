package handlers

import (
	"database/sql"
	"log"
	"net/http"
	data "rtf/Data"
)

func CreatePostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    userID := data.GetUserIDFromSession(db, w, r) // Pass db here
    if userID == -1 {
        log.Println("User not authenticated")
        return
    }

    category := r.FormValue("category")
    content := r.FormValue("content")

    _, err := db.Exec("INSERT INTO posts (user_id, category, content) VALUES (?, ?, ?)", userID, category, content)
    if err != nil {
        log.Printf("Error creating post: %v", err)
        http.Error(w, "Error creating post", http.StatusInternalServerError)
        return
    }

    log.Println("Post created successfully by user:", userID)
    w.WriteHeader(http.StatusCreated)
}
