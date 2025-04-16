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
	// Navigation Keys
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	PanelUp    key.Binding
	PanelDown  key.Binding
	PanelLeft  key.Binding
	PanelRight key.Binding
	NextTab    key.Binding
	PrevTab    key.Binding

	// Selection and Action Keys
	Enter  key.Binding
	Back   key.Binding
	Escape key.Binding

	// Tab Selection Keys
	TabScripts  key.Binding
	TabPackages key.Binding
	TabProject  key.Binding
	TabNpx      key.Binding
	TabLogs     key.Binding

	// Package Management
	Install      key.Binding
	InstallDev   key.Binding
	Delete       key.Binding
	Outdated     key.Binding
	Update       key.Binding
	Link         key.Binding
	UnLink       key.Binding
	Search       key.Binding
	InstallAll   key.Binding
	CheckMissing key.Binding

	// Script Management
	RunScript  key.Binding
	StopScript key.Binding

	// NPX Management
	NewNpx key.Binding
	RunNpx key.Binding

	// General Actions
	Reload  key.Binding
	Edit    key.Binding
	Open    key.Binding
	Build   key.Binding
	Test    key.Binding
	Publish key.Binding

	// UI Controls
	Help         key.Binding
	Quit         key.Binding
	ToggleDetail key.Binding
	FullScreen   key.Binding
	ActionMenu   key.Binding

	// Search Controls
	SearchUp     key.Binding
	SearchDown   key.Binding
	SearchInput  key.Binding
	SearchCancel key.Binding
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
		// Navigation
		{k.Up, k.Down, k.Left, k.Right, k.PanelUp, k.PanelDown, k.PanelLeft, k.PanelRight},
		// Tab selection
		{k.TabScripts, k.TabPackages, k.TabProject, k.TabNpx, k.TabLogs},
		// Package management
		{k.Install, k.InstallDev, k.Delete, k.Outdated, k.Update, k.Link, k.UnLink, k.Search, k.InstallAll, k.CheckMissing},
		// Script and NPX management
		{k.RunScript, k.StopScript, k.NewNpx, k.RunNpx},
		// General actions
		{k.Enter, k.Reload, k.Edit, k.Open, k.Build, k.Test, k.Publish},
		// UI controls
		{k.Help, k.Quit, k.ToggleDetail, k.FullScreen, k.ActionMenu},
	}
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation Keys
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
		PanelLeft: key.NewBinding(
			key.WithKeys("alt+left", "alt+h"),
			key.WithHelp("alt+←", "left panel"),
		),
		PanelRight: key.NewBinding(
			key.WithKeys("alt+right", "alt+l"),
			key.WithHelp("alt+→", "right panel"),
		),
		NextTab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next tab"),
		),
		PrevTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "previous tab"),
		),

		// Selection and Action Keys
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/run"),
		),
		Back: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "go back"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),

		// Tab Selection Keys
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

		// Package Management
		Install: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "install package"),
		),
		InstallDev: key.NewBinding(
			key.WithKeys("I"),
			key.WithHelp("I", "install dev package"),
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
		Link: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "link package"),
		),
		UnLink: key.NewBinding(
			key.WithKeys("L"),
			key.WithHelp("L", "unlink package"),
		),
		Search: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		InstallAll: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "install all dependencies"),
		),
		CheckMissing: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "check missing dependencies"),
		),

		// Script Management
		RunScript: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "run script"),
		),
		StopScript: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "stop script"),
		),

		// NPX Management
		NewNpx: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new npx command"),
		),
		RunNpx: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "run npx command"),
		),

		// General Actions
		Reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload UI"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit package.json"),
		),
		Open: key.NewBinding(
			key.WithKeys("o"),
			key.WithHelp("o", "open in editor"),
		),
		Build: key.NewBinding(
			key.WithKeys("b"),
			key.WithHelp("b", "build"),
		),
		Test: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "test"),
		),
		Publish: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "publish"),
		),

		// UI Controls
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		ToggleDetail: key.NewBinding(
			key.WithKeys("space"),
			key.WithHelp("space", "toggle details"),
		),
		FullScreen: key.NewBinding(
			key.WithKeys("F"),
			key.WithHelp("F", "toggle fullscreen"),
		),
		ActionMenu: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "open action menu"),
		),

		// Search Controls
		SearchUp: key.NewBinding(
			key.WithKeys("ctrl+p"),
			key.WithHelp("ctrl+p", "search up"),
		),
		SearchDown: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "search down"),
		),
		SearchInput: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		SearchCancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel search"),
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
	// Splash screen related fields
	showSplash   bool
	splashScreen SplashModel
	// Quit screen related fields
	showQuit   bool
	quitScreen QuitModel
	// Message queue for background operations
	msgChan chan tea.Msg
}

