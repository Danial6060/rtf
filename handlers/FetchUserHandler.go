package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func FetchUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT nickname FROM users")
	if err != nil {
		log.Printf("something went wrong: %v", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var nickname string
		if err := rows.Scan(&nickname); err != nil {
			log.Printf("something went wrong: %v", err)
			http.Error(w, "Error fetching users", http.StatusInternalServerError)
			return
		}
		users = append(users, nickname)
	}

	// if any errors where encounterd during the reading of the rows
	if err = rows.Err(); err != nil {
		log.Printf("something went wrong: %v", err)
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	// convert the array into json and send it
	if err = json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("error sending json: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
