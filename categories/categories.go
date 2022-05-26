package categories

import (
	"database/sql"
	"fmt"
)

type Category struct {
	CategoryName string
	CategoryID   int
	PostID       int
}

var db *sql.DB

func AddCategory(db *sql.DB, postID int64, categoryname string) {
	stmt, err := db.Prepare("INSERT INTO postcategory (postID, categoryname) VALUES (?, ?)")

	if err != nil {
		fmt.Println("error preparing statment", err)
		return
	}

	result, _ := stmt.Exec(postID, categoryname)

	// checking if the result has been added and the last inserted row
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rows affected:", rowsAff)
	fmt.Println("last inserted:", lastIns)

}

func AdminAddCategory(db *sql.DB, categoryname string) {
	stmt, err := db.Prepare("INSERT INTO categories(categoryname) VALUES(?)")
	if err != nil {
		fmt.Println("error preparing statement", err)
		return
	}

	if CategoryExistsCheck(db, categoryname) {
		fmt.Println("CANNOT ADD CATEGORY AS IT EXISTS")
		return
	}
	result, _ := stmt.Exec(categoryname)

	// checking if result has been added and the last inserted row
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rows affected: ", rowsAff)
	fmt.Println("last inserted: ", lastIns)
}

func CategoryExistsCheck(db *sql.DB, categoryname string) bool {
	stmt := "SELECT categoryname FROM categories WHERE categoryname = ?"
	row := db.QueryRow(stmt, categoryname)
	err := row.Scan(&categoryname)
	if err != nil {
		return false
	}
	return true
}

func GetAllCategories(db *sql.DB) []string {
	rows, err := db.Query("SELECT categoryname FROM categories")
	if err != nil {
		fmt.Println(err)
	}
	var cats []string
	defer rows.Close()
	for rows.Next() {
		var c string
		err2 := rows.Scan(&c)
		if err != nil {
			fmt.Println(err2)
		}
		cats = append(cats, c)
	}

	return cats

}

func DeleteCategory(db *sql.DB, categoryname string) {
	stmt, err := db.Prepare("DELETE FROM categories WHERE categoryname=?")
	if err != nil {
		fmt.Println("error deleting category", err)
	}
	stmt.Exec(categoryname)

	stmt3, err3 := db.Prepare("DELETE FROM postcategory WHERE categoryname = ?")
	if err3 != nil {
		fmt.Println("error deleting reports from the report table", err3)
	}
	stmt3.Exec(categoryname)
}
