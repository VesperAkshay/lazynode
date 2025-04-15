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
	helpContent := `LazyNode Help

━━━ Navigation ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ↑/k, ↓/j       : Navigate up/down
  ←/h, →/l       : Navigate left/right
  alt+↑/↓/←/→    : Navigate between panels
  tab, shift+tab : Cycle through panels
  1-5            : Switch to panel by number
  
━━━ Package Management ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  i              : Install package
  I              : Install as dev dependency
  d              : Uninstall package
  u              : Update package
  o              : Check for outdated packages
  l              : Link package
  L              : Unlink package
  g              : Link package globally
  G              : Unlink package globally
  b              : Build package/project
  t              : Run tests
  p              : Publish package
  e              : Edit package.json
  /              : Search packages
  
━━━ Script Management ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  enter          : Run selected script
  ctrl+c         : Stop running script
  r              : Reload scripts list
  
━━━ NPX Commands ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  n              : New npx command
  enter          : Run selected npx command
  
━━━ Project Actions ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  r              : Reload project info
  
━━━ Panel Layout ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Left           : Scripts (top) & Project (bottom)
  Middle         : Packages
  Right          : NPX
  Bottom         : Logs
  
━━━ UI Controls ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  ?              : Toggle help
  q, ctrl+c      : Quit
  space          : Toggle details
  F              : Toggle fullscreen
  esc            : Cancel/back

Press any key to close this help screen.`

	// Apply terminal styling
	return lipgloss.NewStyle().
		Foreground(terminalBrightWhite).
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
