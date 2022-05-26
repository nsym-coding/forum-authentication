package likes

import (
	"database/sql"
	"fmt"
	"log"
)

type Like struct {
	LikeID    int
	UserID    int
	PostID    int
	CommentID int
}

type Notifier struct {
	CID int
}

var (
	db      *sql.DB
	LastIns int64
	ID      Notifier
)

func PostCreatorID(db *sql.DB, postID int) int {

	rows, err := db.Query(`SELECT post.userID 
	FROM post 
	WHERE post.postID = ?
	;`, postID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err2 := rows.Scan(&ID.CID)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return ID.CID
}

//POST LIKE BUTTON
func LikeButton(db *sql.DB, userID int, postID int, creator int) {

	fmt.Println("------POSTCREATORID", creator)

	stmt, err := db.Prepare("INSERT INTO likes(userID, postID, notified, creatorID) VALUES(?, ?, 0, ?) ")
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

func CommentLikeButton(db *sql.DB, userID int, commentID int) {
	stmt, err := db.Prepare("INSERT INTO likes(userID, commentID) VALUES(?, ?) ")
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

func DeleteLike(db *sql.DB, userID int, postID int) {
	stmt, err := db.Prepare("DELETE FROM likes WHERE userID = ? AND postID = ?")
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

func DeleteCommentLike(db *sql.DB, userID int, commentID int) {
	stmt, err := db.Prepare("DELETE FROM likes WHERE userID = ? AND commentID = ?")
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


func GetPostLikes(db *sql.DB, postID int) int {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE postID = ?;", postID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func GetCommentLikes(db *sql.DB, commentID int) int {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM likes WHERE commentID = ?;", commentID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func GetNumComment(db *sql.DB, postID int) int {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM comments WHERE postID =?;", postID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

