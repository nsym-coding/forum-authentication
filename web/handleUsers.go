package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/users"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type GitOAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

type GoogleOAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

type GoogleEmail struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

type GitUserName struct {
	Username string `json:"login"`
	//Email    string                 `json:"email"`
	//X map[string]interface{} `json:"-"`
}

type GitEmail struct {
	Email string `json:"email"`
	//X     map[string]interface{} `json:"-"`
}

var GuserId int
var GitLoginName string
var GoogleUserName string

const GitclientID = "3c0054a12dde712c09fb"
const GitclientSecret = "d44989291dc0c3eb40bf37dc84f948310aa10672"

const GoogleClientID = "994843129537-tmvgh06pi2fpuaite0b5em56so3nng9h.apps.googleusercontent.com"
const GoogleClientSecret = "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"

func (s *myServer) GoogleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Tpl.ExecuteTemplate(w, "googlelogin.html", nil)
	}
}

func (s *myServer) GoogleOAUTHLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprint(w, "checking it redirects")
		s.Db, _ = sql.Open("sqlite3", "forum.db")
		httpClient := http.Client{}
		err := r.ParseForm()

		if err != nil {
			fmt.Fprint(w, "Couldn't parse query", err)
		}

		code := r.FormValue("code")

		reqURL := fmt.Sprintf("https://oauth2.googleapis.com/token?&code=%s&client_id=%s&client_secret=%s&redirect_uri=https://localhost:8080/google/redirect&grant_type=authorization_code", code, GoogleClientID, GoogleClientSecret)

		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprint(w, "Couldn't create HTTP request", err)
		}
		//req.Header.Add("Authorization", "Basic"+GoogleClientID+":"+GoogleClientSecret)

		//Setting this header so the response is in json format
		req.Header.Set("accept", "application/json")

		//Sending out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprint(w, "couldn't send http request", err)
		}
		defer res.Body.Close()

		var g GoogleOAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&g); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		var bearer = "Bearer " + g.AccessToken

		reqURL2 := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", g.AccessToken)

		// Create a new request using http
		req1, err := http.NewRequest("GET", reqURL2, nil)

		req1.Header.Add("Authorization", bearer)

		client := &http.Client{}
		resp, err := client.Do(req1)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		body1, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}

		var GMailAdd GoogleEmail

		err3 := json.Unmarshal(body1, &GMailAdd)
		if err3 != nil {
			log.Fatal(err3)
		}

		for i := 0; i < len(GMailAdd.Email); i++ {
			if GMailAdd.Email[i] == '@' {
				GoogleUserName = GMailAdd.Email[0:i] + GMailAdd.ID[0:4]
			}
		}

		var (
			UnameExists  bool
			UemailExists bool
		)

		// check if email already exists
		emailStmt := "SELECT userID FROM users WHERE email = ?"
		rowE := s.Db.QueryRow(emailStmt, GMailAdd.Email)
		var uID string
		err5 := rowE.Scan(&uID)
		if err5 != sql.ErrNoRows {
			fmt.Println("email already exists, err:", err5)
			UemailExists = true

		}

		// check if username already exists
		userStmt := "SELECT userID FROM users WHERE username = ?"
		rowU := s.Db.QueryRow(userStmt, GoogleUserName)
		var uIDs string
		error := rowU.Scan(&uIDs)
		if error != sql.ErrNoRows {
			fmt.Println("username already exists, err:", error)
			UnameExists = true

		}

		var blank []byte

		if !UnameExists && !UemailExists {

			users.RegisterUser(s.Db, GoogleUserName, blank, GMailAdd.Email)

		} else if !UnameExists && UemailExists {
			Tpl.ExecuteTemplate(w, "register.html", "email already exists")
			return
			//fmt.Print("This is an existing user")
		} else if UnameExists && !UemailExists {
			Tpl.ExecuteTemplate(w, "register.html", "username already exists")
			return
		}

		http.Redirect(w, r, "/loginauth", http.StatusSeeOther)

	}
}

func (s *myServer) GitLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Tpl.ExecuteTemplate(w, "githublogin.html", nil)
	}
}

