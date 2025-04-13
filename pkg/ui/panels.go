package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/VesperAkshay/lazynode/pkg/npm"
	"github.com/VesperAkshay/lazynode/pkg/scripts"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PanelStyle defines the style for a panel
var PanelStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#3c3836")).
	Padding(0, 0)

// TitleStyle defines the style for a panel title
var TitleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fabd2f")).
	Bold(true).
	PaddingLeft(1)

// HighlightStyle defines the style for highlighted items
var HighlightStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#b8bb26")).
	Bold(true)

// ErrorStyle defines the style for error messages
var ErrorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fb4934")).
	Bold(true)

// SelectedItemStyle defines the style for selected items in lists
var SelectedItemStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#83a598")).
	Background(lipgloss.Color("#3c3836")).
	Bold(true)

// HeaderStyle defines the style for panel headers
var HeaderStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#3c3836")).
	Foreground(lipgloss.Color("#ebdbb2")).
	Bold(true).
	Padding(0, 1)

// StatusStyle defines the style for status bars
var StatusStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("#282828")).
	Foreground(lipgloss.Color("#a89984")).
	Padding(0, 1)

// Panel represents a UI panel
type Panel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Panel, tea.Cmd)
	View() string
	Width() int
	Height() int
	SetSize(width, height int)
	Title() string
}

// ScriptsPanel displays and manages npm scripts
type ScriptsPanel struct {
	title        string
	width        int
	height       int
	scriptList   list.Model
	scriptRunner *scripts.ScriptRunner
	loading      bool
	error        string
	activeScript string
	logsPanel    *LogsPanel
}

// NewScriptsPanel creates a new scripts panel
func NewScriptsPanel(scriptRunner *scripts.ScriptRunner) *ScriptsPanel {
	// Create a list for the scripts
	scriptItems := []list.Item{}

	// Add scripts to the list
	for _, script := range scriptRunner.Scripts {
		scriptItems = append(scriptItems, scriptItem{script})
	}

	// Create compact custom list delegate
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.SetSpacing(0) // Reduce spacing between items
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(lipgloss.Color("#b8bb26")).
		Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(lipgloss.Color("#83a598")).
		Bold(false)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(lipgloss.Color("#ebdbb2"))

	// Create the list model
	scriptList := list.New(scriptItems, delegate, 0, 0)
	scriptList.Title = "🧪 Scripts"
	scriptList.SetShowStatusBar(false)    // Hide status bar to save space
	scriptList.SetFilteringEnabled(false) // Disable filtering to save space
	scriptList.SetShowHelp(false)         // Hide help to save space
	scriptList.Styles.Title = scriptList.Styles.Title.
		Foreground(lipgloss.Color("#b8bb26")).
		Background(lipgloss.Color("#3c3836")).
		Bold(true)

	return &ScriptsPanel{
		title:        "Scripts",
		scriptList:   scriptList,
		scriptRunner: scriptRunner,
	}
}

// scriptItem represents a script item in the list
type scriptItem struct {
	script scripts.Script
}

func (i scriptItem) Title() string       { return i.script.Name }
func (i scriptItem) Description() string { return i.script.Command }
func (i scriptItem) FilterValue() string { return i.script.Name }

// SetLogsPanel sets the logs panel for script output
func (p *ScriptsPanel) SetLogsPanel(logsPanel *LogsPanel) {
	p.logsPanel = logsPanel
}

// Init initializes the panel
func (p *ScriptsPanel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (p *ScriptsPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch msg.String() {
		case "enter":
			// Run the selected script
			if i, ok := p.scriptList.SelectedItem().(scriptItem); ok {
				p.loading = true
				p.activeScript = i.script.Name

				// Log the script execution
				if p.logsPanel != nil {
					p.logsPanel.AddLog(fmt.Sprintf("Running script: npm run %s", i.script.Name))
				}

				// Run the script in the background
				go func() {
					cmd, err := p.scriptRunner.RunScript(i.script.Name)
					if err != nil {
						p.error = fmt.Sprintf("Error running script: %v", err)
						if p.logsPanel != nil {
							p.logsPanel.AddLog(fmt.Sprintf("Error: %v", err))
						}
					} else {
						// Capture output from the command
						if p.logsPanel != nil {
							go func() {
								// Wait for the command to finish
								err := cmd.Wait()
								if err != nil {
									p.logsPanel.AddLog(fmt.Sprintf("Script exited with error: %v", err))
								} else {
									p.logsPanel.AddLog(fmt.Sprintf("Script completed: %s", i.script.Name))
								}
							}()
						}
					}
					p.loading = false
				}()
			}

		case "k", "j":
			// Navigate the list but also handle updating the script info
			if i, ok := p.scriptList.SelectedItem().(scriptItem); ok {
				p.activeScript = i.script.Name
			}
		}
	}

	// Update the list model
	p.scriptList, cmd = p.scriptList.Update(msg)

	return p, cmd
}

