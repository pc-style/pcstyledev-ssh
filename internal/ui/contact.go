package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pcstyle/ssh-server/internal/api"
)

// Field indices
const (
	fieldMessage = iota
	fieldName
	fieldEmail
	fieldDiscord
	fieldPhone
	fieldSubmit
	fieldBack
	fieldCount
)

// BackMsg is sent when user wants to go back
type BackMsg struct{}

// ContactModel represents the contact form
type ContactModel struct {
	inputs        []textinput.Model
	focusIndex    int
	width         int
	height        int
	apiClient     *api.Client
	submitting    bool
	submitted     bool
	submitSuccess bool
	submitMessage string
}

// NewContactModel creates a new contact form model
func NewContactModel(apiClient *api.Client) ContactModel {
	m := ContactModel{
		inputs:    make([]textinput.Model, fieldCount-2), // Exclude submit and back buttons
		apiClient: apiClient,
	}

	// Message field (required)
	m.inputs[fieldMessage] = textinput.New()
	m.inputs[fieldMessage].Placeholder = "Enter your message here..."
	m.inputs[fieldMessage].CharLimit = 2000
	m.inputs[fieldMessage].Width = 60
	m.inputs[fieldMessage].Focus()

	// Name field
	m.inputs[fieldName] = textinput.New()
	m.inputs[fieldName].Placeholder = "Your name (optional)"
	m.inputs[fieldName].CharLimit = 100
	m.inputs[fieldName].Width = 60

	// Email field
	m.inputs[fieldEmail] = textinput.New()
	m.inputs[fieldEmail].Placeholder = "your.email@example.com (optional)"
	m.inputs[fieldEmail].CharLimit = 100
	m.inputs[fieldEmail].Width = 60

	// Discord field
	m.inputs[fieldDiscord] = textinput.New()
	m.inputs[fieldDiscord].Placeholder = "@yourusername (optional)"
	m.inputs[fieldDiscord].CharLimit = 100
	m.inputs[fieldDiscord].Width = 60

	// Phone field
	m.inputs[fieldPhone] = textinput.New()
	m.inputs[fieldPhone].Placeholder = "+1234567890 (optional)"
	m.inputs[fieldPhone].CharLimit = 50
	m.inputs[fieldPhone].Width = 60

	return m
}

// SubmitMsg is sent when the form is being submitted
type SubmitMsg struct{}

// SubmitResultMsg is sent when submission completes
type SubmitResultMsg struct {
	Success bool
	Message string
	Error   error
}

// Init initializes the contact model
func (m ContactModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the contact model
func (m ContactModel) Update(msg tea.Msg) (ContactModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if !m.submitting {
				return m, func() tea.Msg {
					return BackMsg{}
				}
			}
			return m, nil

		case "tab", "shift+tab", "up", "down":
			if m.submitting {
				return m, nil
			}

			// Navigate between fields
			if msg.String() == "up" || msg.String() == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > fieldBack {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = fieldBack
			}

			// Update focus
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)

		case "enter":
			if m.submitting {
				return m, nil
			}

			// Handle submit button
			if m.focusIndex == fieldSubmit {
				// Validate message field
				if strings.TrimSpace(m.inputs[fieldMessage].Value()) == "" {
					m.submitSuccess = false
					m.submitMessage = "Message is required!"
					m.submitted = true
					return m, nil
				}

				m.submitting = true
				m.submitted = false
				return m, m.submitForm()
			}

			// Handle back button
			if m.focusIndex == fieldBack {
				return m, func() tea.Msg {
					return BackMsg{}
				}
			}
		}

	case SubmitResultMsg:
		m.submitting = false
		m.submitted = true
		m.submitSuccess = msg.Success
		if msg.Error != nil {
			m.submitMessage = msg.Error.Error()
		} else {
			m.submitMessage = msg.Message
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update the focused input
	if m.focusIndex < len(m.inputs) {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}

	return m, nil
}

// submitForm submits the contact form
func (m ContactModel) submitForm() tea.Cmd {
	return func() tea.Msg {
		req := api.ContactRequest{
			Message:  strings.TrimSpace(m.inputs[fieldMessage].Value()),
			Name:     strings.TrimSpace(m.inputs[fieldName].Value()),
			Email:    strings.TrimSpace(m.inputs[fieldEmail].Value()),
			Discord:  strings.TrimSpace(m.inputs[fieldDiscord].Value()),
			Phone:    strings.TrimSpace(m.inputs[fieldPhone].Value()),
		}

		resp, err := m.apiClient.SubmitContact(req)
		if err != nil {
			return SubmitResultMsg{
				Success: false,
				Error:   err,
			}
		}

		return SubmitResultMsg{
			Success: resp.Success,
			Message: resp.Message,
		}
	}
}

// View renders the contact form
func (m ContactModel) View() string {
	var b strings.Builder

	// Title
	title := "Contact Form"
	b.WriteString(TitleStyle.Render(title))
	b.WriteString("\n\n")

	// Show submission status
	if m.submitted {
		if m.submitSuccess {
			b.WriteString(SuccessStyle.Render("✓ " + m.submitMessage))
		} else {
			b.WriteString(ErrorStyle.Render("✗ " + m.submitMessage))
		}
		b.WriteString("\n\n")
	}

	// Show loading state
	if m.submitting {
		b.WriteString(HelpStyle.Render("Submitting..."))
		b.WriteString("\n")
		return BaseStyle.Render(b.String())
	}

	// Form fields
	fields := []string{
		"Message *",
		"Name",
		"Email",
		"Discord",
		"Phone",
	}

	for i, label := range fields {
		// Label
		labelStr := LabelStyle.Render(label + ":")
		b.WriteString(labelStr)
		b.WriteString("\n")

		// Input
		inputStyle := InputStyle
		if i == m.focusIndex {
			inputStyle = InputFocusedStyle
		}
		b.WriteString(inputStyle.Render(m.inputs[i].View()))
		b.WriteString("\n\n")
	}

	// Buttons
	submitButton := ButtonStyle.Render("Submit")
	if m.focusIndex == fieldSubmit {
		submitButton = ButtonActiveStyle.Render("Submit")
	}

	backButton := ButtonStyle.Render("Back")
	if m.focusIndex == fieldBack {
		backButton = ButtonActiveStyle.Render("Back")
	}

	b.WriteString(fmt.Sprintf("%s  %s\n", submitButton, backButton))

	// Help text
	b.WriteString("\n")
	helpText := "Use Tab/↑/↓ to navigate • Enter to submit/go back • Esc to cancel"
	b.WriteString(HelpStyle.Render(helpText))

	return BaseStyle.Render(b.String())
}
