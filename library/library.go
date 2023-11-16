package library

import (
	"database/sql"
	"fmt"
	"golibrary/utils"
	"log"
)

type LibraryFacade struct {
	DB *sql.DB
}

func NewLibraryFacade(db *sql.DB) *LibraryFacade {
	return &LibraryFacade{DB: db}
}

func (lf *LibraryFacade) StartLibrary() {
	userCount, bookCount, authorCount := lf.getCounts()
	lf.generateAndInsertDataIfNeeded(userCount, bookCount, authorCount)
}

func (lf *LibraryFacade) PrintLibraryUsers() {
	bookIDs, err := lf.getBookIDs()
	if err != nil {
		log.Println("Ошибка при получении идентификаторов книг:", err)
		return
	}

	users, err := lf.getUsers(bookIDs)
	if err != nil {
		log.Println("Ошибка при получении пользователей:", err)
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Name)
		fmt.Println("Арендованные книги:")
		for _, book := range user.RentedBooks {
			fmt.Printf("  ID: %d, Название: %s\n", book.ID, book.Name)
		}
		fmt.Println("---------------")
	}
}

func (lf *LibraryFacade) getCounts() (int, int, int) {
	var userCount, bookCount, authorCount int

	err := lf.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		log.Fatal(err)
	}

	err = lf.DB.QueryRow("SELECT COUNT(*) FROM books").Scan(&bookCount)
	if err != nil {
		log.Fatal(err)
	}

	err = lf.DB.QueryRow("SELECT COUNT(*) FROM authors").Scan(&authorCount)
	if err != nil {
		log.Fatal(err)
	}

	return userCount, bookCount, authorCount
}

func (lf *LibraryFacade) generateAndInsertDataIfNeeded(userCount, bookCount, authorCount int) (string, error) {
	var result string

	if userCount == 0 {
		users := utils.GenerateAndInsertUsers(lf.DB, 50, bookCount)
		result += fmt.Sprintf("Сгенерировано и добавлено пользователей: %d\n", len(users))
	}

	if bookCount == 0 {
		books := utils.GenerateAndInsertBooks(lf.DB, 100)
		result += fmt.Sprintf("Сгенерировано и добавлено книг: %d\n", len(books))
	}

	if authorCount == 0 {
		authors := utils.GenerateAndInsertAuthors(lf.DB, 10, bookCount)

		result += fmt.Sprintf("Сгенерировано и добавлено авторов: %d\n", len(authors))
	}

	return result, nil
}

func (lf *LibraryFacade) getBookIDs() ([]int, error) {
	rows, err := lf.DB.Query("SELECT id FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookIDs []int
	for rows.Next() {
		var bookID int
		err := rows.Scan(&bookID)
		if err != nil {
			return nil, err
		}
		bookIDs = append(bookIDs, bookID)
	}

	return bookIDs, nil
}

func (lf *LibraryFacade) getUsers(bookIDs []int) ([]utils.User, error) {
	rows, err := lf.DB.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []utils.User
	for rows.Next() {
		var user utils.User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}

		user.RentedBooks = utils.GetRandomBooks(lf.DB, bookIDs)

		users = append(users, user)
	}

	return users, nil
}
