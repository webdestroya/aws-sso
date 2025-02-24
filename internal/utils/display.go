package utils

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(9))
	WarningStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(220))
	SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(40))
	HeaderStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.ANSIColor(0))
)

// var (
// 	// ErrorStatus   = ErrorStyle.Render("ERROR")
// 	// SuccessStatus = SuccessStyle.Render("SUCCESS")
// )
