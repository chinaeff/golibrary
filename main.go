package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golibrary/library"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTablesSQL := `
		CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);

		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			author_id INTEGER,
			FOREIGN KEY (author_id) REFERENCES authors(id)
		);

		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);
	`

	_, err = db.Exec(createTablesSQL)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Таблицы успешно созданы.")
	}

	lf := library.NewLibraryFacade(db)

	lf.StartLibrary()
	lf.PrintLibraryUsers()

	fmt.Println("Библиотека запущена и готова к использованию!")
}

///
