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

	// Pass db to handlers
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/register", handlers.MakeHandler(db, handlers.RegisterHandler))
	http.HandleFunc("/login", handlers.MakeHandler(db, handlers.LoginHandler))
	http.HandleFunc("/logout", handlers.MakeHandler(db, handlers.LogoutHandler))
	http.HandleFunc("/create_post", handlers.MakeHandler(db, handlers.CreatePostHandler))
	http.HandleFunc("/comment_post", handlers.MakeHandler(db, handlers.CommentPostHandler))
	http.HandleFunc("/fetch_posts", handlers.MakeHandler(db, handlers.FetchPostsHandler))
	http.HandleFunc("/fetch_comments", handlers.MakeHandler(db, handlers.FetchCommentsHandler))
	http.HandleFunc("/is_logged_in", handlers.MakeHandler(db, handlers.IsUserLoggedInHandler))
	http.HandleFunc("/fetch_users", handlers.MakeHandler(db, handlers.FetchUserHandler))
	http.HandleFunc("/fetch_chat_history", handlers.MakeHandler(db, handlers.FetchChatHistoryHandler))
	http.HandleFunc("/ws", handlers.MakeHandler(db, handlers.WSHandler))

	fmt.Printf("Server started at %s", Port)
	http.ListenAndServe(Port, nil)
}