// NewModel initializes a new model with default values
func NewModel() Model {
	keys := DefaultKeyMap()
	helpPanel := NewHelpPanel()
	splashScreen := NewSplashModel()
	quitScreen := NewQuitModel()

	// Initialize logs panel early to avoid nil references
	logs := NewLogsPanel()

	return Model{
		keys:         keys,
		help:         help.New(),
		activeTab:    "scripts",
		panels:       make(map[string]Panel),
		logs:         logs,
		helpPanel:    helpPanel,
		showHelp:     false,
		ready:        false,
		showSplash:   true,
		splashScreen: splashScreen,
		showQuit:     false,
		quitScreen:   quitScreen,
		msgChan:      make(chan tea.Msg, 100), // Larger buffer for background messages
		error:        "",
		width:        80, // Default initial width
		height:       24, // Default initial height
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.splashScreen.Init(),
		m.detectProject,
		m.startTicker(),
		m.listenForBackgroundMsgs,
	)
}

// listenForBackgroundMsgs listens for messages from background goroutines
func (m Model) listenForBackgroundMsgs() tea.Msg {
	select {
	case msg := <-m.msgChan:
		return msg
	case <-time.After(100 * time.Millisecond):
		// If no message received after timeout, return a no-op message
		// and let the system continue processing
		return tea.Cmd(m.listenForBackgroundMsgs)
	}
}

// Load packages in the background to avoid blocking the UI
func (m *Model) loadPackagesAsync() tea.Cmd {
	return func() tea.Msg {
		// Load packages (optimized for speed)
		err := m.packageMgr.LoadPackages()
		if err != nil {
			m.logs.AddLog(fmt.Sprintf("Warning: Failed to load packages: %v", err))
		}

		// Create packages panel now that data is loaded
		packagesPanel := NewPackagesPanel(m.packageMgr)

		return packageLoadedMsg{
			panel: packagesPanel,
			count: len(m.packageMgr.Packages),
		}
	}
}