func (s *myServer) GitOAUTHLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Db, _ = sql.Open("sqlite3", "forum.db")
		httpClient := http.Client{}
		err := r.ParseForm()

		if err != nil {
			fmt.Fprint(w, "Couldn't parse query", err)
		}

		code := r.FormValue("code")

		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", GitclientID, GitclientSecret, code)

		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprint(w, "Couldn't create HTTP request", err)
		}

		//Setting this header so the response is in json format
		req.Header.Set("accept", "application/json")

		//Sending out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprint(w, "couldn't send http request", err)
		}
		defer res.Body.Close()

		// Parse the request body into the `Git` struct
		var t GitOAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		var bearer = "Bearer " + t.AccessToken

		// Create a new request using http
		req1, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		req2, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
		// add authorization header to the req
		req1.Header.Add("Authorization", bearer)
		req2.Header.Add("Authorization", bearer)
		//	req1.SetBasicAuth(t.AccessToken, "x-oauth-basic")

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req1)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}

		resp1, err := client.Do(req2)

		body1, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}

		body2, err := ioutil.ReadAll(resp1.Body)

		//fmt.Fprintln(w, string([]byte(body1)))
		var GitUserMail []string

		defer resp.Body.Close()
		var Guser GitUserName
		var Gmail []GitEmail

		err3 := json.Unmarshal(body1, &Guser)
		if err3 != nil {
			log.Fatal(err3)
		}

		err4 := json.Unmarshal(body2, &Gmail)
		if err4 != nil {
			log.Fatal(err4)
		}
		for _, v := range Gmail {
			GitUserMail = append(GitUserMail, v.Email)
		}

		GitLoginName = Guser.Username

		var (
			UnameExists  bool
			UemailExists bool
		)

		// check if email already exists
		emailStmt := "SELECT userID FROM users WHERE email = ?"
		rowE := s.Db.QueryRow(emailStmt, GitUserMail[0])
		var uID string
		err5 := rowE.Scan(&uID)
		if err5 != sql.ErrNoRows {
			fmt.Println("email already exists, err:", err5)
			UemailExists = true

		}

		// check if username already exists
		userStmt := "SELECT userID FROM users WHERE username = ?"
		rowU := s.Db.QueryRow(userStmt, GitLoginName)
		var uIDs string
		error := rowU.Scan(&uIDs)
		if error != sql.ErrNoRows {
			fmt.Println("username already exists, err:", error)
			UnameExists = true

		}

		var blank []byte

		if !UnameExists && !UemailExists {

			users.RegisterUser(s.Db, GitLoginName, blank, GitUserMail[0])
		} else if !UnameExists && UemailExists {
			Tpl.ExecuteTemplate(w, "register.html", "email already exists")
			return
			//fmt.Print("This is an existing user")
		} else if UnameExists && !UemailExists {
			Tpl.ExecuteTemplate(w, "register.html", "username already exists")
			return
		}

		http.Redirect(w, r, "/loginauth", http.StatusSeeOther)

	}
}

func (s *myServer) RegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Tpl.ExecuteTemplate(w, "register.html", nil)
	}
}

func (s *myServer) RegisterAuthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Db, _ = sql.Open("sqlite3", "forum.db")
		fmt.Println("********registerAuthHandler running*******")

		r.ParseForm()

		username := r.FormValue("username")
		fmt.Println("username", username)

		nameAlphaNumeric := true

		// check whether each character in username AlphaNumeric

		for _, char := range username {
			if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
				nameAlphaNumeric = false
			}
		}

		// check username length is 6 < x < 10

		var usernameLength bool

		if len(username) <= 10 && len(username) >= 6 {
			usernameLength = true
		}

		// check password is valid given all conditions

		password := r.FormValue("password")
		fmt.Println("password", password, "\npswdLength", len(password))

		var pswdLowercase, pswdUpperCase, pswdNumber, pswdSpecical, pswdNoSpaces, passwordLength bool

		pswdNoSpaces = true

		for _, char := range password {
			switch {
			case unicode.IsLower(char):
				pswdLowercase = true
			case unicode.IsUpper(char):
				pswdUpperCase = true
			case unicode.IsNumber(char):
				pswdNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				pswdSpecical = true
			case unicode.IsSpace(int32(char)):
				pswdNoSpaces = false
			}
		}
		// check password length is 5 < x < 20
		if len(password) <= 20 && len(password) >= 5 {
			passwordLength = true
		}

		if !nameAlphaNumeric || !usernameLength || !pswdLowercase || !pswdUpperCase || !pswdNumber || !pswdSpecical || !pswdNoSpaces || !passwordLength {
			Tpl.ExecuteTemplate(w, "register.html", "Please check your username or password")
			return
		}

		email := r.FormValue("email")
		email = strings.ToLower(email)
		fmt.Println("email", email)

		if !users.ValidEmail(email) {
			Tpl.ExecuteTemplate(w, "register.html", "Please enter a valid email address")
			return
		}

		// check if email already exists
		emailStmt := "SELECT userID FROM users WHERE email = ?"
		rowE := s.Db.QueryRow(emailStmt, email)
		var uID string
		err := rowE.Scan(&uID)
		if err != sql.ErrNoRows {
			fmt.Println("email already exists, err:", err)
			Tpl.ExecuteTemplate(w, "register.html", "email already exists")
			return
		}

		// check if username already exists
		userStmt := "SELECT userID FROM users WHERE username = ?"
		rowU := s.Db.QueryRow(userStmt, username)
		var uIDs string
		error := rowU.Scan(&uIDs)
		if error != sql.ErrNoRows {
			fmt.Println("username already exists, err:", error)
			Tpl.ExecuteTemplate(w, "register.html", "username already exists")
			return
		}

		//  create hash from password
		var hash []byte

		hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("bcrypt err:", err)
			Tpl.ExecuteTemplate(w, "register.html", "there was a problem registering user")
			return
		}

		fmt.Println("hash:", hash)
		fmt.Println("string(hash)", string(hash))
		users.RegisterUser(s.Db, username, hash, email)
		// fmt.Fprintf(w, "congrats your account has been successfully created")
		Tpl.ExecuteTemplate(w, "login.html", "User succesfully registered. Login to forum")
	}
}

