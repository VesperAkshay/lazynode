package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lazynode/lazynode/pkg/ui"
)

// Version information set during build
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	// Define command-line flags
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("LazyNode version %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildDate)
		os.Exit(0)
	}

	// Initialize the LazyNode application
	fmt.Printf("Starting LazyNode %s - TUI for Node.js, npm, and npx\n", Version)

	// Create our model
	model := ui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running LazyNode: %v\n", err)
		os.Exit(1)
	}
}
