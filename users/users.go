package users

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	//"net/mail"
	"strings"
)

type User struct {
	UserID          int
	Email           string
	Username        string
	Password        string
	Usertype        string
	ExternalLoginId string
}

type AuthUser struct {
	Email        string
	PasswordHash string
}

// var tpl *template.Template
var CurrentUser string

// var dbUsers = map[string]User{}
var DbSessions = make(map[string]string)

// dbUsers["yonas@hotmail.com"] = User{Username: "yonas123", Email: "yonas@hotmail.com"}

// nil

// this func registers a users username, password(as a hash) and email
func RegisterUser(db *sql.DB, username string, hash []byte, email string) {
	// db, _ = sql.Open("sqlite3", "forum.db")
	stmt, err := db.Prepare("INSERT INTO users (username, hash, email, usertype, becomemod) VALUES (?, ?, ?, 'user', 0)")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	// defer stmt.Close()

	result, _ := stmt.Exec(username, hash, email)

	// checking if the result has been added and the last inserted row
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted:", lastIns)
}

func ValidEmail(email string) bool {
	i := strings.Index(email, "@")
	fmt.Println("i:", i)

	domain := email[i+1:]
	fmt.Println("Domain: ", domain)

	_, err := net.LookupMX(domain)
	// _, err2 := mail.ParseAddress(email)
	if err != nil {
		fmt.Println("invalid email")
		return false
	}

	return true
}

func AlreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie(CurrentUser)
	if err != nil {
		return false
	}
	_, ok := DbSessions[c.Name]
	return ok
}

func SessionExists(s string) bool {
	if _, user := DbSessions[s]; user {
		return true
	}
	return false
}

func GetUserType(db *sql.DB, userId int) string {
	var userType string
	err := db.QueryRow("SELECT userType FROM users WHERE userId= ?;", userId).Scan(&userType)
	if err != nil {
		fmt.Println("error from get user", err)
	}
	return userType
}

func BecomeAMod(db *sql.DB, userId int) {
	stmt, err := db.Prepare("UPDATE users SET becomemod = 1 WHERE userId= ?")

	if err != nil {
		fmt.Println("Error updating the become mod column in users table", err)
	}

	stmt.Exec(userId)
}

func GetModRequests(db *sql.DB) []string {
	var modRequests []string
	rows, err := db.Query("SELECT username FROM users WHERE becomemod = 1 AND userType= 'user'")
	if err != nil {
		fmt.Println("Error with getting all mod requests", err)
	}

	defer rows.Close()
	for rows.Next() {
		var username string
		err2 := rows.Scan(&username)
		modRequests = append(modRequests, username)
		if err2 != nil {
			fmt.Println("cannot range through to get usernames of requested mods", err2)
		}
	}

	return modRequests
}

func AcceptMod(db *sql.DB, username string) {
	stmt, err := db.Prepare("UPDATE users set userType= 'mod' where username = ?")

	if err != nil {
		fmt.Println("User could not be updated to mod", err)
	}

	stmt.Exec(username)
}

func DeclineMod(db *sql.DB, username string) {
	stmt, err := db.Prepare("UPDATE users SET becomemod = 0 WHERE username = ?")

	if err != nil {
		fmt.Println("Error updating the become mod column in users table when declining mod", err)
	}
	stmt.Exec(username)
}

func GetAllMods(db *sql.DB) []string {
	var allMods []string

	rows, err := db.Query("SELECT username FROM users WHERE userType= 'mod'")
	if err != nil {
		fmt.Println("Error with getting all mod usernames", err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		err2 := rows.Scan(&username)
		allMods = append(allMods, username)
		if err2 != nil {
			fmt.Println("cannot range through to get usernames of moderators", err2)
		}
	}
	return allMods
}

func DemoteMod(db *sql.DB, username string) {
	stmt, err := db.Prepare("UPDATE users SET becomemod = 0, usertype = 'user' WHERE username = ?")

	if err != nil {
		fmt.Println("Error demoting a mod", err)
	}
	stmt.Exec(username)
}
