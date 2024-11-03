package mysql

import (
	"context"
	"database/sql" // interact with a sql database
	"log"          // self explanatory
	"time"
)

type Comment struct {
	Title string `json:"Title"`
	Body  string `json:"Body"`
}

func connectAndReturn() *sql.DB {
	const (
		DB_HOST = "tcp(mysql:3307)"
		DB_NAME = "demo-db"
		DB_USER = "root"
		DB_PASS = "secret"
	)

	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"

	for {
		log.Print(".")
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Error: " + err.Error())
			return nil
		}
		err = db.Ping()
		if err != nil {
			log.Println("Error: " + err.Error())
			time.Sleep(2 * time.Second)
			continue
		} else {
			log.Println("Success!")
			initializeTable(db)
			return db
		}
	}
}

func initializeTable(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS `comments`")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE `comments`")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS `comments` ( id MEDIUMINT NOT NULL AUTO_INCREMENT, title varchar(32), body varchar(32), PRIMARY KEY (id) )")
	if err != nil {
		panic(err)
	}

	query := "INSERT INTO `comments` (`title`, `body`) VALUES (?, ?)"
	result, err := db.ExecContext(context.Background(), query, "Very cool kanye", "but can it do sql?")
	if err != nil {
		log.Fatalf("impossible insert into %s: %s", "comments", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve id: %d from %s: %s", id, "comments", err)
	}
}

func GetAComment() Comment {
	db := connectAndReturn()
	rows, err := db.Query("SELECT * FROM `comments`")
	if err != nil {
		log.Fatal("Error during query: " + err.Error())
	}
	var comment Comment

	for rows.Next() {
		rows.Scan(&comment.Title, &comment.Body)
		log.Println("the comment: " + comment.Title + " " + comment.Body)
		return comment
	}
	return Comment{Title: "Title", Body: "Body"}

}
