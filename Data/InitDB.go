package data

import (
	"database/sql"
	"log"
)

func InitDB(db *sql.DB) {
	// Create User Table
	createUserTable := `CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nickname TEXT UNIQUE,
    age INTEGER,
    gender TEXT,
    first_name TEXT,
    last_name TEXT,
    email TEXT UNIQUE,
    password TEXT
);`
	_, err := db.Exec(createUserTable)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return
	}

	// Create Post Table
	createTablePosts := `CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    category TEXT,
    content TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err = db.Exec(createTablePosts)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return
	}

	// Create sessions Table
	createTablesessions := `CREATE TABLE sessions (
    		id INTEGER PRIMARY KEY AUTOINCREMENT,
    		user_id INTEGER,
    		token TEXT,
   		 FOREIGN KEY(user_id) REFERENCES users(id)
		);`
	_, err = db.Exec(createTablesessions)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return
	}

	// Create Comment Table
	createTableComments := `CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    content TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err = db.Exec(createTableComments)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return
	}
	// Create Private Messages Table
	createTablePrivateMessages := `CREATE TABLE private_messages (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	sender_id INTEGER,
    	receiver_id INTEGER,
   		content TEXT,
    	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY(sender_id) REFERENCES users(id),
    	FOREIGN KEY(receiver_id) REFERENCES users(id)
	);`
	_, err = db.Exec(createTablePrivateMessages)
	if err != nil {
		log.Fatal("Failed to create table:", err)
		return
	}
	// Create Online Users Table
	// createTableOnlineUsers := `CREATE TABLE online_users (
	// 	user_id INTEGER PRIMARY KEY,
	// 	is_online BOOLEAN DEFAULT FALSE,
	// 	last_activity DATETIME DEFAULT CURRENT_TIMESTAMP,
	// 	FOREIGN KEY (user_id) REFERENCES users(id)
	// );`
	// _, err = db.Exec(createTableOnlineUsers)
	// if err != nil {
	// 	log.Fatal("Failed to create table:", err)
	// 	return
	// }

}
