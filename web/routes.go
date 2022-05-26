package web

import (
	"fmt"
	"net/http"
	"time"
)

var Requests []*http.Request

func Rate(a <-chan time.Time, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Requests = append(Requests, r)

		if len(Requests) > 0 {
			func() {
				<-a
			}()

			next.ServeHTTP(w, Requests[0])
			Requests = Requests[1:]

		}

		fmt.Println()
	}
}

func (s *myServer) Routes(a <-chan time.Time) {
	// http.HandleFunc("/register", srv.LoginAuthHandler())
	http.HandleFunc("/home", Rate(a, HomepageSessionChecker(s.HomepageHandler())))
	http.HandleFunc("/register/", Rate(a, s.RegisterUserHandler()))
	http.HandleFunc("/registerauth", Rate(a, s.RegisterAuthHandler()))
	http.HandleFunc("/login", Rate(a, s.LoginHandler()))
	http.HandleFunc("/githublogin", Rate(a, s.GitLogin()))
	http.HandleFunc("/oauth/redirect/", Rate(a, s.GitOAUTHLogin()))
	http.HandleFunc("/googlelogin", Rate(a, s.GoogleLogin()))
	http.HandleFunc("/google/redirect", Rate(a, s.GoogleOAUTHLogin()))
	http.HandleFunc("/activitypage", Rate(a, Auth(s.ActivityPage())))
	http.HandleFunc("/loginauth", Rate(a, s.LoginAuthHandler()))
	http.HandleFunc("/logout", Rate(a, s.LogoutHandler()))
	http.HandleFunc("/createpost/", Rate(a, Auth(SessionChecker(s.CreatePostHandler()))))
	http.HandleFunc("/storepost", Rate(a, Auth(SessionChecker(s.StorePostHandler()))))
	http.HandleFunc("/createcomment/", Rate(a, Auth(SessionChecker(s.CreateCommentHandler()))))
	http.HandleFunc("/storecomment", Rate(a, Auth(SessionChecker(s.StoreCommentHandler()))))
	http.HandleFunc("/likes", Rate(a, Auth(SessionChecker(s.LikeHandler()))))
	http.HandleFunc("/delete", Rate(a, Auth(SessionChecker(s.DeletePost()))))
	http.HandleFunc("/report", Rate(a, Auth(SessionChecker(s.ReportHandler()))))
	http.HandleFunc("/deletecomment", Rate(a, Auth(SessionChecker(s.DeleteCommentHandler()))))
	http.HandleFunc("/dislikes", Rate(a, Auth(SessionChecker(s.DislikeHandler()))))
	http.HandleFunc("/commentlikes", Rate(a, Auth(SessionChecker(s.CommentLikeHandler()))))
	http.HandleFunc("/commentdislikes", Rate(a, Auth(SessionChecker(s.CommentDislikeHandler()))))
	http.HandleFunc("/showpost/", Rate(a, s.ShowPostHandler()))
	http.HandleFunc("/becomeamod", Rate(a, Auth(SessionChecker(s.BecomeAModHandler()))))
	http.HandleFunc("/acceptmod", Rate(a, Auth(SessionChecker(s.AcceptAModHandler()))))
	http.HandleFunc("/declinemod", Rate(a, Auth(SessionChecker(s.DeclineAModHandler()))))
	http.HandleFunc("/demotemod", Rate(a, Auth(SessionChecker(s.DemoteAModHandler()))))
	http.HandleFunc("/addcategory", Rate(a, Auth(SessionChecker(s.AdminAddingCategory()))))
	http.HandleFunc("/deletecategory", Rate(a, Auth(SessionChecker(s.AdminDeletingCategory()))))

	// http.HandleFunc("/emptycommentpost/", Rate(a, Auth(SessionChecker(s.EmptyCommentPost()))))
	// http.HandleFunc("/showcomment/", Rate(a, Auth(SessionChecker(s.ShowCommentHandler()))))
	// this is for the template css files to run.
	http.HandleFunc("/deleteactpost", Rate(a, Auth(SessionChecker(s.DeleteActPost()))))
	http.HandleFunc("/deleteactcomment", Rate(a, Auth(SessionChecker(s.DeleteActComment()))))
	http.HandleFunc("/editactpost", Rate(a, Auth(SessionChecker(s.EditActComment()))))
	http.HandleFunc("/postedited", Rate(a, Auth(SessionChecker(s.EditPostHandler()))))
	http.HandleFunc("/editactcomment", Rate(a, Auth(SessionChecker(s.EditCommentHandler()))))
	http.HandleFunc("/commentedited", Rate(a, Auth(SessionChecker(s.EditedCommentHandler()))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
}
