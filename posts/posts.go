package posts

import (
	"database/sql"
	"fmt"

	"forum/dislikes"
	"forum/likes"
)

type Post struct {
	PostID       int
	UserID       int
	CreationDate string
	PostTitle    string
	PostContent  string
	Image        string
	Edited       bool
}

type HomepagePosts struct {
	PostID       int
	PostTitle    string
	PostContent  string
	PostUsername string
	CreationDate string
	PostLike     int
	PostDislike  int
	NetLikes     int
	CommentNum   int
}

type ActPage struct {
	UserID       int
	PostTitle    string
	CommentID    int
	CommentText  string
	PostID       int
	CreationDate string
	PostLike     int
}

var db *sql.DB

// type s *web.Server

var LastIns int64

func CreatePosts(db *sql.DB, userID int, title string, content string, image string) {
	stmt, err := db.Prepare("INSERT INTO post (userID, postTitle, postContent, image, creationDate) VALUES (?, ?, ?, ?, strftime('%H:%M %d/%m/%Y','now','localtime'))")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	// defer stmt.Close()
	result, _ := stmt.Exec(userID, title, content, image)

	// checking if the result has been added and the last inserted row
	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted:", LastIns)
}

// Get all the data needed for the hompage
func GetHomepageData(db *sql.DB) []HomepagePosts {
	rows, err := db.Query(`SELECT postID, postTitle, postContent, username, creationDate FROM post 
	INNER JOIN users ON users.userID = post.userID;`)
	if err != nil {
		fmt.Println(err)
	}

	postdata := []HomepagePosts{}
	defer rows.Close()
	for rows.Next() {
		var p HomepagePosts
		// fmt.Println(&p.PostID)
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.PostContent, &p.PostUsername, &p.CreationDate)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		p.PostDislike = dislikes.GetPostDislikes(db, p.PostID)
		p.NetLikes = NetLikes(db, p.PostID)
		p.CommentNum = likes.GetNumComment(db, p.PostID)
		postdata = append(postdata, p)
		if err2 != nil {
			fmt.Println(err2)
		}
	}

	return postdata
}

func NetLikes(db *sql.DB, PostID int) int {
	NetLikes := likes.GetPostLikes(db, PostID) - dislikes.GetPostDislikes(db, PostID)

	return NetLikes
}

// func CommentLikes(db *sql.DB, PostID int) int {
// 	CommentLikes := likes.PCLikes(db, PostID)

// 	return CommentLikes
// }

