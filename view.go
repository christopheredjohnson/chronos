package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500")).Render("⏳ Chronos – Projects")
	return fmt.Sprintf("%s\n\n%s\n\n[↑/↓] select • [enter] start/stop • [q] quit", title, m.table.View())
}
