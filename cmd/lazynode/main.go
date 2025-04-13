package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/VesperAkshay/lazynode/pkg/ui"
	"github.com/VesperAkshay/lazynode/pkg/version"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Define command-line flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// If version flag is set, print version and exit
	if *versionFlag {
		fmt.Printf("LazyNode version %s\n", version.GetVersion())
		os.Exit(0)
	}

	// Initialize the LazyNode application
	fmt.Printf("Starting LazyNode %s - TUI for Node.js, npm, and npx\n", version.GetVersion())

	// Create our model
	model := ui.NewModel()
	p := tea.NewProgram(model, tea.WithAltScreen())

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running LazyNode: %v\n", err)
		os.Exit(1)
	}
}