// returns user's comments with their corresponding posts
func ActivityComments(db *sql.DB, userid int) []ActPage {
	rows, err := db.Query(`SELECT post.userID, post.postTitle, comments.commentID, comments.commentText, post.postID 
	FROM post, comments
	WHERE comments.userID = ?
	AND post.postID = comments.postID
	;`, userid)
	if err != nil {
		fmt.Println(err)
	}
	pac := []ActPage{}
	defer rows.Close()
	for rows.Next() {
		var p ActPage
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.CommentID, &p.CommentText, &p.PostID)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		pac = append(pac, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return pac
}

// returns user's liked posts
func ActivityPostLikes(db *sql.DB, userid int) []ActPage {
	rows, err := db.Query(`SELECT DISTINCT likes.userID, likes.postID, post.postTitle
	FROM likes, post
	WHERE likes.userID = ?
	AND likes.postID = post.postID 
	;`, userid)
	if err != nil {
		fmt.Println(err)
	}
	actlikes := []ActPage{}
	defer rows.Close()
	for rows.Next() {
		var p ActPage
		err2 := rows.Scan(&p.UserID, &p.PostID, &p.PostTitle)
		actlikes = append(actlikes, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return actlikes
}

// returns user's disliked posts
func ActivityPostDislikes(db *sql.DB, userid int) []ActPage {
	rows, err := db.Query(`SELECT DISTINCT dislikes.userID, dislikes.postID, post.postTitle
	FROM dislikes, post
	WHERE dislikes.userID = ?
	AND dislikes.postID = post.postID 
	;`, userid)
	if err != nil {
		fmt.Println(err)
	}
	actpostdislikes := []ActPage{}
	defer rows.Close()
	for rows.Next() {
		var p ActPage
		err2 := rows.Scan(&p.UserID, &p.PostID, &p.PostTitle)
		actpostdislikes = append(actpostdislikes, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return actpostdislikes
}

// returns user's liked comments
func ActivityCommentLikes(db *sql.DB, userid int) []ActPage {
	rows, err := db.Query(`SELECT DISTINCT likes.userID, likes.commentID, comments.commentText
	FROM likes, comments
	WHERE likes.userID = ?
	AND likes.commentID = comments.commentID 
	;`, userid)
	if err != nil {
		fmt.Println(err)
	}
	actcommentlikes := []ActPage{}
	defer rows.Close()
	for rows.Next() {
		var p ActPage
		err2 := rows.Scan(&p.UserID, &p.CommentID, &p.CommentText)
		actcommentlikes = append(actcommentlikes, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return actcommentlikes
}

// returns user's disliked comments
func ActivityCommentDislikes(db *sql.DB, userid int) []ActPage {
	rows, err := db.Query(`SELECT DISTINCT dislikes.userID, dislikes.commentID, comments.commentText
	FROM dislikes, comments
	WHERE dislikes.userID = ?
	AND dislikes.commentID = comments.commentID 
	;`, userid)
	if err != nil {
		fmt.Println(err)
	}
	actcommentdislikes := []ActPage{}
	defer rows.Close()
	for rows.Next() {
		var p ActPage
		err2 := rows.Scan(&p.UserID, &p.CommentID, &p.CommentText)
		actcommentdislikes = append(actcommentdislikes, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return actcommentdislikes
}

func UsersPostsHomepageData(db *sql.DB, userID int) []HomepagePosts {
	rows, err := db.Query("SELECT postID, postTitle, postContent, username, creationDate FROM post INNER JOIN users ON users.userID =  post.userID WHERE users.userID = ?;", userID)
	if err != nil {
		fmt.Println(err)
	}
	postdata := []HomepagePosts{}
	defer rows.Close()
	for rows.Next() {
		var p HomepagePosts
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.PostContent, &p.PostUsername, &p.CreationDate)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		p.PostDislike = dislikes.GetPostDislikes(db, p.PostID)
		p.NetLikes = NetLikes(db, p.PostID)
		p.CommentNum = likes.GetNumComment(db, p.PostID)
		postdata = append(postdata, p)
		if err2 != nil {
			fmt.Println(err2)
		}
	}

	return postdata
}

func UsersLikesHomepageData(db *sql.DB, userID int) []HomepagePosts {
	rows, err := db.Query(`SELECT DISTINCT post.postID, post.postTitle, post.creationDate, username
	FROM likes, post, users
	WHERE likes.userID = ?
	AND likes.postID = post.postID 
	AND post.userID = users.userID
	;`, userID)
	if err != nil {
		fmt.Println(err)
	}
	postdata := []HomepagePosts{}
	defer rows.Close()
	for rows.Next() {
		var p HomepagePosts
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.CreationDate, &p.PostUsername)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		p.PostDislike = dislikes.GetPostDislikes(db, p.PostID)
		p.NetLikes = NetLikes(db, p.PostID)
		p.CommentNum = likes.GetNumComment(db, p.PostID)
		postdata = append(postdata, p)
		if err2 != nil {
			fmt.Println(err2)
		}
	}

	return postdata
}

// getting the data from one post, and storing the values in the post struct
func GetPostData(db *sql.DB, postID int) []Post {
	rows, err := db.Query("SELECT * FROM post WHERE postID = ?;", postID)
	if err != nil {
		fmt.Println(err)
	}
	postdata := []Post{}
	defer rows.Close()
	for rows.Next() {
		var p Post
		err2 := rows.Scan(&p.PostID, &p.UserID, &p.CreationDate, &p.PostTitle, &p.PostContent, &p.Image, &p.Edited)
		postdata = append(postdata, p)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
	return postdata
}

// selects posts with the same category, based on the category name

func CategoryPagePosts(db *sql.DB, name string) []HomepagePosts {
	rows, err := db.Query(`SELECT post.postID, postTitle, username, creationDate FROM post 
	INNER JOIN postcategory ON postcategory.postID = post.postID 
	INNER JOIN users ON users.userID = post.userID
	WHERE postcategory.categoryname = ?;`, name)
	if err != nil {
		fmt.Println(err)
	}
	posts := []HomepagePosts{}
	defer rows.Close()
	for rows.Next() {
		var p HomepagePosts
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.PostUsername, &p.CreationDate)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		p.PostDislike = dislikes.GetPostDislikes(db, p.PostID)
		p.NetLikes = NetLikes(db, p.PostID)
		p.CommentNum = likes.GetNumComment(db, p.PostID)
		posts = append(posts, p)
		if err2 != nil {
			fmt.Println(err2)
		}

	}
	return posts
}

// func should delete post and all comments relating to post
func DeletePost(db *sql.DB, postID int) {
	// deleting the post
	stmt, err := db.Prepare("DELETE FROM post WHERE postID = ?")
	if err != nil {
		fmt.Println("error preparing delete post statement", err)
	}

	stmt.Exec(postID)

	// deleting the comments related to the post
	stmt2, err2 := db.Prepare("DELETE FROM comments WHERE postID = ?")

	if err2 != nil {
		fmt.Println("error deleting comments relating to post", err2)
	}
	stmt2.Exec(postID)

	// deleting the reports connected to a post
	stmt3, err3 := db.Prepare("DELETE FROM report WHERE postID = ?")
	if err3 != nil {
		fmt.Println("error deleting reports from the report table", err3)
	}
	stmt3.Exec(postID)

	stmt4, err4 := db.Prepare("DELETE FROM postcategory WHERE postID = ?")
	if err4 != nil {
		fmt.Println("error deleting post from the category table", err4)
	}
	stmt4.Exec(postID)
	// deleting the likes connected to a post
	stmt7, err7 := db.Prepare("DELETE FROM likes WHERE postID = ?")
	if err7 != nil {
		fmt.Println("error deleting reports from the likes table", err4)
	}
	stmt7.Exec(postID)

	// deleting the dislikes connected to a post
	stmt5, err5 := db.Prepare("DELETE FROM dislikes WHERE postID = ?")
	if err4 != nil {
		fmt.Println("error deleting reports from the dislikes table", err5)
	}
	stmt5.Exec(postID)

	// deleting the dislikes connected to a post
	stmt6, err6 := db.Prepare("DELETE FROM category WHERE postID = ?")
	if err3 != nil {
		fmt.Println("error deleting reports from the dislikes table", err6)
	}
	stmt6.Exec(postID)
}

func ReportedPostsHomepageData(db *sql.DB) []HomepagePosts {
	rows, err := db.Query(`SELECT post.postID, post.postTitle, users.username, post.creationDate 
						FROM post INNER JOIN report on report.postID = post.postID
						INNER JOIN users ON users.userID = post.userID`)
	if err != nil {
		fmt.Println(err)
	}
	postdata := []HomepagePosts{}
	defer rows.Close()
	for rows.Next() {
		var p HomepagePosts
		err2 := rows.Scan(&p.PostID, &p.PostTitle, &p.PostUsername, &p.CreationDate)
		p.PostLike = likes.GetPostLikes(db, p.PostID)
		p.PostDislike = dislikes.GetPostDislikes(db, p.PostID)
		p.NetLikes = NetLikes(db, p.PostID)
		p.CommentNum = likes.GetNumComment(db, p.PostID)
		postdata = append(postdata, p)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
	return postdata
}

func DenyReportRequest(db *sql.DB, postID int) {
	// delete the report
	stmt, err := db.Prepare("DELETE FROM report WHERE postID = ?")
	if err != nil {
		fmt.Println("error preparing delete report statement", err)
	}

	stmt.Exec(postID)
}
