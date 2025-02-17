package cmd

import "github.com/charmbracelet/lipgloss"

var (
	errorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
)

var (
	errorStatus   = errorStyle.Render("ERROR")
	statusSuccess = successStyle.Render("SUCCESS")
)
