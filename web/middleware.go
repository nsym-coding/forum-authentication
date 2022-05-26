package web

import (
	"fmt"
	"forum/users"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie(users.CurrentUser)
		if err != nil {
			var errorMessage string
			// http.Redirect(w, r, "/login", http.StatusSeeOther)
			if r.URL.Path == "/likes" || r.URL.Path == "/dislikes" {
				errorMessage = "To like/dislike a post, please log in/register"
			} else if r.URL.Path == "/storecomment" {
				errorMessage = "To comment on this post please log in/register"
			} else if r.URL.Path == "/commentlikes" || r.URL.Path == "/commentdislikes" {
				errorMessage = "To like/dislike this comment, please log in/register"
			}
			Tpl.ExecuteTemplate(w, "login.html", errorMessage)
			return
		}
		fmt.Println(c.Value)
		HandlerFunc.ServeHTTP(w, r)
	}
}

func SessionChecker(HandlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if users.SessionExists(users.CurrentUser) {
			fmt.Println()
			fmt.Println("Session Exists/ID", users.DbSessions)
			fmt.Println()

			// Read the cookie,if there are any
			cookie, err := r.Cookie(users.CurrentUser)

			//i.e session exists (on map) but no active cookie on (presumably) new client
			if err != nil {

				// create new cookie for user on this client
				id := uuid.Must(uuid.NewV4())
				c := &http.Cookie{
					Name:  users.CurrentUser,
					Value: id.String(),
				}

				http.SetCookie(w, c)
				users.DbSessions[users.CurrentUser] = c.Value
				HandlerFunc.ServeHTTP(w, r)
				fmt.Println("Map Values reassigned for new client log in:  ", users.DbSessions)
				fmt.Println()

				return

				//session exists but differing UUID,  logout/close session
			} else if cookie.Value != users.DbSessions[users.CurrentUser] {

				//expire cookie as there's an active session elsewhere
				for _, cookie := range r.Cookies() {
					fmt.Println("-------Test Delete 3")

					if cookie.Name == users.CurrentUser {
						fmt.Println("checking if loop entered--------------------------------------------------------")
						fmt.Println("1st max age check ----> ", cookie.MaxAge, cookie.Value)

						cookie.MaxAge = -1
						//cookie.Value = ""
						//cookie.Expires = time.Date(2022, time.April, 21, 0000, 00, 00, 00, time.UTC)
						fmt.Println("2nd max age check ----> ", cookie.MaxAge, cookie.Value)
					}
					http.SetCookie(w, cookie)
				}

				//redirects to log in page when prior session exists
				http.Redirect(w, r, "/logout", http.StatusSeeOther)

				fmt.Println("Cookie deleted due to pre-existing session")
				fmt.Println()
				return

				//UUID matches that within map, active session, no conflicts
			} else if cookie.Value == users.DbSessions[users.CurrentUser] {
				HandlerFunc.ServeHTTP(w, r)
				fmt.Println("Active session on client, no changes made")
				fmt.Println()
				return

			}

		}
	}
}

func HomepageSessionChecker(HandlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println(users.DbSessions)
		if users.SessionExists(users.CurrentUser) {
			fmt.Println()
			//fmt.Println("Session Exists/ID", users.DbSessions)
			fmt.Println()

			// Read the cookie,if there are any
			cookie, err := r.Cookie(users.CurrentUser)

			//i.e session exists (on map) but no active cookie on (presumably) new client
			if err != nil {
				// create new cookie for user on this client
				HandlerFunc.ServeHTTP(w, r)
				return

				//session exists but differing UUID,  logout/close session
			} else if cookie.Value != users.DbSessions[users.CurrentUser] {

				//expire cookie as there's an active session elsewhere
				for _, cookie := range r.Cookies() {

					if cookie.Name == users.CurrentUser {

						cookie.MaxAge = -1
						//cookie.Value = ""
						//cookie.Expires = time.Date(2022, time.April, 21, 0000, 00, 00, 00, time.UTC)
						fmt.Println("2nd max age check ----> ", cookie.MaxAge, cookie.Value)
					}
					http.SetCookie(w, cookie)
				}

				//redirects to log in page when prior session exists
				http.Redirect(w, r, "/logout", http.StatusSeeOther)

				fmt.Println("Cookie deleted due to pre-existing session")
				fmt.Println()
				return

				//UUID matches that within map, active session, no conflicts
			} else if cookie.Value == users.DbSessions[users.CurrentUser] {
				HandlerFunc.ServeHTTP(w, r)
				fmt.Println("Active session on client, no changes made")
				fmt.Println()
				return

			}

		}
		HandlerFunc.ServeHTTP(w, r)
	}
}
