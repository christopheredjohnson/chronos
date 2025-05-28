package main

import (
	"database/sql"
	"time"
)

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

func loadProjects(db *sql.DB) []Project {
	rows, err := db.Query("SELECT id, name, elapsed FROM projects")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var elapsedSeconds int64
		if err := rows.Scan(&p.ID, &p.Name, &elapsedSeconds); err != nil {
			continue
		}
		p.Elapsed = time.Duration(elapsedSeconds) * time.Second
		p.Tracking = false
		projects = append(projects, p)
	}
	return projects
}

func (m *model) saveProject(p Project) {
	_, _ = m.db.Exec("UPDATE projects SET elapsed = ? WHERE id = ?", int(p.Elapsed.Seconds()), p.ID)
}

func (m *model) addProject(name string) {
	res, err := m.db.Exec("INSERT INTO projects (name, elapsed) VALUES (?, ?)", name, 0)
	if err != nil {
		return
	}
	id, _ := res.LastInsertId()
	p := Project{ID: int(id), Name: name, Elapsed: 0}
	m.projects = append(m.projects, p)
}
