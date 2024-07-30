package handlers

import (
	"database/sql"
	"log"
	"net/http"
	data "rtf/Data"
)

func CommentPostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := data.GetUserIDFromSession(w, r)
	if userID == -1 {
		log.Println("User not authenticated")
		return
	}

	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	if err != nil {
		log.Printf("Error commenting post: %v", err)
		http.Error(w, "Error commenting post", http.StatusInternalServerError)
		return
	}

	log.Println("Comment created successfully by user:", userID)
	w.WriteHeader(http.StatusCreated)
}
