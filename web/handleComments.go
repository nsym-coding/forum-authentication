package web

import (
	"fmt"
	"forum/comments"
	"forum/likes"
	"net/http"
	"strconv"
)

type ActPage struct {
	PostID       int
	PostTitle    string
	CommentText  string
	CreationDate string
	PostLike     int
}

var CommentData comments.Comment
var CUserIdint int
var ContentComment string
var CommentId int

//var CPostIdint int

func (s *myServer) CreateCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		userID := r.URL.Query().Get("userid")

		PostIDInt, _ = strconv.Atoi(userID)

		Tpl.ExecuteTemplate(w, "createcomment.html", nil)

	}
}

func (s *myServer) StoreCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		ContentComment = r.FormValue("content")

		CreatorsID := likes.PostCreatorID(s.Db, PostIDInt)

		comments.CreateComment(s.Db, GuserId, PostIDInt, ContentComment, CreatorsID)

		var commentData comments.Comment
		SPostID := strconv.Itoa(PostIDInt)

		commentData.CommentText = ContentComment
		commentData.PostID = PostIDInt
		commentData.UserID = GuserId

		CommentId = comments.GetCommentID(s.Db, PostIDInt)

		fmt.Println("testing method to get comment id", CommentId)

		http.Redirect(w, r, "/showpost/?postid="+SPostID, http.StatusSeeOther)

	}
}

func (s *myServer) DeleteCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, id := range v {
				idInt, _ := strconv.Atoi(id)
				comments.DeleteComment(s.Db, idInt)
			}
		}
		SPostID := strconv.Itoa(PostIDInt)
		http.Redirect(w, r, "/showpost/?postid="+SPostID, http.StatusSeeOther)
	}
}
