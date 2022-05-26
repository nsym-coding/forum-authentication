package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"forum/categories"
	"forum/comments"
	"forum/posts"
	userimages "forum/templates/userImages"
	"forum/users"

	"forum/report"
)

type ActivityPage struct {
	Posts             []posts.HomepagePosts
	CommentsWithPosts []posts.ActPage
	LikedPosts        []posts.ActPage
	DislikedPosts     []posts.ActPage
	LikedComments     []posts.ActPage
	DislikedComments  []posts.ActPage
	Comments          []posts.Post
	Username          string
	LoggedIn          bool
	UserID            int
	Nbool             bool
	Notification      int
	CommentNote       []Notify
	LikeNote          []Notify
	DisLikeNote       []Notify
	ModRequests       []string
	UserType          string
	AllMods           []string
	Categories        []string
	EditFormID        string
	EditFormTitle     string
	EditFormContent   string
	EditCommentID     int
	ReportedPosts     []report.Report
}

var (
	editformID      string
	editformTitle   string
	editformContent string
	editbool        bool
)

func (s *myServer) EditedCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limits requests to 20MB (x is the limiter where x<<20)
		r.Body = http.MaxBytesReader(w, r.Body, 20<<20)

		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, "Images must be less than 20MB!!", http.StatusBadRequest)
			return
		}

		comment := r.FormValue("comment")
		commentid, _ := strconv.Atoi(r.FormValue(("editpage")))

		UpdateComment(s.Db, comment, commentid)

		Suserid := strconv.Itoa(GuserId)

		http.Redirect(w, r, "/activitypage?userid="+Suserid, http.StatusSeeOther)
	}
}

func (s *myServer) EditCommentHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			var commID string

			for _, value := range r.Form["editcomment"] {
				commID += value
			}

			coMMInt, _ := strconv.Atoi(commID)

			comment := comments.GetCommentEdit(s.Db, coMMInt)

			notify := len(CommentNotify(s.Db)) + len(LikesNotify(s.Db)) + len(DisLikesNotify(s.Db))

			if notify > 0 {
				editbool = true
			}

			var EditComment PostPageData

			EditComment.LoggedIn = users.AlreadyLoggedIn(r)
			EditComment.Username = users.CurrentUser
			EditComment.UserID = GuserId
			EditComment.Notification = notify
			EditComment.Nbool = editbool
			EditComment.EditComment = comment
			EditComment.EditCommentID = coMMInt

			Tpl.ExecuteTemplate(w, "editactcomment.html", EditComment)

		}
	}
}

func (s *myServer) EditPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// limits requests to 20MB (x is the limiter where x<<20)
		r.Body = http.MaxBytesReader(w, r.Body, 20<<20)

		err := r.ParseMultipartForm(20 << 20)
		if err != nil {
			http.Error(w, "Images must be less than 20MB!!", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		postid, _ := strconv.Atoi(r.FormValue(("editpage")))

		x, _, _ := r.FormFile("userimage")
		if x != nil {
			// Get handler for filename, size and headers
			// file, handler, err := r.FormFile("userimage2") //Change it to this to test internal error.
			file, handler, err := r.FormFile("userimage")
			if err != nil {
				//	Tpl.ExecuteTemplate(w, "error.html", nil)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, " An internal server error has occurred: ", http.StatusInternalServerError)
				return
				// fmt.Println("Error Retrieving the File")
				// fmt.Println(err)
				// return
			}

			defer file.Close()

			Imagename = handler.Filename

			userimages.SaveImage(file, handler.Filename)
		}

		// update database with new values
		UpdatePost(s.Db, title, content, Imagename, postid)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) EditActComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			var formString string

			formSlice := r.Form["editpost"]

			for _, char := range formSlice {
				formString += char
			}

			formSplit := strings.Split(formString, "&")

			editformID = formSplit[0]
			editformTitle = formSplit[1]
			editformContent = formSplit[2]

		}

		notify := len(CommentNotify(s.Db)) + len(LikesNotify(s.Db)) + len(DisLikesNotify(s.Db))

		if notify > 0 {
			editbool = true
		}

		var EditPage PostPageData

		EditPage.LoggedIn = users.AlreadyLoggedIn(r)
		EditPage.Username = users.CurrentUser
		EditPage.UserID = GuserId
		EditPage.EditFormID = editformID
		EditPage.EditFormTitle = editformTitle
		EditPage.EditFormContent = editformContent
		EditPage.Notification = notify
		EditPage.Nbool = editbool

		Tpl.ExecuteTemplate(w, "editactpost.html", EditPage)
	}
}

