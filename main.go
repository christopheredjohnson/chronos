package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"

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

	if err != nil {
		fmt.Println("Failed finind suitable data dir")
	}

	var dir string

	if len(dirs) > 1 {
		dir = dirs[0]
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}

	path := path.Join(dir, "chronos.db")

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
