package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	db             *sql.DB
	projects       []Project
	table          table.Model
	input          textinput.Model
	addingProject  bool
	timerStartedAt map[int]time.Time
	width          int
	height         int
}

func initialModel(db *sql.DB) model {
	columns := []table.Column{
		{Title: "Project", Width: 30},
		{Title: "Elapsed", Width: 12},
		{Title: "Running", Width: 8},
	}

	projects := loadProjects(db)
	rows := buildRows(projects, make(map[int]time.Time))

	t := table.New(table.WithColumns(columns), table.WithRows(rows), table.WithFocused(true))

	// Customize styles
	styles := table.DefaultStyles()

	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFA500")). // orange border
		Foreground(lipgloss.Color("#FFFFFF")).       // white text
		Background(lipgloss.Color("#333333"))        // dark background

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("#000000")). // black text
		Background(lipgloss.Color("#FFA500"))  // orange highlight

	styles.Cell = styles.Cell.
		Foreground(lipgloss.Color("#DDDDDD")) // light gray text

	t.SetStyles(styles)

	input := textinput.New()
	input.Placeholder = "New project name"
	input.CharLimit = 255
	input.Width = 30

	return model{
		db:             db,
		projects:       projects,
		table:          t,
		input:          input,
		timerStartedAt: make(map[int]time.Time),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), textinput.Blink, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Reserve vertical space for title, help text, and spacing
		availableRows := msg.Height - 8
		if availableRows < 1 {
			availableRows = 1
		}
		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(availableRows)
		return m, nil
	case tea.KeyMsg:
		key := msg.String()

		if m.addingProject {
			switch key {
			case "enter":
				name := strings.TrimSpace(m.input.Value())
				if name != "" {
					m.addProject(name)
				}
				m.addingProject = false
				m.input.Reset()
				m.updateTableRows()
			case "esc":
				m.addingProject = false
				m.input.Reset()
			default:
				m.input, cmd = m.input.Update(msg)
				return m, cmd
			}
			return m, nil
		}

		switch key {
		case "ctrl+c", "q":
			m.saveAll()
			return m, tea.Quit
		case "a":
			m.addingProject = true
			m.input.Focus()
			return m, nil
		case "enter":
			i := m.table.Cursor()
			p := &m.projects[i]
			if p.Tracking {
				started := m.timerStartedAt[p.ID]
				p.Elapsed += time.Since(started)
				p.Tracking = false
				delete(m.timerStartedAt, p.ID)
				m.saveProject(*p)
			} else {
				m.timerStartedAt[p.ID] = time.Now()
				p.Tracking = true
			}
			m.updateTableRows()
		}

		m.table, cmd = m.table.Update(msg)
		return m, cmd

	case tickMsg:
		m.updateTableRows()
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500")).Render("⏳ Chronos – Projects")

	if m.addingProject {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n[enter] save • [esc] cancel", title, m.table.View(), m.input.View())
	}

	return fmt.Sprintf("%s\n\n%s\n\n[↑/↓] select • [enter] toggle • [a] add • [q] quit", title, m.table.View())
}

func (m *model) saveAll() {
	for _, p := range m.projects {
		if p.Tracking {
			started := m.timerStartedAt[p.ID]
			p.Elapsed += time.Since(started)
			p.Tracking = false
		}
		m.saveProject(p)
	}
}

func (m *model) updateTableRows() {
	m.table.SetRows(buildRows(m.projects, m.timerStartedAt))
}

type Project struct {
	ID       int
	Name     string
	Elapsed  time.Duration
	Tracking bool
}

type tickMsg time.Time

func buildRows(projects []Project, startedMap map[int]time.Time) []table.Row {
	var rows []table.Row
	for _, p := range projects {
		elapsed := formatDuration(func() time.Duration {
			if p.Tracking {
				return time.Since(startedMap[p.ID]) + p.Elapsed
			}
			return p.Elapsed
		}())

		status := "⏸"

		if p.Tracking {
			status = "⏱"
		}
		rows = append(rows, table.Row{p.Name, elapsed, status})
	}
	return rows
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
