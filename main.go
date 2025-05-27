package main

import (
	"database/sql"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./chronos.db")
	if err != nil {
		fmt.Println("error opening db:", err)
		os.Exit(1)
	}
	initDB(db)

	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func initDB(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			elapsed INTEGER
		)
	`)
	if err != nil {
		panic(err)
	}
}
