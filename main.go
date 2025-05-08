package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyDef struct {
	key    string
	label  string
	offset int
}

var keyboard = [][]keyDef{
	// Top row (no offset)
	{{key: "Q"}, {key: "W"}, {key: "E"}, {key: "R"}, {key: "T"}, {key: "Y"}, {key: "U"}, {key: "I"}, {key: "O"}, {key: "P"}, {key: "[", label: "["}, {key: "]", label: "]"}},
	// Home row
	{{key: "A", offset: 2}, {key: "S"}, {key: "D"}, {key: "F"}, {key: "G"}, {key: "H"}, {key: "J"}, {key: "K"}, {key: "L"}, {key: ";", label: ";"}, {key: "'", label: "'"}},
	// Bottom row
	{{key: "Z", offset: 4}, {key: "X"}, {key: "C"}, {key: "V"}, {key: "B"}, {key: "N"}, {key: "M"}, {key: ",", label: ","}, {key: ".", label: "."}, {key: "/", label: "/"}},
}

type model struct {
	pressedKeys map[string]time.Time
}

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg{}
		}),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

		var k string
		switch msg.Type {
		case tea.KeySpace:
			k = "SPACE"
		default:
			k = strings.ToUpper(msg.String())
		}

		m.pressedKeys[k] = time.Now()
		return m, nil
	case tickMsg:
		now := time.Now()
		for k, t := range m.pressedKeys {
			if now.Sub(t) > 300*time.Millisecond {
				delete(m.pressedKeys, k)
			}
		}
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg{}
		})
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	styleNormal := lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("252"))
	styleActive := styleNormal.Background(lipgloss.Color("12")).Bold(true).Foreground(lipgloss.Color("0"))

	for _, row := range keyboard {
		if len(row) > 0 {
			b.WriteString(strings.Repeat(" ", row[0].offset))
		}
		for i, keyDef := range row {
			if i > 0 && keyDef.offset > 0 {
				b.WriteString(strings.Repeat(" ", keyDef.offset))
			}

			display := keyDef.label
			if display == "" {
				display = keyDef.key
			}

			style := styleNormal
			if _, ok := m.pressedKeys[keyDef.key]; ok {
				style = styleActive
			}

			b.WriteString(style.Render(display))
		}
		b.WriteString("\n")
	}

	space := "SPACE"
	spaceStyle := styleNormal.Width(30).Align(lipgloss.Center)
	if _, ok := m.pressedKeys["SPACE"]; ok {
		spaceStyle = styleActive.Width(30).Align(lipgloss.Center)
	}
	b.WriteString("\n" + spaceStyle.Render(space) + "\n\n")

	b.WriteString("Press ESC or Ctrl+C to exit.\n")
	return b.String()
}

func main() {
	m := model{pressedKeys: make(map[string]time.Time)}

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	fmt.Println("Starting Keyboard Visualizer...")
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
