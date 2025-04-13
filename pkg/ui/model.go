package ui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lazynode/lazynode/pkg/npm"
	"github.com/lazynode/lazynode/pkg/npx"
	"github.com/lazynode/lazynode/pkg/project"
	"github.com/lazynode/lazynode/pkg/scripts"
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

		// Check if Up/Down keys are for panel selection (left side) or content navigation (right side)
		// This happens when Alt/Option key is held with Up/Down
		if m.ready && !m.showHelp && (key.Matches(msg, m.keys.PanelUp) || key.Matches(msg, m.keys.PanelDown)) {
			// Navigation between panels in left side
			panels := []string{"scripts", "packages", "project", "npx"}
			if key.Matches(msg, m.keys.PanelUp) {
				for i, panel := range panels {
					if panel == m.activeTab && i > 0 {
						m.activeTab = panels[i-1]
						return m, nil
					}
				}
			} else { // PanelDown
				for i, panel := range panels {
					if panel == m.activeTab && i < len(panels)-1 {
						m.activeTab = panels[i+1]
						return m, nil
					}
				}
			}
			return m, nil
		}

		// For regular Up/Down keys, pass them to the panel for content navigation
		if m.ready && !m.showHelp {
			if panel, ok := m.panels[m.activeTab]; ok {
				updatedPanel, cmd := panel.Update(msg)
				m.panels[m.activeTab] = updatedPanel
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		// Store the window size for responsive rendering
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

		// Always update help panel size
		if m.helpPanel != nil {
			m.helpPanel.SetSize(m.width, m.height)
		}

		// Update panel sizes
		contentWidth := m.width - 6
		contentHeight := m.height - 8

		if m.ready {
			for _, panel := range m.panels {
				panel.SetSize(contentWidth, contentHeight)
			}
		}

	case errorMsg:
		m.error = string(msg)

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

	// Update panel dimensions after updating the model dimensions
	if m.width > 0 && m.height > 0 && m.ready {
		// Calculate heights with minimal padding
		mainHeight := int(float64(m.height)*0.85) - 1 // -1 for header
		if mainHeight < 10 {
			mainHeight = 10
		}
		logsHeight := m.height - mainHeight - 1
		if logsHeight < 3 {
			logsHeight = 3
		}

		// Set panel sizes with full width
		m.panels["scripts"].SetSize(m.width, mainHeight)
		m.panels["packages"].SetSize(m.width, mainHeight)
		m.panels["logs"].SetSize(m.width, logsHeight)

		// Always set help panel size regardless of ready state
		if m.helpPanel != nil {
			m.helpPanel.SetSize(m.width, m.height)
		}
		m.help.Width = m.width
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
		Width(m.width).
		Padding(0, 1)

	statusStyle := topBarStyle.Copy()

	// Panel style similar to Lazygit
	panelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4d4d4d")).
		Padding(0, 0).
		Margin(0, 0)

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#b8bb26"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fabd2f")).
		Padding(0, 1)

	// Calculate panel dimensions
	leftPanelWidth := m.width / 4
	if leftPanelWidth < 20 {
		leftPanelWidth = 20
	}

	rightPanelWidth := m.width - leftPanelWidth - 3 // Account for borders

	topPanelsHeight := m.height * 2 / 3
	if topPanelsHeight < 10 {
		topPanelsHeight = 10
	}

	bottomPanelHeight := m.height - topPanelsHeight - 4 // Account for status bars and borders
	if bottomPanelHeight < 5 {
		bottomPanelHeight = 5
	}

	// Top status bar
	statusBar := topBarStyle.Render(fmt.Sprintf("LazyNode - %s", m.project.Name))

	// Left panel (navigation)
	leftPanelTitle := titleStyle.Render("Panel Selection (Alt+↑/↓)")
	leftPanelContent := fmt.Sprintf("%s\n\n", leftPanelTitle)
	panels := []string{"scripts", "packages", "project", "npx"}
	panelNames := []string{"Scripts", "Packages", "Project", "NPX"}

	for i, panel := range panels {
		if panel == m.activeTab {
			leftPanelContent += selectedStyle.Render(fmt.Sprintf("▶ %s\n", panelNames[i]))
		} else {
			leftPanelContent += fmt.Sprintf("  %s\n", panelNames[i])
		}
	}

	leftPanel := panelStyle.
		Width(leftPanelWidth).
		Height(topPanelsHeight).
		Render(leftPanelContent)

	// Right panel (content)
	var mainContent string
	var rightPanelTitle string

	if panel, ok := m.panels[m.activeTab]; ok {
		// Temporarily set the panel size to fit in our layout
		panel.SetSize(rightPanelWidth-2, topPanelsHeight-2) // Adjust for borders
		mainContent = panel.View()
		rightPanelTitle = panel.Title()
	}

	// Add a custom header to the right panel
	rightPanelContent := fmt.Sprintf("%s (↑/↓ to navigate)\n%s",
		titleStyle.Render(rightPanelTitle),
		mainContent,
	)

	rightPanel := panelStyle.
		Width(rightPanelWidth).
		Height(topPanelsHeight).
		Render(rightPanelContent)

	// Command log title with real-time animation
	activeSpinner := ""

	// Add spinner animation for active operations
	if m.activeTab == "packages" {
		if packagesPanel, ok := m.panels["packages"].(*PackagesPanel); ok && packagesPanel.loading {
			spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26"))
			activeSpinner = spinnerStyle.Render(packagesPanel.spinnerFrames[packagesPanel.spinner]) + " "
		}
	} else if m.activeTab == "scripts" {
		if scriptsPanel, ok := m.panels["scripts"].(*ScriptsPanel); ok && scriptsPanel.activeScript != "" {
			spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26"))
			// Use logs panel's spinner for animation
			if logsPanel, ok := m.panels["logs"].(*LogsPanel); ok {
				activeSpinner = spinnerStyle.Render(logsPanel.spinnerFrames[logsPanel.spinner]) + " "
			}
		}
	}

	commandLogTitle := titleStyle.Render("Command Log " + activeSpinner)

	// Show the active command if one is running, or just a prompt
	commandLogContent := ""

	// Check if there's an active input in packages panel
	if m.activeTab == "packages" {
		if packagesPanel, ok := m.panels["packages"].(*PackagesPanel); ok && packagesPanel.showInput {
			inputMode := ""
			switch packagesPanel.inputMode {
			case "install":
				inputMode = "install"
			case "install-dev":
				inputMode = "install --save-dev"
			case "uninstall":
				inputMode = "uninstall"
			case "update":
				inputMode = "update"
			}
			commandLogContent = fmt.Sprintf("$ npm %s %s", inputMode, packagesPanel.input.Value())
		} else if packagesPanel.loading {
			commandLogContent = "$ npm install ... working"
		} else {
			commandLogContent = "$ npm install <package>"
		}
	} else if m.activeTab == "scripts" {
		if scriptsPanel, ok := m.panels["scripts"].(*ScriptsPanel); ok && scriptsPanel.activeScript != "" {
			commandLogContent = fmt.Sprintf("$ npm run %s", scriptsPanel.activeScript)
		} else {
			commandLogContent = "$ npm run <script>"
		}
	} else if m.activeTab == "npx" {
		if npxPanel, ok := m.panels["npx"].(*NpxPanel); ok && npxPanel.showInput {
			commandLogContent = fmt.Sprintf("$ npx %s", npxPanel.input.Value())
		} else {
			commandLogContent = "$ npx <command>"
		}
	} else {
		commandLogContent = "$ npm <command>"
	}

	// Format the logs panel with a command log section
	commandLogStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#21262d")).
		Foreground(lipgloss.Color("#c9d1d9")).
		Width(m.width-4).
		Height(3).
		Padding(0, 1)

	commandLogFormatted := lipgloss.JoinVertical(
		lipgloss.Left,
		commandLogTitle,
		commandLogStyle.Render(commandLogContent),
	)

	// Bottom panel (logs/commands)
	logsPanel := m.panels["logs"]
	logsPanel.SetSize(m.width-2, bottomPanelHeight-6) // Adjust for command log section

	formattedLogs := lipgloss.JoinVertical(
		lipgloss.Left,
		commandLogFormatted,
		logsPanel.View(),
	)

	bottomPanel := panelStyle.
		Width(m.width).
		Height(bottomPanelHeight).
		Render(formattedLogs)

	// Bottom status bar with keybindings
	helpBar := statusStyle.Render("[q]Quit [?]Help [↑↓]Navigate panel content [alt+↑↓]Switch panels [↵]Select [i]Install [d]Delete")

	// Layout the UI
	topPanels := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		topPanels,
		bottomPanel,
		helpBar,
	)
}
