package ui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/VesperAkshay/lazynode/pkg/npm"
	"github.com/VesperAkshay/lazynode/pkg/npx"
	"github.com/VesperAkshay/lazynode/pkg/project"
	"github.com/VesperAkshay/lazynode/pkg/scripts"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines the keybindings for the application
type KeyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	PanelUp     key.Binding
	PanelDown   key.Binding
	Enter       key.Binding
	Install     key.Binding
	Delete      key.Binding
	Outdated    key.Binding
	Update      key.Binding
	Reload      key.Binding
	Edit        key.Binding
	TabScripts  key.Binding
	TabPackages key.Binding
	TabProject  key.Binding
	TabNpx      key.Binding
	TabLogs     key.Binding
	Quit        key.Binding
	Help        key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
// This is part of the help.KeyMap interface.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
// This is part of the help.KeyMap interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.PanelUp, k.PanelDown},
		{k.TabScripts, k.TabPackages, k.TabProject, k.TabNpx, k.TabLogs},
		{k.Install, k.Delete, k.Outdated, k.Update},
		{k.Enter, k.Reload, k.Edit},
		{k.Help, k.Quit},
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		PanelUp: key.NewBinding(
			key.WithKeys("alt+up", "alt+k"),
			key.WithHelp("alt+↑", "previous panel"),
		),
		PanelDown: key.NewBinding(
			key.WithKeys("alt+down", "alt+j"),
			key.WithHelp("alt+↓", "next panel"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/run"),
		),
		Install: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "install package"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete package"),
		),
		Outdated: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "check outdated"),
		),
		Update: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "update package"),
		),
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload UI"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit package.json"),
		),
		TabScripts: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "scripts panel"),
		),
		TabPackages: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "packages panel"),
		),
		TabProject: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "project panel"),
		),
		TabNpx: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "npx panel"),
		),
		TabLogs: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "logs panel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}

// Model contains the state of the application
type Model struct {
	keys         KeyMap
	help         help.Model
	width        int
	height       int
	activeTab    string
	panels       map[string]Panel
	projectPath  string
	project      *project.Project
	packageMgr   *npm.PackageManager
	scriptRunner *scripts.ScriptRunner
	npxRunner    *npx.Runner
	logs         *LogsPanel
	helpPanel    *HelpPanel
	showHelp     bool
	ready        bool
	error        string
}

// NewModel initializes a new model with default values
func NewModel() Model {
	keys := DefaultKeyMap()
	helpPanel := NewHelpPanel()

	return Model{
		keys:      keys,
		help:      help.New(),
		activeTab: "scripts",
		panels:    make(map[string]Panel),
		helpPanel: helpPanel,
		showHelp:  false,
		ready:     false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.detectProject,
		m.startTicker(),
	)
}

// detectProject is a command that tries to find a package.json file
func (m Model) detectProject() tea.Msg {
	packageJSONPath, err := project.Detect()
	if err != nil {
		return errorMsg(fmt.Sprintf("Could not find a package.json file: %v", err))
	}

	// Initialize the project
	proj, err := project.NewProject(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error loading project: %v", err))
	}

	// Initialize package manager
	pkgMgr, err := npm.NewPackageManager(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error initializing package manager: %v", err))
	}

	// Initialize script runner
	scriptRunner, err := scripts.NewScriptRunner(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error initializing script runner: %v", err))
	}

	// Initialize npx runner
	npxRunner, err := npx.NewRunner(filepath.Dir(packageJSONPath))
	if err != nil {
		return errorMsg(fmt.Sprintf("Error initializing npx runner: %v", err))
	}

	return projectDetectedMsg{
		path:         packageJSONPath,
		project:      proj,
		packageMgr:   pkgMgr,
		scriptRunner: scriptRunner,
		npxRunner:    npxRunner,
	}
}

// errorMsg represents an error message
type errorMsg string

// projectDetectedMsg is sent when a project is detected
type projectDetectedMsg struct {
	path         string
	project      *project.Project
	packageMgr   *npm.PackageManager
	scriptRunner *scripts.ScriptRunner
	npxRunner    *npx.Runner
}

