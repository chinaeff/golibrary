package main

// @title Library API
// @version 1.0
// @description API for managing a library.
// @termsOfService http://library.com/terms/
// @contact.name API Support
// @contact.url http://library.com/support
// @host localhost:8080
// @BasePath /v1
// @schemes http

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"golibrary/library"
	"log"
	"net/http"
)

// @router /users [get]
// @summary Get the list of users
// @description Get the list of all users in the library
// @tags Users
// @produce json
// @success 200 {array} User

// RunSwaggerCmd генерирует документацию Swagger JSON
func RunSwaggerCmd() {
	swag.NewFormatter()
}

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

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.Handle("/swagger/doc.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RunSwaggerCmd()
		http.ServeFile(w, r, "docs/doc.json")
	}))

	router.HandleFunc("/authors", lf.GetAuthorsHandler).Methods("GET")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
