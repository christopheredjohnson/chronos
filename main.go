package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

func main() {

	path := initFiles()

	db, err := sql.Open("sqlite3", path)
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

func initFiles() string {
	scope := gap.NewScope(gap.User, "chronos")

	dirs, err := scope.DataDirs()
	if err != nil || len(dirs) == 0 {
		log.Fatal("Failed finding suitable data dir")
	}

	dir := dirs[0]

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}

	path := filepath.Join(dir, "chronos.db")
	fmt.Println(path)

	return path
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
