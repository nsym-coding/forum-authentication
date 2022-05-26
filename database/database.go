package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "forum.db"

var PostIDInt int

func CreateDB() {
	var db *sql.DB
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("create table if not exists users (userID integer primary key AUTOINCREMENT, email text, username text, hash CHAR(60), usertype text, externalloginid text, becomemod INTEGER)")
	db.Exec(`create table if not exists post (
		postID integer primary key AUTOINCREMENT, 
		userID integer REFERENCES users(userID), 
		creationDate integer,
		postTitle CHAR(50),
		postContent CHAR(250), 
		image CHAR(100), 
		edited integer);`)
	db.Exec("create table if not exists postcategory (categoryID integer PRIMARY KEY AUTOINCREMENT, postID integer REFERENCES post(postID), categoryname text)")
	db.Exec("create table if not exists categories (categoryID integer PRIMARY KEY AUTOINCREMENT, categoryname text)")
	db.Exec(`create table if not exists comments (
		commentID integer primary key AUTOINCREMENT, 
		userID integer REFERENCES users(userID), 
		postID integer REFERENCES post(postID), 
		commentText CHAR(250), 
		edited integer, 
		creationDate integer,
		notified integer,
		creatorID integer);`)
	db.Exec(`create table if not exists likes (likeID integer PRIMARY KEY AUTOINCREMENT, 
		userID integer REFERENCES users(userID), postID integer REFERENCES post(postID), 
		commentID integer REFERENCES comments(commentID),
		notified integer,
		creatorID integer);`)
	db.Exec(`create table if not exists dislikes (dislikeID integer PRIMARY KEY AUTOINCREMENT, 
			userID integer REFERENCES users(userID), postID integer REFERENCES post(postID), 
			commentID integer REFERENCES comments(commentID),
			notified integer,
			creatorID integer);`)

	db.Exec(`create table if not exists report(reportID integer PRIMARY KEY AUTOINCREMENT, reporttype text,
				username text REFERENCES users(username), postID integer REFERENCES post(postID));`)
}

//
