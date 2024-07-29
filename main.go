package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	data "rtf/Data"
	"rtf/handlers"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Port   string = ":8080"
	dbPath string = "./forum.db"
)

func main() {

	data.DbExists(dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/create_post", handlers.CreatePostHandler)
	http.HandleFunc("/comment_post", handlers.CommentPostHandler)
	http.HandleFunc("/fetch_posts", handlers.FetchPostsHandler)
	http.HandleFunc("/fetch_comments", handlers.FetchCommentsHandler)
	http.HandleFunc("/is_logged_in", handlers.IsUserLoggedInHandler)
	http.HandleFunc("/fetch_users", handlers.FetchUserHandler)
	http.HandleFunc("/fetch_chat_history", handlers.FetchChatHistoryHandler)
	http.HandleFunc("/ws", handlers.WSHandler)

	fmt.Printf("Server started at %s", Port)
	http.ListenAndServe(Port, nil)
}