// View renders the panel
func (p *ScriptsPanel) View() string {
	// In a 4-panel grid, we need to be more economical with space
	availableHeight := p.height - 2 // Reserve 2 lines for status
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Update the list dimensions for compact display
	p.scriptList.SetSize(p.width, availableHeight)

	// Create a minimal compact view for small panels
	var statusInfo string
	if p.loading {
		statusInfo = HighlightStyle.Render("⟳ " + p.activeScript)
	} else if p.error != "" {
		statusInfo = ErrorStyle.Render(p.error)
	} else if _, ok := p.scriptList.SelectedItem().(scriptItem); ok {
		statusInfo = fmt.Sprintf("[↵]Run")
	}

	// Ultra compact view with minimal status line
	return fmt.Sprintf("%s\n%s",
		p.scriptList.View(),
		statusInfo)
}

// Width returns the panel width
func (p *ScriptsPanel) Width() int {
	return p.width
}

// Height returns the panel height
func (p *ScriptsPanel) Height() int {
	return p.height
}

// SetSize sets the panel size
func (p *ScriptsPanel) SetSize(width, height int) {
	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}

	p.width = width
	p.height = height

	// Adjust list size to account for borders and status line
	listWidth := width - 2
	listHeight := height - 3

	if listWidth < 5 {
		listWidth = 5
	}
	if listHeight < 2 {
		listHeight = 2
	}

	p.scriptList.SetSize(listWidth, listHeight)
}

// Title returns the panel title
func (p *ScriptsPanel) Title() string {
	return p.title
}

// PackageAction represents different package management actions
type PackageAction struct {
	Name        string
	Description string
	Key         string
	Command     string
}

// GetPackageActions returns the available package actions
func GetPackageActions() []PackageAction {
	return []PackageAction{
		{Name: "Install", Description: "Install a package", Key: "i", Command: "install"},
		{Name: "Install Dev", Description: "Install as dev dependency", Key: "I", Command: "install-dev"},
		{Name: "Uninstall", Description: "Uninstall a package", Key: "d", Command: "uninstall"},
		{Name: "Update", Description: "Update a package", Key: "u", Command: "update"},
		{Name: "Check Outdated", Description: "Check for outdated packages", Key: "o", Command: "outdated"},
	}
}

// PackagesPanel displays and manages npm packages
type PackagesPanel struct {
	title          string
	width          int
	height         int
	packageList    list.Model
	packageManager *npm.PackageManager
	loading        bool
	error          string
	input          textinput.Model
	showInput      bool
	inputMode      string // "install", "uninstall", "update", etc.
	statusMessage  string
	statusTime     time.Time
	spinner        int
	spinnerFrames  []string
	lastUpdate     time.Time
	showActions    bool       // New field to show action selection mode
	actionList     list.Model // New field for action list
	actions        []PackageAction
	showConfirm    bool          // Field for confirmation dialog
	confirmMessage string        // Message for confirmation dialog
	confirmAction  PackageAction // Action to perform if confirmed
	confirmPackage string        // Package to act on if confirmed
}

