package utils

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
)

var (
	ErrorStatus   = ErrorStyle.Render("ERROR")
	SuccessStatus = SuccessStyle.Render("SUCCESS")
)
