package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/categories"
	"forum/comments"
	"forum/posts"
	"forum/report"
	userimages "forum/templates/userImages"
	"forum/users"
)

type PostPageData struct {
	Posts           []posts.Post
	Comments        []comments.Comment
	LoggedIn        bool
	Liked           bool
	Disliked        bool
	CommentLiked    bool
	CommentDisliked bool
	Username        string
	Image           string
	UserID          int
	UserType        string
	Reported        bool
	Nbool           bool
	Notification    int
	Categories      []string
	EditFormID        string
	EditFormTitle     string
	EditFormContent   string
	EditComment string
	EditCommentID int

	ReportType      []string
}

var (
	UserIdint  int
	PostIDInt  int
	Imagename  string
	ReportType string
)

func (s *myServer) CreatePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// getting the user id from the url
		userId := r.URL.Query().Get("userid")
		UserIdint, _ = strconv.Atoi(userId)
		Tpl.ExecuteTemplate(w, "createpost.html", PostPageData{LoggedIn: users.AlreadyLoggedIn(r), Username: users.CurrentUser, UserID: UserIdint, Categories: categories.GetAllCategories(s.Db)})
	}
}

func (s *myServer) StorePostHandler() http.HandlerFunc {
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
			// fmt.Printf("Uploaded Image: %+v\n", handler.Filename)
			// fmt.Printf("File Size: %+v\n", handler.Size)
			// fmt.Printf("MIME Header: %+v\n", handler.Header)

			userimages.SaveImage(file, handler.Filename)
		}
		// adding the post to the database

		posts.CreatePosts(s.Db, UserIdint, title, content, Imagename)

		// formvalue for buttons. If they have been clicked, the form value returned will be "on"
		for k, v := range r.Form {
			if k == "category" {

				for _, c := range v {

					categories.AddCategory(s.Db, posts.LastIns, c)
				}

			}
		}
		// manutd := r.FormValue("manutd")
		// arsenal := r.FormValue("arsenal")
		// chelsea := r.FormValue("chelsea")
		// tottenham := r.FormValue("tottenham")
		// newcastle := r.FormValue("newcastle")
		// mancity := r.FormValue("mancity")

		// // use if statements because we need to enter the cat name instead of the returned value "on"
		// if manutd == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "manutd")
		// }
		// if arsenal == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "arsenal")
		// }
		// if chelsea == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "chelsea")
		// }
		// if newcastle == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "newcastle")
		// }
		// if tottenham == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "tottenham")
		// }
		// if mancity == "on" {
		// 	categories.AddCategory(s.Db, posts.LastIns, "mancity")
		// }
		fmt.Println("title:", title, "content:", content)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) ShowPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, v)
		}
		// had to open the database here as it wasnt picking the correct post everytime without this.
		s.Db, _ = sql.Open("sqlite3", "forum.db")
		// get the postId and display the post and its contents
		postID := r.URL.Query().Get("postid")
		PostIDInt, _ = strconv.Atoi(postID)

		NotificationInt := len(CommentNotify(s.Db)) + len(LikesNotify(s.Db)) + len(DisLikesNotify(s.Db))

		data := PostPageData{
			Posts:           posts.GetPostData(s.Db, PostIDInt),
			Comments:        comments.GetCommentData(s.Db, PostIDInt, GuserId),
			LoggedIn:        users.AlreadyLoggedIn(r),
			Liked:           UserLiked(s.Db),
			Disliked:        UserDisliked(s.Db),
			CommentLiked:    CommentUserLiked(s.Db),
			CommentDisliked: CommentUserDisliked(s.Db),
			Username:        users.CurrentUser,
			Image:           Imagename,
			UserID:          GuserId,
			UserType:        users.GetUserType(s.Db, GuserId),
			Notification:    NotificationInt,
			Reported:        ModReported(s.Db),
			ReportType:      report.GetReportType(s.Db, PostIDInt),
		}

		if NotificationInt > 0 {
			data.Nbool = true
		}

		Tpl.ExecuteTemplate(w, "showpost.html", data)
	}
}

func (s *myServer) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("ACTIVTY URL PATH", r.URL.Path)
		posts.DeletePost(s.Db, PostIDInt)

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) ReportHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SPostID := strconv.Itoa(PostIDInt)
		r.ParseForm()

		if ModReported(s.Db) {
			fmt.Fprint(w, "You have already reported this post")
			// http.Redirect(w, r, "/showpost/?postid="+SPostID, http.StatusSeeOther)
		} else {
			postPageFilter := r.FormValue("reportfilter")

			if postPageFilter == "Irrelevant" {
				report.ReportButton(s.Db, users.CurrentUser, "irrelevant", PostIDInt)
			}
			if postPageFilter == "Obscene" {
				report.ReportButton(s.Db, users.CurrentUser, "obscene", PostIDInt)
			}
			if postPageFilter == "Illegal" {
				report.ReportButton(s.Db, users.CurrentUser, "illegal", PostIDInt)
			}
			if postPageFilter == "Insulting" {
				report.ReportButton(s.Db, users.CurrentUser, "insulting", PostIDInt)
			}

			http.Redirect(w, r, "/showpost/?postid="+SPostID, http.StatusSeeOther)
		}
	}
}

func ModReported(db *sql.DB) bool {
	// check if post already been reportd by mod
	userStmt := "SELECT username FROM report WHERE username = ? AND postID = ?"
	row := db.QueryRow(userStmt, users.CurrentUser, PostIDInt)

	var uID string
	var pID string
	err := row.Scan(&uID, &pID)
	if err != sql.ErrNoRows {
		fmt.Println("Mod has already reported this post", err)
		return true
	}
	return false
}
