package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)


func main() {
	sqliteDatabase, _ := sql.Open("sqlite3", "./db.db")
	go Publisher()
	Subscriber(sqliteDatabase)
	select {}
}


func Select(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var password string
		var age int
		var hours_spent int
		err := rows.Scan(&name, &password, &age, &hours_spent)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("name: "+name, "password: "+password, "age: "+strconv.Itoa(age), "hours_spent: "+strconv.Itoa(hours_spent))
	}
}

func Insert(db *sql.DB, data Data) {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO USER (name, password, age, hours_spent) VALUES ('%v', '%v', %v, %v)",data.name, data.password, data.age, data.hours_spent))
	if err != nil {
		log.Fatal(err)
	}
}

