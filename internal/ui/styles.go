package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Brand colors (using basic ANSI colors for maximum compatibility)
	ColorPrimary   = lipgloss.AdaptiveColor{Light: "51", Dark: "51"}     // Cyan
	ColorSecondary = lipgloss.AdaptiveColor{Light: "198", Dark: "198"}   // Magenta
	ColorSuccess   = lipgloss.AdaptiveColor{Light: "46", Dark: "46"}     // Bright Green
	ColorError     = lipgloss.AdaptiveColor{Light: "196", Dark: "196"}   // Bright Red
	ColorMuted     = lipgloss.AdaptiveColor{Light: "243", Dark: "243"}   // Gray
	ColorWhite     = lipgloss.AdaptiveColor{Light: "255", Dark: "255"}   // White

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title/banner style
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Align(lipgloss.Center)

	// Navigation styles
	NavItemStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Padding(0, 2)

	NavItemSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Bold(true).
				Padding(0, 2)

	NavArrowStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	// Form styles
	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			MarginRight(1)

	InputStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(lipgloss.AdaptiveColor{Light: "235", Dark: "235"}).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(ColorPrimary).
				Background(lipgloss.AdaptiveColor{Light: "235", Dark: "235"}).
				Padding(0, 1).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(ColorPrimary)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Background(ColorSecondary).
			Padding(0, 3).
			Bold(true)

	ButtonActiveStyle = lipgloss.NewStyle().
				Foreground(ColorWhite).
				Background(ColorPrimary).
				Padding(0, 3).
				Bold(true)

	// Message styles
	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true).
			Padding(1, 2)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true).
			Padding(1, 2)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true).
			MarginTop(1)

	// Box/container style
	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			MarginTop(1)
)
