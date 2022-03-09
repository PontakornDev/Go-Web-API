package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(db *sql.DB) {
	query := `CREATE TABLE users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		create_at DATETIME,
		PRIMARY KEY (id)
		);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func Insert(db *sql.DB) {
	var username string
	var password string
	fmt.Scan(&username)
	fmt.Scan(&password)
	createAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, create_at) VALUES (?, ?, ?)`, username, password, createAt)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	fmt.Println("id:", id)
}

func delete(db *sql.DB) {
	var deleteID int
	fmt.Scan(&deleteID)
	_, err := db.Exec(`DELETE FROM users WHERE id= ?`, deleteID)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("delete id:", delete)
	}
}

func query(db *sql.DB) {

	var (
		id         int
		coursename string
		price      float64
		instructor string
	)
	for {
		var inputID int
		fmt.Scan(&inputID)
		query := "SELECT * FROM onlinecourse WHERE id = ?"
		if err := db.QueryRow(query, inputID).Scan(&id, &coursename, &price, &instructor); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, coursename, price, instructor)
	}
}

func main() {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		fmt.Println("failed to connect")
	} else {
		fmt.Println("connect successfully")
	}
	// createTable(db)
	// Insert(db)
	delete(db)
}
