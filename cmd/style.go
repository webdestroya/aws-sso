package cmd

import "github.com/charmbracelet/lipgloss"

var errorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("9"))

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("10"))
