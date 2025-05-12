package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zczqas/keyboard-cli/internal/ui"
)

func main() {
	m := ui.NewModel()

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