func (s *myServer) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("login handler running")
		fmt.Println("checking bool-----> ", users.AlreadyLoggedIn(r))

		Tpl.ExecuteTemplate(w, "login.html", nil)
	}
}

func (s *myServer) LoginAuthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if GitLoginName == "" && GoogleUserName == "" {
			fmt.Println("CHECKING IF GIT LOGIN NAME IS EMPTY", GitLoginName)
			// we need to figure out whether we have to close the database at some point to save resources.
			s.Db, _ = sql.Open("sqlite3", "forum.db")
			fmt.Println("login authHandler running")
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")

			users.CurrentUser = username
			fmt.Println("username:", username, "password", password)
			// get password(hash form) from db to compare with users supplied password
			var hash string
			stmt := "SELECT hash FROM users WHERE username = ?"
			row := s.Db.QueryRow(stmt, username)
			err := row.Scan(&hash)
			fmt.Println("hash from db:", hash)
			if err != nil {
				fmt.Println("error with username, may not exist")
				// keep the message to the user a little more vague, so hackers dont know whether you entered an incorrect username or password
				Tpl.ExecuteTemplate(w, "login.html", "check username and password")
				return
			}

			// get userId to pass onto createpost handler

			var userID int

			stmt2 := "SELECT userID FROM users WHERE username = ?"
			row2 := s.Db.QueryRow(stmt2, username)
			err2 := row2.Scan(&userID)
			fmt.Println("userID from db:", userID)
			GuserId = userID
			if err2 != nil {
				fmt.Println("user not found in db")
			}

			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
			if err == nil {

				//NO ACTIVE SESSION/FIRST TIME
				id := uuid.Must(uuid.NewV4())
				c := &http.Cookie{
					Name:  username,
					Value: id.String(),
				}

				http.SetCookie(w, c)
				users.DbSessions[username] = c.Value
				// tpl.ExecuteTemplate(w, "loginauth.html", userID)
				http.Redirect(w, r, "/home", http.StatusSeeOther)

				/////////remove///////////////////
				fmt.Println("sessionbool", users.SessionExists(username))

				for _, cookie := range r.Cookies() {
					fmt.Println()
					fmt.Println("Name : ", cookie.Name)
					fmt.Println("Value/UUID : ", cookie.Value)
				}

				fmt.Println("First time log-in successful")
				fmt.Println()
				/////////////////////////////////////
				return

			}

			fmt.Println("incorrect password")
			Tpl.ExecuteTemplate(w, "login.html", "check username and password")
		}
		//Git log in
		if GitLoginName != "" {
			s.Db, _ = sql.Open("sqlite3", "forum.db")
			fmt.Println("login authHandler running")
			r.ParseForm()
			username := GitLoginName
			GitLoginName = ""
			// password := r.FormValue("password")

			users.CurrentUser = username
			fmt.Println("username:", username)
			// get password(hash form) from db to compare with users supplied password
			// var hash string
			// stmt := "SELECT hash FROM users WHERE username = ?"
			// row := s.Db.QueryRow(stmt, username)
			// err := row.Scan(&hash)
			// fmt.Println("hash from db:", hash)
			// if err != nil {
			// 	fmt.Println("error with username, may not exist")
			// 	// keep the message to the user a little more vague, so hackers dont know whether you entered an incorrect username or password
			// 	Tpl.ExecuteTemplate(w, "login.html", "check username and password")
			// 	return
			// }

			// get userId to pass onto createpost handler

			var userID int

			stmt2 := "SELECT userID FROM users WHERE username = ?"
			row2 := s.Db.QueryRow(stmt2, username)
			err2 := row2.Scan(&userID)
			fmt.Println("userID from db:", userID)
			GuserId = userID
			if err2 != nil {
				fmt.Println("user not found in db")
			}

			//NO ACTIVE SESSION/FIRST TIME
			id := uuid.Must(uuid.NewV4())
			c := &http.Cookie{
				Name:  username,
				Value: id.String(),
			}

			http.SetCookie(w, c)
			users.DbSessions[username] = c.Value
			// tpl.ExecuteTemplate(w, "loginauth.html", userID)
			http.Redirect(w, r, "/home", http.StatusSeeOther)

			/////////remove///////////////////
			fmt.Println("sessionbool", users.SessionExists(username))

			for _, cookie := range r.Cookies() {
				fmt.Println()
				fmt.Println("Name : ", cookie.Name)
				fmt.Println("Value/UUID : ", cookie.Value)
			}

			fmt.Println("First time log-in successful")
			fmt.Println()
			/////////////////////////////////////
			return
		}
		if GoogleUserName != "" {
			s.Db, _ = sql.Open("sqlite3", "forum.db")
			fmt.Println("login authHandler running")
			r.ParseForm()
			username := GoogleUserName
			GoogleUserName = ""
			// password := r.FormValue("password")

			users.CurrentUser = username
			fmt.Println("username:", username)
			// get password(hash form) from db to compare with users supplied password
			// var hash string
			// stmt := "SELECT hash FROM users WHERE username = ?"
			// row := s.Db.QueryRow(stmt, username)
			// err := row.Scan(&hash)
			// fmt.Println("hash from db:", hash)
			// if err != nil {
			// 	fmt.Println("error with username, may not exist")
			// 	// keep the message to the user a little more vague, so hackers dont know whether you entered an incorrect username or password
			// 	Tpl.ExecuteTemplate(w, "login.html", "check username and password")
			// 	return
			// }

			// get userId to pass onto createpost handler

			var userID int

			stmt2 := "SELECT userID FROM users WHERE username = ?"
			row2 := s.Db.QueryRow(stmt2, username)
			err2 := row2.Scan(&userID)
			fmt.Println("userID from db:", userID)
			GuserId = userID
			if err2 != nil {
				fmt.Println("user not found in db")
			}

			//NO ACTIVE SESSION/FIRST TIME
			id := uuid.Must(uuid.NewV4())
			c := &http.Cookie{
				Name:  username,
				Value: id.String(),
			}

			http.SetCookie(w, c)
			users.DbSessions[username] = c.Value
			// tpl.ExecuteTemplate(w, "loginauth.html", userID)
			http.Redirect(w, r, "/home", http.StatusSeeOther)

			/////////remove///////////////////
			fmt.Println("sessionbool", users.SessionExists(username))

			for _, cookie := range r.Cookies() {
				fmt.Println()
				fmt.Println("Name : ", cookie.Name)
				fmt.Println("Value/UUID : ", cookie.Value)
			}

			fmt.Println("First time log-in successful")
			fmt.Println()
			/////////////////////////////////////
			return

		}
		// fmt.Println("incorrect password")
		// Tpl.ExecuteTemplate(w, "login.html", "check username and password")

	}
}

func (s *myServer) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(users.CurrentUser)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// delete the session
		if c.Value == users.DbSessions[c.Name] {

			delete(users.DbSessions, c.Name)
		}
		// remove the cookie
		c = &http.Cookie{
			Name:   users.CurrentUser,
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, c)

		fmt.Println("User logged out and redirected to the log-in page")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

func (s *myServer) BecomeAModHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users.BecomeAMod(s.Db, GuserId)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) AcceptAModHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, name := range v {
				users.AcceptMod(s.Db, name)
			}
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) DeclineAModHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, name := range v {
				users.DeclineMod(s.Db, name)
			}
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

func (s *myServer) DemoteAModHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		for _, v := range r.Form {
			for _, name := range v {
				users.DemoteMod(s.Db, name)
			}
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}