func (s *myServer) DeleteActComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, id := range v {
				idInt, _ := strconv.Atoi(id)
				comments.DeleteComment(s.Db, idInt)
			}
		}
		stringGID := strconv.Itoa(GuserId)
		http.Redirect(w, r, "/activitypage?userid="+stringGID, http.StatusSeeOther)
	}
}

func (s *myServer) DeleteActPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		for _, v := range r.Form {
			for _, id := range v {
				idInt, _ := strconv.Atoi(id)
				posts.DeletePost(s.Db, idInt)
			}
		}

		stringGID := strconv.Itoa(GuserId)

		http.Redirect(w, r, "/activitypage?userid="+stringGID, http.StatusSeeOther)
	}
}

func (s *myServer) ActivityPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data ActivityPage

		data.Notification = (len(CommentNotify(s.Db)) + len(LikesNotify(s.Db)) + len(DisLikesNotify(s.Db)))

		if data.Notification > 0 {
			data.Nbool = true
		}

		data.CommentNote = CommentNotify(s.Db)

		data.LikeNote = LikesNotify(s.Db)

		data.DisLikeNote = DisLikesNotify(s.Db)

		data.Posts = posts.UsersPostsHomepageData(s.Db, GuserId)

		data.CommentsWithPosts = posts.ActivityComments(s.Db, GuserId)

		data.LikedPosts = posts.ActivityPostLikes(s.Db, GuserId)

		data.DislikedPosts = posts.ActivityPostDislikes(s.Db, GuserId)

		data.LikedComments = posts.ActivityCommentLikes(s.Db, GuserId)

		data.ReportedPosts = report.GetReportData(s.Db)

		data.DislikedComments = posts.ActivityCommentDislikes(s.Db, GuserId)
		data.Username = users.CurrentUser
		data.LoggedIn = users.AlreadyLoggedIn(r)
		data.ModRequests = users.GetModRequests(s.Db)
		data.UserType = users.GetUserType(s.Db, GuserId)
		data.UserID = GuserId
		data.AllMods = users.GetAllMods(s.Db)
		data.Categories = categories.GetAllCategories(s.Db)
		SuserID := strconv.Itoa(GuserId)

		if string(r.URL.RawQuery[len(r.URL.RawQuery)-1]) != SuserID {
			http.Error(w, "Incorrect user request made!", http.StatusBadRequest)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println("------------------------")
		fmt.Println(data.Posts)

		Tpl.ExecuteTemplate(w, "activitypage.html", data)

		func() {
			ResetCommentNotified(s.Db)
		}()

		func() {
			ResetLikesNotified(s.Db)
		}()

		func() {
			ResetDisLikesNotified(s.Db)
		}()
	}
}

func (s *myServer) AdminAddingCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println("key: ", k, "value: ", v)
		}
		SuserID := strconv.Itoa(GuserId)
		categoryname := r.FormValue("category")
		if categoryname != "" {
			categories.AdminAddCategory(s.Db, categoryname)
		}
		// http.Redirect(w, r, "/home", http.StatusSeeOther)

		http.Redirect(w, r, "/activitypage?userid="+SuserID, http.StatusSeeOther)
	}
}

func (s *myServer) AdminDeletingCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, c := range v {
				categories.DeleteCategory(s.Db, c)
			}
		}
		SuserID := strconv.Itoa(GuserId)

		http.Redirect(w, r, "/activitypage?userid="+SuserID, http.StatusSeeOther)
	}
}
