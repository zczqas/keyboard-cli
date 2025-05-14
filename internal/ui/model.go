package ui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zczqas/keyboard-cli/internal/game"
	"github.com/zczqas/keyboard-cli/internal/keyboard"
)

// Model represents the application state
type Model struct {
	PressedKeys  map[string]time.Time
	TypedStrokes []string
	MaxStrokes   int
	Keyboard     [][]keyboard.KeyDef
	Challenge    *game.Challenge
	Mode         string // "practice" or "visual"
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

// Update handles incoming messages and updates the model state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		// Handle mode switch with F1/F2
		if msg.Type == tea.KeyF1 {
			m.Mode = "visual"
			return m, nil
		} else if msg.Type == tea.KeyF2 {
			if m.Mode != "practice" {
				m.Mode = "practice"
				m.Challenge = game.NewChallenge("medium")
			}
			return m, nil
		} else if msg.Type == tea.KeyF3 && m.Mode == "practice" {
			// Create a new challenge
			m.Challenge = game.NewChallenge("medium")
			return m, nil
		}

		var k string
		var displayKey string

		switch msg.Type {
		case tea.KeySpace:
			k = "SPACE"
			displayKey = " "
		case tea.KeyBackspace:
			k = "BACKSPACE"
			if len(m.TypedStrokes) > 0 {
				m.TypedStrokes = m.TypedStrokes[:len(m.TypedStrokes)-1]
			}
		case tea.KeyEnter:
			k = "ENTER"
			displayKey = "\n"
		default:
			k = strings.ToUpper(msg.String())
			displayKey = msg.String()
		}

		m.PressedKeys[k] = time.Now()

		// Process key for the typing challenge if in practice mode
		if m.Mode == "practice" && m.Challenge != nil && displayKey != "" && msg.Type != tea.KeyBackspace {
			m.Challenge.ProcessKey(displayKey)
		}

		if msg.Type != tea.KeyBackspace && msg.Type != tea.KeyEsc && msg.Type != tea.KeyCtrlC {
			m.TypedStrokes = append(m.TypedStrokes, displayKey)
			if len(m.TypedStrokes) > m.MaxStrokes {
				m.TypedStrokes = m.TypedStrokes[1:]
			}
		}
		return m, nil

	case tickMsg:
		now := time.Now()
		for k, t := range m.PressedKeys {
			if now.Sub(t) > 300*time.Millisecond {
				delete(m.PressedKeys, k)
			}
		}
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return m, nil
}

// View renders the current model state
func (m Model) View() string {
	var b strings.Builder

	// Display different content based on mode
	if m.Mode == "practice" && m.Challenge != nil {
		b.WriteString(renderPracticeMode(m))
	} else {
		b.WriteString(renderVisualMode(m))
	}

	b.WriteString("\n\nPress F1 for Visual Mode, F2 for Practice Mode, F3 for New Text\n")
	b.WriteString("Press ESC or Ctrl+C to exit.\n")
	return b.String()
}

// renderVisualMode renders the keyboard visualization
func renderVisualMode(m Model) string {
	var b strings.Builder

	styleNormal := lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("252"))
	styleActive := styleNormal.Background(lipgloss.Color("12")).Bold(true).Foreground(lipgloss.Color("0"))

	for _, row := range m.Keyboard {
		if len(row) > 0 {
			b.WriteString(strings.Repeat(" ", row[0].Offset))
		}
		for i, keyDef := range row {
			if i > 0 && keyDef.Offset > 0 {
				b.WriteString(strings.Repeat(" ", keyDef.Offset))
			}

			display := keyDef.Label
			if display == "" {
				display = keyDef.Key
			}

			style := styleNormal
			if _, ok := m.PressedKeys[keyDef.Key]; ok {
				style = styleActive
			}

			b.WriteString(style.Render(display))
		}
		b.WriteString("\n")
	}

	space := "SPACE"
	spaceStyle := styleNormal.Width(30).Align(lipgloss.Center)
	if _, ok := m.PressedKeys["SPACE"]; ok {
		spaceStyle = styleActive.Width(30).Align(lipgloss.Center)
	}
	b.WriteString("\n" + spaceStyle.Render(space) + "\n\n")

	typedText := strings.Join(m.TypedStrokes, "")
	b.WriteString(lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(50).
		Render("Typed: "+typedText) + "\n")

	return b.String()
}

// renderPracticeMode renders the typing challenge
func renderPracticeMode(m Model) string {
	var b strings.Builder

	// Display typing challenge text
	challengeStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Width(80)

	// Ensure we have a valid challenge
	if m.Challenge != nil {
		formattedText := m.Challenge.GetFormattedText()
		b.WriteString(challengeStyle.Render("Type this text:\n\n"+formattedText) + "\n\n")

		// Display stats
		progressStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)
		b.WriteString(progressStyle.Render(m.Challenge.GetProgressText()) + "\n\n")
	}

	// Show the keyboard visualization below
	styleNormal := lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("252"))
	styleActive := styleNormal.Background(lipgloss.Color("12")).Bold(true).Foreground(lipgloss.Color("0"))

	for _, row := range m.Keyboard {
		if len(row) > 0 {
			b.WriteString(strings.Repeat(" ", row[0].Offset))
		}
		for i, keyDef := range row {
			if i > 0 && keyDef.Offset > 0 {
				b.WriteString(strings.Repeat(" ", keyDef.Offset))
			}

			display := keyDef.Label
			if display == "" {
				display = keyDef.Key
			}

			style := styleNormal
			if _, ok := m.PressedKeys[keyDef.Key]; ok {
				style = styleActive
			}

			// Highlight the next key to press in practice mode
			if m.Challenge != nil && m.Challenge.CurrentPosition < len(m.Challenge.Text) {
				nextChar := string(m.Challenge.Text[m.Challenge.CurrentPosition])
				if (nextChar == " " && keyDef.Key == "SPACE") ||
					(strings.ToUpper(nextChar) == keyDef.Key) {
					style = styleNormal.Background(lipgloss.Color("39")).Bold(true)
				}
			}

			b.WriteString(style.Render(display))
		}
		b.WriteString("\n")
	}

	space := "SPACE"
	spaceStyle := styleNormal.Width(30).Align(lipgloss.Center)
	if _, ok := m.PressedKeys["SPACE"]; ok {
		spaceStyle = styleActive.Width(30).Align(lipgloss.Center)
	}

	// Highlight space if it's the next key
	if m.Challenge != nil && m.Challenge.CurrentPosition < len(m.Challenge.Text) {
		nextChar := string(m.Challenge.Text[m.Challenge.CurrentPosition])
		if nextChar == " " {
			spaceStyle = styleNormal.Background(lipgloss.Color("39")).Bold(true).Width(30).Align(lipgloss.Center)
		}
	}

	b.WriteString("\n" + spaceStyle.Render(space) + "\n")

	return b.String()
}

// Custom message types
type tickMsg struct{}

// NewModel creates a new Model with default settings
func NewModel() Model {
	return Model{
		PressedKeys:  make(map[string]time.Time),
		TypedStrokes: []string{},
		MaxStrokes:   100,
		Keyboard:     keyboard.GetKeyboardLayout(),
		Challenge:    game.NewChallenge("medium"),
		Mode:         "practice", // Start in practice mode by default
	}
}
