package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpPanel represents the help panel
type HelpPanel struct {
	width  int
	height int
}

// NewHelpPanel creates a new help panel
func NewHelpPanel() *HelpPanel {
	return &HelpPanel{}
}

// Init initializes the help panel
func (p *HelpPanel) Init() tea.Cmd {
	return nil
}

// Update handles key events
func (p *HelpPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	return p, nil
}

// SetSize sets the size of the help panel
func (p *HelpPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// View returns the view for the help panel
func (p *HelpPanel) View() string {
	helpContent := `LazyNode Help:

Navigation:
  q           : Quit
  ?           : Toggle help
  tab         : Switch panels
  ↑/↓         : Navigate within panel
  alt+↑/↓     : Switch between panels
  enter       : Select/run script
  
Scripts:
  r           : Refresh scripts
  
Packages:
  a           : Show all actions
  i           : Install package
  d           : Uninstall package
  u           : Update package
  o           : Check for outdated
`

	// Render without borders or padding
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Render(helpContent)
}

// Width returns the width of the help panel
func (p *HelpPanel) Width() int {
	return p.width
}

// Height returns the height of the help panel
func (p *HelpPanel) Height() int {
	return p.height
}

// Title returns the title of the help panel
func (p *HelpPanel) Title() string {
	return "Help"
}
