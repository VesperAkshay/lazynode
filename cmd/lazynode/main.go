package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lazynode/lazynode/pkg/ui"
)

func main() {
	// Initialize the LazyNode application
	fmt.Println("Starting LazyNode - TUI for Node.js, npm, and npx")

	// Create our model
	model := ui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running LazyNode: %v\n", err)
		os.Exit(1)
	}
}