// NewPackagesPanel creates a new packages panel
func NewPackagesPanel(packageManager *npm.PackageManager) *PackagesPanel {
	// Create a delegate for custom list item rendering
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Copy().Foreground(lipgloss.Color("#b8bb26"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Copy().Foreground(lipgloss.Color("#a89984"))

	// Create the list
	packageList := list.New([]list.Item{}, delegate, 0, 0)
	packageList.Title = "Packages"
	packageList.SetShowStatusBar(false)
	packageList.SetFilteringEnabled(false) // Disable filtering for simplicity
	packageList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fabd2f")).
		Bold(true)

	// Create action delegate
	actionDelegate := list.NewDefaultDelegate()
	actionDelegate.Styles.SelectedTitle = actionDelegate.Styles.SelectedTitle.Copy().Foreground(lipgloss.Color("#b8bb26"))
	actionDelegate.Styles.SelectedDesc = actionDelegate.Styles.SelectedDesc.Copy().Foreground(lipgloss.Color("#a89984"))

	// Create the action list
	actions := GetPackageActions()
	actionItems := make([]list.Item, len(actions))
	for i, action := range actions {
		actionItems[i] = packageActionItem{action}
	}

	actionList := list.New(actionItems, actionDelegate, 0, 0)
	actionList.Title = "Package Actions"
	actionList.SetShowStatusBar(false)
	actionList.SetFilteringEnabled(false)
	actionList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#fabd2f")).
		Bold(true)

	// Create the input for package installation
	input := textinput.New()
	input.Placeholder = "Package name"
	input.Focus()

	// Spinner animation frames
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	// Create the panel
	panel := &PackagesPanel{
		title:          "Packages",
		packageList:    packageList,
		packageManager: packageManager,
		loading:        true, // Start with loading state
		error:          "",
		input:          input,
		showInput:      false,
		showActions:    false,
		spinnerFrames:  spinnerFrames,
		lastUpdate:     time.Now(),
		actionList:     actionList,
		actions:        actions,
	}

	// Immediately load packages when panel is created
	go func() {
		// Ensure packages are loaded initially
		if len(packageManager.Packages) == 0 {
			err := packageManager.LoadPackages()
			if err != nil {
				panel.error = fmt.Sprintf("Error loading packages: %v", err)
			}
		}

		// Add all packages to the list
		for _, pkg := range packageManager.Packages {
			panel.packageList.InsertItem(len(panel.packageList.Items()), packageItem{pkg})
		}

		// Set loading to false after initial load
		panel.loading = false
	}()

	return panel
}

// packageItem represents a package item in the list
type packageItem struct {
	pkg npm.Package
}

func (i packageItem) Title() string {
	title := i.pkg.Name
	if i.pkg.LatestVersion != "" {
		title += fmt.Sprintf(" (%s → %s)", i.pkg.Version, i.pkg.LatestVersion)
		return HighlightStyle.Render(title)
	} else {
		title += fmt.Sprintf(" (%s)", i.pkg.Version)
	}
	return title
}

func (i packageItem) Description() string {
	if i.pkg.Type == "devDependency" {
		return fmt.Sprintf("[dev] %s", i.pkg.Description)
	}
	return i.pkg.Description
}

func (i packageItem) FilterValue() string {
	return i.pkg.Name
}

// packageActionItem represents an action item in the action list
type packageActionItem struct {
	action PackageAction
}

func (i packageActionItem) Title() string {
	return fmt.Sprintf("[%s] %s", i.action.Key, i.action.Name)
}

func (i packageActionItem) Description() string {
	return i.action.Description
}

func (i packageActionItem) FilterValue() string {
	return i.action.Name
}

// executeAction executes a package management action
func (p *PackagesPanel) executeAction(action PackageAction, packageName string) {
	p.loading = true
	p.error = ""

	go func() {
		var err error
		defer func() {
			// Always reset UI state when done, whether successful or not
			p.loading = false
			p.showInput = false
			p.input.SetValue("")

			// Force UI update
			p.refreshPackageList()
		}()

		switch action.Command {
		case "install":
			err = p.packageManager.InstallPackage(packageName, false)
			if err == nil {
				p.statusMessage = fmt.Sprintf("✅ Installed %s", packageName)
				p.statusTime = time.Now()
			}
		case "install-dev":
			err = p.packageManager.InstallPackage(packageName, true)
			if err == nil {
				p.statusMessage = fmt.Sprintf("✅ Installed %s (dev)", packageName)
				p.statusTime = time.Now()
			}
		case "uninstall":
			err = p.packageManager.UninstallPackage(packageName)
			if err == nil {
				p.statusMessage = fmt.Sprintf("✅ Uninstalled %s", packageName)
				p.statusTime = time.Now()
			}
		case "update":
			err = p.packageManager.UpdatePackage(packageName)
			if err == nil {
				p.statusMessage = fmt.Sprintf("✅ Updated %s", packageName)
				p.statusTime = time.Now()
			}
		case "outdated":
			outdated, err := p.packageManager.CheckOutdatedPackages()
			if err != nil {
				p.error = fmt.Sprintf("Error checking outdated packages: %v", err)
			} else if len(outdated) > 0 {
				p.statusMessage = fmt.Sprintf("Found %d outdated packages", len(outdated))
				p.statusTime = time.Now()
			} else {
				p.statusMessage = "All packages are up to date"
				p.statusTime = time.Now()
			}
		}

		if err != nil {
			p.error = fmt.Sprintf("Error: %v", err)
		}
	}()
}

// Init initializes the panel
func (p *PackagesPanel) Init() tea.Cmd {
	// Load packages immediately in a non-blocking way
	go func() {
		// Force a reload of packages
		err := p.packageManager.LoadPackages()
		if err != nil {
			p.error = fmt.Sprintf("Error loading packages: %v", err)
		} else {
			// Load packages into the list
			p.refreshPackageList()
		}
		p.loading = false
	}()

	return nil
}

// Update handles messages
func (p *PackagesPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd

	// Update spinner animation
	if time.Since(p.lastUpdate) > 100*time.Millisecond {
		p.spinner = (p.spinner + 1) % len(p.spinnerFrames)
		p.lastUpdate = time.Now()
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input
		switch {
		case p.showConfirm:
			// Handle confirmation dialog
			switch msg.String() {
			case "y", "Y":
				// User confirmed the action
				p.showConfirm = false
				p.executeAction(p.confirmAction, p.confirmPackage)
				p.confirmPackage = ""
			case "n", "N", "esc":
				// User cancelled the action
				p.showConfirm = false
				p.confirmPackage = ""
			}

		case p.showInput:
			// Handle input mode
			switch msg.String() {
			case "enter":
				// Process the input
				value := p.input.Value()

				if value != "" {
					// Find the action by command
					for _, action := range p.actions {
						if action.Command == p.inputMode {
							// Special case for uninstall - show confirmation
							if action.Command == "uninstall" {
								p.showConfirm = true
								p.confirmMessage = fmt.Sprintf("Are you sure you want to uninstall '%s'?", value)
								p.confirmAction = action
								p.confirmPackage = value
								p.showInput = false
								return p, nil
							}

							p.executeAction(action, value)
							break
						}
					}
				}

				p.showInput = false

			case "esc":
				// Cancel the input
				p.showInput = false
				p.input.SetValue("")
			}

			// Update the input
			p.input, cmd = p.input.Update(msg)

		case p.showActions:
			// Handle action selection mode
			switch msg.String() {
			case "enter":
				// Get selected action
				if i, ok := p.actionList.SelectedItem().(packageActionItem); ok {
					action := i.action

					// If uninstall action, pre-fill with selected package
					if action.Command == "uninstall" && p.packageList.SelectedItem() != nil {
						if pkgItem, ok := p.packageList.SelectedItem().(packageItem); ok {
							// Show confirmation dialog directly
							p.showConfirm = true
							p.confirmMessage = fmt.Sprintf("Are you sure you want to uninstall '%s'?", pkgItem.pkg.Name)
							p.confirmAction = action
							p.confirmPackage = pkgItem.pkg.Name
						}
					} else if action.Command == "outdated" {
						// Directly run outdated check
						for _, a := range p.actions {
							if a.Command == "outdated" {
								p.executeAction(a, "")
								break
							}
						}
					} else {
						// Show input for other actions
						p.showInput = true
						p.inputMode = action.Command
						p.input.Placeholder = fmt.Sprintf("%s package", action.Name)
						p.input.Focus()
					}
				}

				// Hide action list
				p.showActions = false

			case "esc":
				// Cancel action selection
				p.showActions = false
			}

			// Update the action list
			p.actionList, cmd = p.actionList.Update(msg)

		case !p.showInput && !p.showActions:
			// Handle normal mode
			switch msg.String() {
			case "a":
				// Show action menu
				p.showActions = true
				p.actionList.Select(0)

			case "i":
				// Install a package
				p.showInput = true
				p.inputMode = "install"
				p.input.Placeholder = "Package name to install"
				p.input.Focus()

			case "shift+i":
				// Install a dev package
				p.showInput = true
				p.inputMode = "install-dev"
				p.input.Placeholder = "Package name to install as dev dependency"
				p.input.Focus()

			case "d":
				// Uninstall a package (with confirmation)
				if i, ok := p.packageList.SelectedItem().(packageItem); ok {
					p.showConfirm = true
					p.confirmMessage = fmt.Sprintf("Are you sure you want to uninstall '%s'?", i.pkg.Name)

					// Find the uninstall action
					for _, action := range p.actions {
						if action.Command == "uninstall" {
							p.confirmAction = action
							break
						}
					}

					p.confirmPackage = i.pkg.Name
				}

			case "o":
				// Check for outdated packages
				for _, action := range p.actions {
					if action.Command == "outdated" {
						p.executeAction(action, "")
						break
					}
				}

			case "u":
				// Update a package
				if i, ok := p.packageList.SelectedItem().(packageItem); ok && i.pkg.LatestVersion != "" {
					p.showInput = true
					p.inputMode = "update"
					p.input.Placeholder = "Confirm update (enter package name)"
					p.input.SetValue(i.pkg.Name)
					p.input.Focus()
				}

			case "/":
				// Search for a package
				p.showInput = true
				p.inputMode = "search"
				p.input.Placeholder = "Search for package"
				p.input.Focus()
			}
		}
	}

	// Update the list model
	p.packageList, cmd = p.packageList.Update(msg)

	return p, cmd
}

// refreshPackageList updates the package list with the latest data
func (p *PackagesPanel) refreshPackageList() {
	// Remember the currently selected index
	selectedIndex := p.packageList.Index()

	// Remember the currently selected package name (if any)
	var selectedPkg string
	if selected, ok := p.packageList.SelectedItem().(packageItem); ok {
		selectedPkg = selected.pkg.Name
	}

	// Clear the list
	p.packageList.SetItems([]list.Item{})

	// Add packages to the list with a micro-delay to allow UI updates
	for _, pkg := range p.packageManager.Packages {
		p.packageList.InsertItem(len(p.packageList.Items()), packageItem{pkg})
	}

	// Try to restore selection
	items := p.packageList.Items()
	if len(items) > 0 {
		// First try to find the same package by name
		if selectedPkg != "" {
			for i, item := range items {
				if pkgItem, ok := item.(packageItem); ok && pkgItem.pkg.Name == selectedPkg {
					p.packageList.Select(i)
					return
				}
			}
		}

		// Otherwise just restore the index if possible
		if selectedIndex < len(items) {
			p.packageList.Select(selectedIndex)
		} else {
			// Select the last item if the index is out of bounds
			p.packageList.Select(len(items) - 1)
		}
	}
}

// View renders the panel
func (p *PackagesPanel) View() string {
	// In a 4-panel grid, we need to be more economical with space
	availableHeight := p.height - 2 // Reserve 2 lines for status
	if availableHeight < 1 {
		availableHeight = 1
	}

	// Update the list dimensions for compact display
	p.packageList.SetSize(p.width, availableHeight)
	p.actionList.SetSize(p.width, availableHeight)

	// Spinner style for loading animation
	spinnerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")).Bold(true)

	// Show appropriate content based on panel state
	if p.loading {
		spinnerChar := spinnerStyle.Render(p.spinnerFrames[p.spinner])
		return fmt.Sprintf("%s\n%s",
			p.packageList.View(),
			spinnerChar+" Working...")
	}

	if p.error != "" {
		return fmt.Sprintf("%s\n%s",
			p.packageList.View(),
			ErrorStyle.Render(p.error))
	}

	// Show confirmation dialog more compactly
	if p.showConfirm {
		return fmt.Sprintf("%s\n%s",
			p.confirmMessage,
			"[y]Yes [n]No")
	}

	// Show action selection mode
	if p.showActions {
		return fmt.Sprintf("%s\n%s",
			p.actionList.View(),
			"[↵]Select [esc]Cancel")
	}

	// Show input mode
	if p.showInput {
		return fmt.Sprintf("%s\n%s",
			p.packageList.View(),
			p.input.View())
	}

	// Get the selected package details for status line
	var statusInfo string
	if i, ok := p.packageList.SelectedItem().(packageItem); ok {
		if i.pkg.Type == "devDependency" {
			statusInfo = "[i]Install [d]Del [dev]"
		} else {
			statusInfo = "[i]Install [d]Del"
		}
	}

	// Ultra compact view
	return fmt.Sprintf("%s\n%s",
		p.packageList.View(),
		statusInfo)
}

// Width returns the panel width
func (p *PackagesPanel) Width() int {
	return p.width
}

// Height returns the panel height
func (p *PackagesPanel) Height() int {
	return p.height
}

// SetSize sets the panel size
func (p *PackagesPanel) SetSize(width, height int) {
	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}

	p.width = width
	p.height = height

	// Adjust list sizes to account for borders and status line
	listWidth := width - 2
	listHeight := height - 3

	if listWidth < 5 {
		listWidth = 5
	}
	if listHeight < 2 {
		listHeight = 2
	}

	p.packageList.SetSize(listWidth, listHeight)
	p.actionList.SetSize(listWidth, height-4)
}

// Title returns the panel title
func (p *PackagesPanel) Title() string {
	return p.title
}

// LogsPanel displays command logs
type LogsPanel struct {
	title         string
	width         int
	height        int
	viewport      viewport.Model
	logs          []string
	spinner       int
	spinnerFrames []string
	lastUpdate    time.Time
	maxLogHistory int
}

// NewLogsPanel creates a new logs panel
func NewLogsPanel() *LogsPanel {
	viewport := viewport.New(0, 0)
	viewport.Style = lipgloss.NewStyle().
		PaddingLeft(1).
		PaddingRight(1)

	// Spinner animation frames
	spinnerFrames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	return &LogsPanel{
		title:         "Logs",
		viewport:      viewport,
		logs:          []string{},
		spinner:       0,
		spinnerFrames: spinnerFrames,
		lastUpdate:    time.Now(),
		maxLogHistory: 100, // Limit log history to prevent memory bloat
	}
}

// Init initializes the panel
func (p *LogsPanel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (p *LogsPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd

	// Update spinner animation
	if time.Since(p.lastUpdate) > 100*time.Millisecond {
		p.spinner = (p.spinner + 1) % len(p.spinnerFrames)
		p.lastUpdate = time.Now()
	}

	p.viewport, cmd = p.viewport.Update(msg)

	return p, cmd
}

// AddLog adds a log message
func (p *LogsPanel) AddLog(message string) {
	timestamp := time.Now().Format("15:04:05")
	logEntry := fmt.Sprintf("[%s] %s", timestamp, message)

	// Add the log to the top for newest-first ordering
	p.logs = append([]string{logEntry}, p.logs...)

	// Limit log history to prevent memory bloat
	if len(p.logs) > p.maxLogHistory {
		p.logs = p.logs[:p.maxLogHistory]
	}

	// Update the viewport content
	p.updateViewportContent()
}

// updateViewportContent updates the viewport content with the current logs
func (p *LogsPanel) updateViewportContent() {
	// Join all logs with newlines, most recent first
	content := strings.Join(p.logs, "\n")

	// Update the viewport content
	p.viewport.SetContent(content)

	// Auto-scroll to top (most recent logs)
	p.viewport.GotoTop()
}

// View renders the panel
func (p *LogsPanel) View() string {
	// Ensure we have reasonable dimensions
	if p.width < 10 {
		p.width = 10
	}
	if p.height < 3 {
		p.height = 3
	}

	// Make sure viewport size is properly set
	viewportWidth := p.width - 2   // Account for borders
	viewportHeight := p.height - 1 // Minimal height for logs

	if viewportWidth < 5 {
		viewportWidth = 5
	}

	if viewportHeight < 1 {
		viewportHeight = 1
	}

	p.viewport.Width = viewportWidth
	p.viewport.Height = viewportHeight

	// For compact view, maybe show only the most recent few logs
	// Create a condensed view with less lines
	maxLogs := 50
	if len(p.logs) > maxLogs {
		recentLogs := p.logs[:maxLogs]
		p.viewport.SetContent(strings.Join(recentLogs, "\n"))
	} else {
		p.viewport.SetContent(strings.Join(p.logs, "\n"))
	}

	// Auto-scroll to top (most recent logs)
	p.viewport.GotoTop()

	// Render the viewport without adding extra space
	return p.viewport.View()
}

// Width returns the panel width
func (p *LogsPanel) Width() int {
	return p.width
}

// Height returns the panel height
func (p *LogsPanel) Height() int {
	return p.height
}

// SetSize sets the panel size
func (p *LogsPanel) SetSize(width, height int) {
	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 2 {
		height = 2
	}

	p.width = width
	p.height = height

	// Update viewport dimensions with constraints
	viewportWidth := width - 2   // Account for borders
	viewportHeight := height - 1 // Minimal space for logs

	if viewportWidth < 5 {
		viewportWidth = 5
	}

	if viewportHeight < 1 {
		viewportHeight = 1
	}

	p.viewport.Width = viewportWidth
	p.viewport.Height = viewportHeight
}

// Title returns the panel title
func (p *LogsPanel) Title() string {
	return p.title
}
