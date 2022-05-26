package report

import (
	"database/sql"
	"fmt"
)

type Report struct {
	ReportType   string
	ModUsername  string
	PostTitle    string
	ReportPostID int
}

var LastIns int64

func ReportButton(db *sql.DB, username string, reporttype string, postID int) {
	stmt, err := db.Prepare("INSERT INTO report(username, reporttype, postID) VALUES(?, ?, ?)")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}

	result, _ := stmt.Exec(username, reporttype, postID)
	rowsAff, _ := result.RowsAffected()
	LastIns, _ = result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted", LastIns)
}

func GetReportType(db *sql.DB, postID int) []string {
	rows, err := db.Query("SELECT reporttype FROM report WHERE postID= ?;", postID)
	if err != nil {
		fmt.Println("error from get report", err)
	}

	var reportType []string
	defer rows.Close()
	for rows.Next() {
		var r string
		err2 := rows.Scan(&r)
		if err != nil {
			fmt.Println(err2)
		}
		reportType = append(reportType, r)
	}
	return reportType
}

// func GetMod(db *sql.DB, reportID int) string {
// 	var username string
// 	err := db.QueryRow("SELECT username FROM users WHERE reportID= ?;", reportID).Scan(&username)
// 	if err != nil {
// 		fmt.Println("error from get mod", err)
// 	}
// 	return username
// }

func GetReportData(db *sql.DB) []Report {
	rows, err := db.Query(`SELECT reporttype, report.username, postTitle, report.postID
	FROM report 
	INNER JOIN post ON post.postID = report.postID 
	INNER JOIN users ON users.username = report.username 
	;`)
	if err != nil {
		fmt.Println(err)
	}
	report := []Report{}
	defer rows.Close()
	for rows.Next() {
		var r Report
		err2 := rows.Scan(&r.ReportType, &r.ModUsername, &r.PostTitle, &r.ReportPostID)
		report = append(report, r)
		if err2 != nil {
			fmt.Println(err2)
		}
	}
	return report
}