// startTicker creates a ticker for real-time updates
func (m Model) startTicker() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// tickMsg is a message sent by the ticker
type tickMsg time.Time

// Update handles key events and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Update the model dimensions
		m.width = msg.Width
		m.height = msg.Height

		// Update panel dimensions when terminal size changes
		if m.ready {
			// Get actual terminal dimensions
			termWidth := m.width
			termHeight := m.height

			// Ensure minimum dimensions
			if termWidth < 80 {
				termWidth = 80
			}
			if termHeight < 24 {
				termHeight = 24
			}

			// Reserve space for top and bottom bars
			availableHeight := termHeight - 2

			// Calculate panel dimensions - 4 equal panels in a grid (2x2)
			topRowHeight := availableHeight * 4 / 5
			if topRowHeight < 8 {
				topRowHeight = 8
			}

			bottomRowHeight := availableHeight - topRowHeight - 1 // -1 for gap
			if bottomRowHeight < 3 {
				bottomRowHeight = 3
			}

			// Calculate panel widths - divide available width into 2 columns
			columnWidth := termWidth/2 - 1 // -1 for gap between columns
			if columnWidth < 35 {
				columnWidth = 35
			}

			// Ensure we don't exceed available width
			if columnWidth*2+1 > termWidth {
				columnWidth = (termWidth - 1) / 2
			}

			// Update all panel sizes based on their position in the layout
			for name, panel := range m.panels {
				if name == "logs" {
					panel.SetSize(termWidth-2, bottomRowHeight-2)
				} else {
					// All other panels get equal size in the grid
					panel.SetSize(columnWidth-2, topRowHeight/2-2)
				}
			}

			// Update help panel size
			if m.helpPanel != nil {
				m.helpPanel.SetSize(termWidth, termHeight)
			}

			// Update help model width
			m.help.Width = termWidth
		}

		return m, nil

	case tea.KeyMsg:
		// Check if any panel has active input or confirmation dialog - if so, pass the event to that panel
		if m.ready && !m.showHelp {
			if panel, ok := m.panels[m.activeTab]; ok {
				// Check if panel has active input/dialogs (for PackagesPanel)
				if packagesPanel, ok := panel.(*PackagesPanel); ok &&
					(packagesPanel.showInput || packagesPanel.showActions || packagesPanel.showConfirm) {
					updatedPanel, cmd := panel.Update(msg)
					m.panels[m.activeTab] = updatedPanel
					return m, cmd
				}
			}
		}

		// Handle global key presses
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, m.keys.TabScripts) && m.ready:
			m.activeTab = "scripts"
			return m, nil

		case key.Matches(msg, m.keys.TabPackages) && m.ready:
			m.activeTab = "packages"
			return m, nil

		case key.Matches(msg, m.keys.TabProject) && m.ready:
			m.activeTab = "project"
			return m, nil

		case key.Matches(msg, m.keys.TabNpx) && m.ready:
			m.activeTab = "npx"
			return m, nil

		case key.Matches(msg, m.keys.TabLogs) && m.ready:
			m.activeTab = "logs"
			return m, nil

		case key.Matches(msg, m.keys.Reload) && m.ready:
			// Reload the project
			return m, m.detectProject
		}

		// Pass Tab key to switch panels
		if msg.String() == "tab" && m.ready {
			panels := []string{"scripts", "packages", "project", "npx", "logs"}
			for i, panel := range panels {
				if panel == m.activeTab {
					// Move to the next panel
					nextIndex := (i + 1) % len(panels)
					m.activeTab = panels[nextIndex]
					break
				}
			}
			return m, nil
		}

		// Check if Up/Down keys are for panel selection (left side) or content navigation (right side)
		// This happens when Alt/Option key is held with Up/Down
		if m.ready && !m.showHelp {
			switch {
			case key.Matches(msg, m.keys.PanelUp):
				// Navigate panels up
				panels := []string{"scripts", "packages", "project", "npx", "logs"}
				for i, panel := range panels {
					if panel == m.activeTab {
						// Move to the previous panel
						prevIndex := i - 1
						if prevIndex < 0 {
							prevIndex = len(panels) - 1
						}
						m.activeTab = panels[prevIndex]
						break
					}
				}
				return m, nil

			case key.Matches(msg, m.keys.PanelDown):
				// Navigate panels down
				panels := []string{"scripts", "packages", "project", "npx", "logs"}
				for i, panel := range panels {
					if panel == m.activeTab {
						// Move to the next panel
						nextIndex := (i + 1) % len(panels)
						m.activeTab = panels[nextIndex]
						break
					}
				}
				return m, nil
			}
		}

		// Pass the key event to the active panel
		if m.ready && !m.showHelp {
			if panel, ok := m.panels[m.activeTab]; ok {
				updatedPanel, cmd := panel.Update(msg)
				m.panels[m.activeTab] = updatedPanel
				return m, cmd
			}
		}

	case errorMsg:
		m.error = string(msg)
		return m, nil

	case projectDetectedMsg:
		// Save the project info
		m.projectPath = msg.path
		m.project = msg.project
		m.packageMgr = msg.packageMgr
		m.scriptRunner = msg.scriptRunner
		m.npxRunner = msg.npxRunner

		// Create logs panel first
		m.logs = NewLogsPanel()

		// Create panels
		scriptsPanel := NewScriptsPanel(m.scriptRunner)
		scriptsPanel.SetLogsPanel(m.logs)
		m.panels["scripts"] = scriptsPanel

		// Ensure packages are loaded before creating the package panel
		if err := m.packageMgr.LoadPackages(); err != nil {
			m.logs.AddLog(fmt.Sprintf("Warning: Failed to load packages: %v", err))
		}

		m.panels["packages"] = NewPackagesPanel(m.packageMgr)
		m.panels["project"] = NewProjectPanel(m.project)
		m.panels["npx"] = NewNpxPanel(m.npxRunner, m.logs)
		m.panels["logs"] = m.logs

		// Initialize panel sizes
		contentWidth := m.width - 6
		contentHeight := m.height - 8

		for _, panel := range m.panels {
			panel.SetSize(contentWidth, contentHeight)
		}

		m.ready = true

		// Add welcome message to logs
		m.logs.AddLog(fmt.Sprintf("Welcome to LazyNode v1.0 - Managing project: %s", m.project.Name))
		m.logs.AddLog(fmt.Sprintf("Found %d packages in package.json", len(m.packageMgr.Packages)))
		m.logs.AddLog("Use tabs 1-5 to navigate between panels")
		m.logs.AddLog("Press ? for help")

	case tickMsg:
		// Real-time updates
		if m.ready {
			// Check for active operations and update UI
			var activeOperation bool

			// Auto-refresh logs panel in real-time
			if logsPanel, ok := m.panels["logs"].(*LogsPanel); ok {
				logsPanel.Update(nil)
			}

			// Update script status in real-time
			if scriptsPanel, ok := m.panels["scripts"].(*ScriptsPanel); ok {
				if scriptsPanel.activeScript != "" {
					scriptsPanel.Update(nil)
					activeOperation = true
					// Add log for long-running scripts periodically
					if time.Now().Second()%5 == 0 {
						m.logs.AddLog(fmt.Sprintf("Script still running: %s", scriptsPanel.activeScript))
					}
				}
			}

			// Update package installation status periodically
			if packagesPanel, ok := m.panels["packages"].(*PackagesPanel); ok {
				if packagesPanel.loading {
					updatedPanel, _ := packagesPanel.Update(nil)
					m.panels["packages"] = updatedPanel
					activeOperation = true
				}
			}

			// If there's an active operation, reduce tick interval for more responsive UI
			if activeOperation {
				cmds = append(cmds, tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
					return tickMsg(t)
				}))
			} else {
				cmds = append(cmds, m.startTicker())
			}
		} else {
			cmds = append(cmds, m.startTicker())
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the user interface
func (m Model) View() string {
	if m.error != "" {
		return fmt.Sprintf("Error: %s\n\nPress q to quit", m.error)
	}

	if !m.ready {
		return "Loading LazyNode..."
	}

	// Show help panel if requested
	if m.showHelp {
		return m.helpPanel.View()
	}

	// Define styles for different panels
	topBarStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ebdbb2")).
		Background(lipgloss.Color("#3c3836")).
		Padding(0, 1)

	statusStyle := topBarStyle.Copy()

	// Panel style similar to Lazygit
	panelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4d4d4d")).
		Padding(0, 0).
		Margin(0, 0)

	selectedPanelStyle := panelStyle.Copy().
		BorderForeground(lipgloss.Color("#b8bb26"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fabd2f")).
		Padding(0, 1)

	// Get actual terminal dimensions
	termWidth := m.width
	termHeight := m.height

	// Ensure minimum dimensions
	if termWidth < 80 {
		termWidth = 80
	}
	if termHeight < 24 {
		termHeight = 24
	}

	// Reserve space for top and bottom bars
	availableHeight := termHeight - 2

	// Calculate panel dimensions - 4 equal panels in a grid (2x2)
	topRowHeight := availableHeight * 4 / 5
	if topRowHeight < 8 {
		topRowHeight = 8
	}

	bottomRowHeight := availableHeight - topRowHeight - 1 // -1 for gap
	if bottomRowHeight < 3 {
		bottomRowHeight = 3
	}

	// Calculate panel widths - divide available width into 2 columns
	columnWidth := termWidth/2 - 1 // -1 for gap between columns
	if columnWidth < 35 {
		columnWidth = 35
	}

	// Ensure we don't exceed available width
	if columnWidth*2+1 > termWidth {
		columnWidth = (termWidth - 1) / 2
	}

	// Top status bar
	statusBar := topBarStyle.Width(termWidth).
		Render(fmt.Sprintf("LazyNode - %s", m.project.Name))

	// Prepare all panels
	panelContents := make(map[string]string)
	panelTitles := make(map[string]string)

	// Set panel sizes and get content
	for name, panel := range m.panels {
		if name == "logs" {
			panel.SetSize(termWidth-2, bottomRowHeight-2)
		} else {
			// All other panels get equal size in the grid
			panel.SetSize(columnWidth-2, topRowHeight/2-2)
		}

		// Get panel content and title
		panelContents[name] = panel.View()
		panelTitles[name] = panel.Title()
	}

	// Render panel with proper styling and highlighting active panel
	renderPanel := func(name string, width, height int) string {
		content := fmt.Sprintf("%s\n%s",
			titleStyle.Render(panelTitles[name]),
			panelContents[name])

		style := panelStyle
		if name == m.activeTab {
			style = selectedPanelStyle
		}

		return style.
			Width(width).
			Height(height).
			Render(content)
	}

	// Render all panels
	topLeftRendered := renderPanel("scripts", columnWidth, topRowHeight/2)
	topRightRendered := renderPanel("packages", columnWidth, topRowHeight/2)
	bottomLeftRendered := renderPanel("project", columnWidth, topRowHeight/2)
	bottomRightRendered := renderPanel("npx", columnWidth, topRowHeight/2)
	logsRendered := renderPanel("logs", termWidth, bottomRowHeight)

	// Layout the panels
	topLeftRow := lipgloss.JoinHorizontal(lipgloss.Top, topLeftRendered, topRightRendered)
	bottomLeftRow := lipgloss.JoinHorizontal(lipgloss.Top, bottomLeftRendered, bottomRightRendered)
	topGrid := lipgloss.JoinVertical(lipgloss.Left, topLeftRow, bottomLeftRow)

	// Help bar at the bottom
	helpText := "[q]Quit [?]Help [Tab]Switch panels [1-5]Select panel [↵]Select"
	if m.activeTab == "scripts" {
		helpText += " | [r]Run script"
	} else if m.activeTab == "packages" {
		helpText += " | [i]Install [d]Delete [u]Update"
	}
	helpBar := statusStyle.Width(termWidth).Render(helpText)

	// Complete layout
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		topGrid,
		logsRendered,
		helpBar,
	)

	// Final trim to ensure we don't exceed terminal dimensions
	return lipgloss.NewStyle().
		MaxWidth(termWidth).
		MaxHeight(termHeight).
		Render(ui)
}
