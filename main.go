package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	gap "github.com/muesli/go-app-paths"
)

func main() {

	exportFlag := flag.String("export", "", "Export data (csv or json)")

	flag.Parse()

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	path := initFiles()

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		fmt.Println("error opening db:", err)
		os.Exit(1)
	}
	initDB(db)

	if *exportFlag != "" {
		projects := loadProjects(db)
		switch *exportFlag {
		case "csv":
			err := exportToCSV(projects, "chronos_export.csv")
			if err != nil {
				log.Fatalf("Error exporting CSV: %v", err)
			}
			fmt.Println("Exported to chronos_export.csv")
		case "json":
			err := exportToJSON(projects, "chronos_export.json")
			if err != nil {
				log.Fatalf("Error exporting JSON: %v", err)
			}
			fmt.Println("Exported to chronos_export.json")
		default:
			log.Fatalf("Unknown export format: %s", *exportFlag)
		}
		return
	}

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

	return path
}
