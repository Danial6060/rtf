package handlers

import (
	"log"
	"net/http"
	"time"
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade:", err)
		return
	}

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read:", err)
			break
		}

		_, err = db.Exec("INSERT INTO private_messages (sender_id, receiver_id, content, created_at) VALUES ($1, $2, $3, $4)", msg.SenderID, msg.ReceiverID, msg.Content, time.Now())
		if err != nil {
			log.Println("Error inserting message:", err)
			break
		}

		msg.CreatedAt = time.Now()
		err = conn.WriteJSON(msg)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
	defer conn.Close()

}
