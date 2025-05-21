package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=consultancy sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Failed to open connection: %v\n", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		return
	}

	fmt.Println("Successfully connected to database!")
}
