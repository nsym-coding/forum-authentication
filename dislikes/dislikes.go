package dislikes

import (
	"database/sql"
	"fmt"
	"log"
)

type Disike struct {
	DislikeID int
	UserID    int
	PostID    int
	CommentID int
}

var db *sql.DB

var LastIns int64

func DislikeButton(db *sql.DB, userID int, postID int, creator int) {
	stmt, err := db.Prepare("INSERT INTO dislikes(userID, postID, notified, creatorID) VALUES(?, ?, 0, ?) ")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}

	result, _ := stmt.Exec(userID, postID, creator)

	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted", LastIns)
}

func CommentDislikeButton(db *sql.DB, userID int, commentID int) {
	stmt, err := db.Prepare("INSERT INTO dislikes(userID, commentID) VALUES(?, ?) ")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}

	result, _ := stmt.Exec(userID, commentID)

	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted", LastIns)
}

func DeleteDislike(db *sql.DB, userID int, postID int) {
	stmt, err := db.Prepare("DELETE FROM dislikes WHERE userID = ? AND postID = ?")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	result, _ := stmt.Exec(userID, postID)

	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted", LastIns)
}

func DeleteCommentDislike(db *sql.DB, userID int, commentID int) {
	stmt, err := db.Prepare("DELETE FROM dislikes WHERE userID = ? AND commentID = ?")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	result, _ := stmt.Exec(userID, commentID)

	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted", LastIns)
}

func GetPostDislikes(db *sql.DB, postID int) int {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE postID = ?;", postID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func GetCommentDislikes(db *sql.DB, commentID int) int {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE commentID = ?;", commentID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}
