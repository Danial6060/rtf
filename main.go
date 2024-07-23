package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	db       *sql.DB
	upgrader = websocket.Upgrader{}
)

type User struct {
	ID        int    `json:"id"`
	Nickname  string `json:"nickname"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Category  string    `json:"category"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Nickname  string    `json:"nickname"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Nickname  string    `json:"nickname"`
}

type Message struct {
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/create_post", createPostHandler)
	http.HandleFunc("/comment_post", commentPostHandler)
	http.HandleFunc("/fetch_posts", fetchPostsHandler)
	http.HandleFunc("/fetch_comments", fetchCommentsHandler)
	http.HandleFunc("/is_logged_in", isUserLoggedInHandler)
	http.HandleFunc("/fetch_users", FetchUserHandler)
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var user *User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	result, err := db.Exec("INSERT INTO users (nickname, age, gender, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, hashedPassword)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	// get the last inserted id from the database
	userID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	sessionToken := uuid.New().String()
	_, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(24 * time.Hour / time.Second), // max age will be 24 hours for the cookie
	})

	log.Println("User registered successfully:", user.Nickname)
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	var hashedPassword string
	var userID int
	err := db.QueryRow("SELECT id, password FROM users WHERE nickname = ? OR email = ?", identifier, identifier).Scan(&userID, &hashedPassword)
	if err != nil {
		log.Printf("Error querying user: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Password mismatch: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.New().String()
	_, err = db.Exec("INSERT INTO sessions (user_id, token) VALUES (?, ?)", userID, sessionToken)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(24 * time.Hour / time.Second), // max age will be 24 hours for the cookie
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
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

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromSession(w, r)
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

func commentPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromSession(w, r)
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

func fetchPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
        SELECT p.id, p.category, p.content, p.created_at, u.nickname
        FROM posts p
        JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `)
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Category, &post.Content, &post.CreatedAt, &post.Nickname); err != nil {
			log.Printf("Error scanning post: %v", err)
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding posts: %v", err)
		http.Error(w, "Error encoding posts", http.StatusInternalServerError)
	}
}

func fetchCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("post_id")
	rows, err := db.Query(`
        SELECT c.id, c.post_id, c.content, c.created_at, u.nickname
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at ASC
    `, postID)
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.Content, &comment.CreatedAt, &comment.Nickname); err != nil {
			log.Printf("Error scanning comment: %v", err)
			http.Error(w, "Error scanning comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error with rows: %v", err)
		http.Error(w, "Error with rows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		log.Printf("Error encoding comments: %v", err)
		http.Error(w, "Error encoding comments", http.StatusInternalServerError)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("Read:", err)
			break
		}

		_, err = db.Exec("INSERT INTO private_messages (sender_id, receiver_id, content) VALUES (?, ?, ?)", msg.SenderID, msg.ReceiverID, msg.Content)
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
}

func getUserIDFromSession(w http.ResponseWriter, r *http.Request) int {
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

func isUserLoggedInHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "User is unauthorized", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Error retrieving token", http.StatusBadRequest)
		return
	}

	sessionToken := cookie.Value
	var token string
	// check if the token exists on the database
	err = db.QueryRow("SELECT token FROM sessions WHERE token = ?", sessionToken).Scan(&token)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// user is authorized
	w.WriteHeader(http.StatusOK)
}
