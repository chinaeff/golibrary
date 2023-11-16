package main

import (
	"database/sql"
	"golibrary/utils"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestGenerateAndInsertUsers(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	utils.GenerateAndInsertUsers(db, 5, 10)

}

func TestGenerateAndInsertBooks(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	utils.GenerateAndInsertBooks(db, 5)

}
