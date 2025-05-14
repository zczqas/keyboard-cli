package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zczqas/keyboard-cli/internal/ui"
)

func main() {
	fmt.Println("Starting Keyboard CLI...")
	fmt.Println("Features:")
	fmt.Println("- F1: Switch to Visual Mode (keyboard visualization only)")
	fmt.Println("- F2: Switch to Practice Mode (typing challenge)")
	fmt.Println("- F3: Get a new typing challenge text")
	fmt.Println("- ESC/Ctrl+C: Exit")
	fmt.Println("\nPress any key to continue...")
	fmt.Scanln()

	m := ui.NewModel()

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}
