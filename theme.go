package main

import "github.com/charmbracelet/lipgloss"

const (
	PRIMARY = lipgloss.Color("#E83151")
)

type Theme struct {
	TitleStyle    lipgloss.Style
	HeaderStyle   lipgloss.Style
	SelectedStyle lipgloss.Style
}

func defaultTheme() Theme {
	return Theme{
		TitleStyle:    lipgloss.NewStyle().Bold(true).Foreground(PRIMARY),
		HeaderStyle:   lipgloss.NewStyle().Foreground(PRIMARY),
		SelectedStyle: lipgloss.NewStyle().Background(PRIMARY),
	}
}
