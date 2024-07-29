package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func FetchChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, sender_id, receiver_id, content, created_at FROM private_messages ORDER BY created_at ASC")
	if err != nil {
		log.Printf("Error fetching chat history: %v", err)
		http.Error(w, "Error fetching chat history", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var chatHistory []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt); err != nil {
			log.Printf("Error scanning chat message: %v", err)
			http.Error(w, "Error scanning chat message", http.StatusInternalServerError)
			return
		}
		chatHistory = append(chatHistory, message)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chatHistory); err != nil {
		log.Printf("Error encoding chat history: %v", err)
		http.Error(w, "Error encoding chat history", http.StatusInternalServerError)
	}
}
