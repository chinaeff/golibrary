package utils

import (
	"database/sql"
	"github.com/brianvoe/gofakeit"
	"log"
	"math/rand"
)

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	RentedBooks   []Book `json:"rented_books"`
	BorrowedBooks []Book `json:"borrowed_books"`
}

type Book struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Author   Author `json:"author"`
	Borrower int    `json:"borrower_id"`
}

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

type LibraryFacade struct {
	DB *sql.DB
}

func GetRandomBooks(db *sql.DB, bookIDs []int) []Book {
	count := gofakeit.Number(1, 5)
	var rentedBooks []Book

	shuffledBookIDs := shuffleIntSlice(bookIDs)
	for i := 0; i < count && i < len(shuffledBookIDs); i++ {
		bookID := shuffledBookIDs[i]
		book := getBookByID(db, bookID)
		rentedBooks = append(rentedBooks, book)
	}

	return rentedBooks
}

func GenerateAndInsertUsers(db *sql.DB, count, bookCount int) []User {
	var users []User
	for i := 0; i < count; i++ {
		user := User{
			Name:        gofakeit.Name(),
			RentedBooks: GetRandomBooks(db, generateRandomBookIDs(bookCount)),
		}

		_, err := db.Exec("INSERT INTO users (name) VALUES (?)", user.Name)
		if err != nil {
			log.Fatal(err)
		}

		lastInsertID, err := getLastInsertID(db)
		if err != nil {
			log.Fatal(err)
		}
		user.ID = int(lastInsertID)

		users = append(users, user)
	}
	return users
}

func getLastInsertID(db *sql.DB) (int64, error) {
	var lastInsertID int64
	err := db.QueryRow("SELECT last_insert_rowid()").Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func GenerateAndInsertAuthors(db *sql.DB, count, bookCount int) []Author {
	var authors []Author
	for i := 0; i < count; i++ {
		author := Author{
			Name: gofakeit.Name(),
		}

		_, err := db.Exec("INSERT INTO authors (name) VALUES (?)", author.Name)
		if err != nil {
			log.Fatal(err)
		}

		lastInsertID, err := getLastInsertID(db)
		if err != nil {
			log.Fatal(err)
		}
		author.ID = int(lastInsertID)

		author.Books = GetRandomBooks(db, generateRandomBookIDs(bookCount))

		authors = append(authors, author)
	}
	return authors
}

func GenerateAndInsertBooks(db *sql.DB, count int) []Book {
	var books []Book
	for i := 0; i < count; i++ {
		book := Book{
			Name:   gofakeit.BeerName(),
			Author: generateRandomAuthor(db),
		}

		result, err := db.Exec("INSERT INTO books (name, author_id) VALUES (?, ?)", book.Name, book.Author.ID)
		if err != nil {
			log.Fatal(err)
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		book.ID = int(lastInsertID)

		books = append(books, book)
	}
	return books
}

func generateRandomBookIDs(count int) []int {
	var bookIDs []int
	for i := 0; i < count; i++ {
		bookIDs = append(bookIDs, gofakeit.Number(1, 1000))
	}
	return bookIDs
}

func generateRandomAuthor(db *sql.DB) Author {
	authors, err := getAuthors(db)
	if err != nil {
		log.Fatal(err)
	}

	if len(authors) == 0 {
		return Author{}
	}

	author := authors[gofakeit.Number(0, len(authors)-1)]

	return author
}

func shuffleIntSlice(slice []int) []int {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	return slice
}

func getBookByID(db *sql.DB, bookID int) Book {
	var book Book
	err := db.QueryRow("SELECT id, name, author_id FROM books WHERE id = ?", bookID).Scan(&book.ID, &book.Name, &book.Author.ID)
	if err != nil {
		log.Fatal(err)
	}
	return book
}

func getAuthors(db *sql.DB) ([]Author, error) {
	rows, err := db.Query("SELECT id, name FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var author Author
		err := rows.Scan(&author.ID, &author.Name)
		if err != nil {
			return nil, err
		}

		author.Books, err = getBooksByAuthorID(db, author.ID)
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

func getBooksByAuthorID(db *sql.DB, authorID int) ([]Book, error) {
	rows, err := db.Query("SELECT id, name FROM books WHERE author_id = ?", authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Name)
		if err != nil {
			return nil, err
		}

		book.Author, err = getAuthorByID(db, authorID)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func getAuthorByID(db *sql.DB, authorID int) (Author, error) {
	var author Author
	err := db.QueryRow("SELECT id, name FROM authors WHERE id = ?", authorID).Scan(&author.ID, &author.Name)
	if err != nil {
		return Author{}, err
	}

	author.Books, err = getBooksByAuthorID(db, authorID)
	if err != nil {
		return Author{}, err
	}

	return author, nil
}