// detectProject is a command that tries to find a package.json file
func (m Model) detectProject() tea.Msg {
	// Display loading message immediately
	m.ready = false

	packageJSONPath, err := project.Detect()
	if err != nil {
		return errorMsg(fmt.Sprintf("Could not find a package.json file: %v", err))
	}

	// Initialize the project (fast operation)
	proj, err := project.NewProject(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error loading project: %v", err))
	}

	// Initialize package manager (fast initial setup)
	pkgMgr, err := npm.NewPackageManager(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error initializing package manager: %v", err))
	}

	// Initialize script runner (fast operation)
	scriptRunner, err := scripts.NewScriptRunner(packageJSONPath)
	if err != nil {
		return errorMsg(fmt.Sprintf("Error initializing script runner: %v", err))
	}

	// Initialize npx runner (fast operation)
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

	// Handle splash screen if it's active
	if m.showSplash {
		// Update splash screen
		updatedSplash, cmd := m.splashScreen.Update(msg)
		m.splashScreen = updatedSplash

		// If splash screen is done, transition to main UI
		if m.splashScreen.IsDone() {
			m.showSplash = false
		} else if cmd != nil {
			return m, cmd
		}
	}

	// Handle quit screen if it's active
	if m.showQuit {
		// Update quit screen
		updatedQuit, cmd := m.quitScreen.Update(msg)
		m.quitScreen = updatedQuit

		// If quit screen is done, exit the application
		if m.quitScreen.IsDone() {
			return m, tea.Quit
		} else if cmd != nil {
			return m, cmd
		}

		// Only show quit screen, nothing else to update
		return m, nil
	}

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

			// Improved panel layout - 3-column grid with logs at bottom
			// Left side: Scripts (top) and Project (bottom)
			// Middle: Packages
			// Right: NPX
			leftColumnWidth := termWidth * 25 / 100                                 // 25% of width
			rightColumnWidth := termWidth * 25 / 100                                // 25% of width
			middleColumnWidth := termWidth - leftColumnWidth - rightColumnWidth - 2 // Remaining space (minus gaps)

			// Make sure no column is too narrow
			minColumnWidth := 25
			if leftColumnWidth < minColumnWidth {
				leftColumnWidth = minColumnWidth
			}
			if rightColumnWidth < minColumnWidth {
				rightColumnWidth = minColumnWidth
			}
			if middleColumnWidth < minColumnWidth {
				middleColumnWidth = minColumnWidth
			}

			// Height calculation - top panels take 3/4 of available space
			topRowHeight := availableHeight * 3 / 4
			if topRowHeight < 8 {
				topRowHeight = 8
			}

			bottomRowHeight := availableHeight - topRowHeight - 1 // -1 for gap
			if bottomRowHeight < 3 {
				bottomRowHeight = 3
			}

			// Individual panel heights
			scriptsPanelHeight := topRowHeight / 2
			projectPanelHeight := topRowHeight - scriptsPanelHeight - 1 // -1 for gap

			// Update all panel sizes based on their position in the layout
			for name, panel := range m.panels {
				switch name {
				case "logs":
					panel.SetSize(termWidth-2, bottomRowHeight-2)
				case "scripts":
					panel.SetSize(leftColumnWidth-2, scriptsPanelHeight-2)
				case "project":
					panel.SetSize(leftColumnWidth-2, projectPanelHeight-2)
				case "packages":
					panel.SetSize(middleColumnWidth-2, topRowHeight-2)
				case "npx":
					panel.SetSize(rightColumnWidth-2, topRowHeight-2)
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
			// Show quit screen instead of immediately quitting
			m.showQuit = true
			m.quitScreen = NewQuitModel()
			return m, m.quitScreen.Init()

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
				switch m.activeTab {
				case "project":
					m.activeTab = "scripts"
				case "packages":
					m.activeTab = "scripts"
				case "npx":
					m.activeTab = "scripts"
				case "logs":
					if m.panels["npx"] != nil {
						m.activeTab = "npx"
					} else if m.panels["packages"] != nil {
						m.activeTab = "packages"
					} else {
						m.activeTab = "project"
					}
				}
				return m, nil

			case key.Matches(msg, m.keys.PanelDown):
				// Navigate panels down
				switch m.activeTab {
				case "scripts":
					m.activeTab = "project"
				case "project":
					m.activeTab = "logs"
				case "packages":
					m.activeTab = "logs"
				case "npx":
					m.activeTab = "logs"
				}
				return m, nil

			case key.Matches(msg, m.keys.PanelLeft):
				// Navigate panels left
				switch m.activeTab {
				case "packages":
					if m.panels["project"] != nil {
						m.activeTab = "project"
					} else {
						m.activeTab = "scripts"
					}
				case "npx":
					m.activeTab = "packages"
				case "logs":
					m.activeTab = "project"
				}
				return m, nil

			case key.Matches(msg, m.keys.PanelRight):
				// Navigate panels right
				switch m.activeTab {
				case "scripts", "project":
					m.activeTab = "packages"
				case "packages":
					if m.panels["npx"] != nil {
						m.activeTab = "npx"
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

		// Create scripts panel (fast operation)
		scriptsPanel := NewScriptsPanel(m.scriptRunner)
		scriptsPanel.SetLogsPanel(m.logs)
		m.panels["scripts"] = scriptsPanel

		// Create project panel (fast operation)
		m.panels["project"] = NewProjectPanel(m.project)

		// Create NPX panel (fast operation)
		m.panels["npx"] = NewNpxPanel(m.npxRunner, m.logs)

		// Create logs panel (fast operation)
		m.panels["logs"] = m.logs

		// Make UI ready immediately to show the interface
		m.ready = true

		// Add welcome message to logs
		m.logs.AddLog(fmt.Sprintf("Welcome to LazyNode v1.0 - Managing project: %s", m.project.Name))
		m.logs.AddLog("Loading packages...")

		// Initialize panel sizes
		m.updatePanelSizes()

		// Load packages asynchronously and restart background message listener
		return m, tea.Batch(m.loadPackagesAsync(), m.listenForBackgroundMsgs)

	case packageLoadedMsg:
		// Update the UI with the loaded package panel
		m.panels["packages"] = msg.panel
		m.logs.AddLog(fmt.Sprintf("Found %d packages in package.json", msg.count))
		m.logs.AddLog("Use tabs 1-5 to navigate between panels")
		m.logs.AddLog("Press ? for help")

		// Make sure panel sizes are updated
		m.updatePanelSizes()

		// Restart background message listener
		return m, m.listenForBackgroundMsgs

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
	// Show splash screen if active
	if m.showSplash {
		return m.splashScreen.View()
	}

	// Show quit screen if active
	if m.showQuit {
		return m.quitScreen.View()
	}

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

	// Improved panel layout - 3-column grid with logs at bottom
	// Left side: Scripts (top) and Project (bottom)
	// Middle: Packages
	// Right: NPX
	leftColumnWidth := termWidth * 25 / 100                                 // 25% of width
	rightColumnWidth := termWidth * 25 / 100                                // 25% of width
	middleColumnWidth := termWidth - leftColumnWidth - rightColumnWidth - 2 // Remaining space (minus gaps)

	// Make sure no column is too narrow
	minColumnWidth := 25
	if leftColumnWidth < minColumnWidth {
		leftColumnWidth = minColumnWidth
	}
	if rightColumnWidth < minColumnWidth {
		rightColumnWidth = minColumnWidth
	}
	if middleColumnWidth < minColumnWidth {
		middleColumnWidth = minColumnWidth
	}

	// Height calculation - top panels take 3/4 of available space
	topRowHeight := availableHeight * 3 / 4
	if topRowHeight < 8 {
		topRowHeight = 8
	}

	bottomRowHeight := availableHeight - topRowHeight - 1 // -1 for gap
	if bottomRowHeight < 3 {
		bottomRowHeight = 3
	}

	// Individual panel heights
	scriptsPanelHeight := topRowHeight / 2
	projectPanelHeight := topRowHeight - scriptsPanelHeight - 1 // -1 for gap

	// Top status bar
	statusBar := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ebdbb2")).
		Background(lipgloss.Color("#3c3836")).
		Padding(0, 1).
		Width(termWidth).
		Render(fmt.Sprintf("LazyNode - %s", m.project.Name))

	// Panel style with rounded borders for a modern look
	panelStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4d4d4d")).
		Padding(0, 0).
		Margin(0, 0)

	// Style for active panel with distinct color
	selectedPanelStyle := panelStyle.Copy().
		BorderForeground(lipgloss.Color("#b8bb26"))

	// Style for panel titles with a pop of color
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fabd2f")).
		Padding(0, 1)

	// Loading style for panels still loading
	loadingStyle := panelStyle.Copy().
		BorderForeground(lipgloss.Color("#83a598"))

	// Prepare all panels
	panelContents := make(map[string]string)
	panelTitles := make(map[string]string)

	// Set panel sizes and get content
	for name, panel := range m.panels {
		switch name {
		case "logs":
			panel.SetSize(termWidth-2, bottomRowHeight-2)
		case "scripts":
			panel.SetSize(leftColumnWidth-2, scriptsPanelHeight-2)
		case "project":
			panel.SetSize(leftColumnWidth-2, projectPanelHeight-2)
		case "packages":
			panel.SetSize(middleColumnWidth-2, topRowHeight-2)
		case "npx":
			panel.SetSize(rightColumnWidth-2, topRowHeight-2)
		}

		// Get panel content and title
		panelContents[name] = panel.View()
		panelTitles[name] = panel.Title()
	}

	// Render panel with proper styling and highlighting active panel
	renderPanel := func(name string, width, height int) string {
		if _, exists := m.panels[name]; exists {
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
		} else if name == "packages" {
			// Special case for packages panel if it's not loaded yet
			content := fmt.Sprintf("%s\n%s",
				titleStyle.Render("Packages"),
				"Loading packages...")

			return loadingStyle.
				Width(width).
				Height(height).
				Render(content)
		}

		// Fallback for other missing panels
		return panelStyle.
			Width(width).
			Height(height).
			Render("Loading...")
	}

	// Render all panels with their new dimensions
	scriptsRendered := renderPanel("scripts", leftColumnWidth, scriptsPanelHeight)
	projectRendered := renderPanel("project", leftColumnWidth, projectPanelHeight)
	packagesRendered := renderPanel("packages", middleColumnWidth, topRowHeight)
	npxRendered := renderPanel("npx", rightColumnWidth, topRowHeight)
	logsRendered := renderPanel("logs", termWidth, bottomRowHeight)

	// Layout the panels - left column stacked vertically
	leftColumn := lipgloss.JoinVertical(lipgloss.Left, scriptsRendered, projectRendered)

	// Arrange the top row with the 3 columns
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, leftColumn, packagesRendered, npxRendered)

	// Enhanced help bar with more details on available shortcuts
	helpText := "[q]Quit [?]Help [Tab]Switch panels [1-5]Select panel [↵]Select"
	if m.activeTab == "scripts" {
		helpText += " | [r]Run script"
	} else if m.activeTab == "packages" {
		helpText += " | [i]Install [d]Delete [u]Update [b]Build [t]Test [p]Publish [e]Edit"
	}

	helpBar := lipgloss.NewStyle().
		Background(lipgloss.Color("#3c3836")).
		Foreground(lipgloss.Color("#ebdbb2")).
		Padding(0, 1).
		Width(termWidth).
		Render(helpText)

	// Complete layout
	ui := lipgloss.JoinVertical(
		lipgloss.Left,
		statusBar,
		topRow,
		logsRendered,
		helpBar,
	)

	// Final trim to ensure we don't exceed terminal dimensions
	return lipgloss.NewStyle().
		MaxWidth(termWidth).
		MaxHeight(termHeight).
		Render(ui)
}

// packageLoadedMsg is sent when packages have been loaded
type packageLoadedMsg struct {
	panel *PackagesPanel
	count int
}

// updatePanelSizes updates the sizes of all panels
func (m *Model) updatePanelSizes() {
	if !m.ready {
		return
	}

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

	// Improved panel layout calculations
	leftColumnWidth := termWidth * 25 / 100
	rightColumnWidth := termWidth * 25 / 100
	middleColumnWidth := termWidth - leftColumnWidth - rightColumnWidth - 2

	// Make sure no column is too narrow
	minColumnWidth := 25
	if leftColumnWidth < minColumnWidth {
		leftColumnWidth = minColumnWidth
	}
	if rightColumnWidth < minColumnWidth {
		rightColumnWidth = minColumnWidth
	}
	if middleColumnWidth < minColumnWidth {
		middleColumnWidth = minColumnWidth
	}

	// Height calculation
	topRowHeight := availableHeight * 3 / 4
	if topRowHeight < 8 {
		topRowHeight = 8
	}

	bottomRowHeight := availableHeight - topRowHeight - 1
	if bottomRowHeight < 3 {
		bottomRowHeight = 3
	}

	// Individual panel heights
	scriptsPanelHeight := topRowHeight / 2
	projectPanelHeight := topRowHeight - scriptsPanelHeight - 1

	// Update all panel sizes based on their position in the layout
	for name, panel := range m.panels {
		switch name {
		case "logs":
			panel.SetSize(termWidth-2, bottomRowHeight-2)
		case "scripts":
			panel.SetSize(leftColumnWidth-2, scriptsPanelHeight-2)
		case "project":
			panel.SetSize(leftColumnWidth-2, projectPanelHeight-2)
		case "packages":
			panel.SetSize(middleColumnWidth-2, topRowHeight-2)
		case "npx":
			panel.SetSize(rightColumnWidth-2, topRowHeight-2)
		}
	}

	// Update help panel size
	if m.helpPanel != nil {
		m.helpPanel.SetSize(termWidth, termHeight)
	}

	// Update help model width
	m.help.Width = termWidth
}
