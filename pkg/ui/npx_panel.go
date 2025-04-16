package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/VesperAkshay/lazynode/pkg/npx"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NpxPanel displays and manages npx commands
type NpxPanel struct {
	title          string
	width          int
	height         int
	commandList    list.Model
	npxRunner      *npx.Runner
	logsPanel      *LogsPanel
	loading        bool
	error          string
	input          textinput.Model
	descInput      textinput.Model
	mode           string // "view", "input", "edit", "category"
	showInput      bool
	editingCommand npx.NpxCommand
	cursor         int
	categoryIndex  int
	categories     []string
	categoryMap    map[string][]list.Item
}

// npxItem represents an npx command item in the list
type npxItem struct {
	command npx.NpxCommand
}

// newNpxItem creates a new npx item
func newNpxItem(cmd npx.NpxCommand) npxItem {
	return npxItem{command: cmd}
}

func (i npxItem) Title() string       { return i.command.Name }
func (i npxItem) Description() string { return i.command.Description }
func (i npxItem) FilterValue() string { return i.command.Name }

// NewNpxPanel creates a new npx panel
func NewNpxPanel(npxRunner *npx.Runner, logsPanel *LogsPanel) *NpxPanel {
	// Create a list for the commands
	var commandItems []list.Item

	// Add popular commands
	for _, cmd := range npx.GetPopularCommands() {
		commandItems = append(commandItems, newNpxItem(cmd))
	}

	// Add recent commands
	for _, cmd := range npxRunner.GetRecentCommands() {
		commandItems = append(commandItems, newNpxItem(cmd))
	}

	// Create the list model
	commandList := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commandList.Title = "npx Commands"
	commandList.SetShowStatusBar(false)
	commandList.SetFilteringEnabled(true)
	commandList.SetShowHelp(false)

	// Create the command input model
	input := textinput.New()
	input.Placeholder = "Enter npx command"
	input.Focus()

	// Create the description input model
	descInput := textinput.New()
	descInput.Placeholder = "Enter command description (optional)"

	// Define categories
	categories := []string{"All Commands", "Recent Commands", "Popular Tools", "Development", "Testing", "Utilities"}

	// Create category map
	categoryMap := make(map[string][]list.Item)
	categoryMap["All Commands"] = commandItems

	// Separate recent commands from npxRunner
	var recentItems []list.Item
	for _, cmd := range npxRunner.GetRecentCommands() {
		recentItems = append(recentItems, newNpxItem(cmd))
	}
	categoryMap["Recent Commands"] = recentItems

	// Separate popular commands
	var popularItems []list.Item
	var devItems []list.Item
	var testingItems []list.Item
	var utilityItems []list.Item

	for _, cmd := range npx.GetPopularCommands() {
		popularItems = append(popularItems, newNpxItem(cmd))

		// Categorize by type
		switch {
		case strings.Contains(cmd.Name, "test") || strings.Contains(cmd.Name, "jest") || strings.Contains(cmd.Name, "mocha"):
			testingItems = append(testingItems, newNpxItem(cmd))
		case strings.Contains(cmd.Name, "create-") || strings.Contains(cmd.Name, "init") || strings.Contains(cmd.Name, "vite") || strings.Contains(cmd.Name, "next"):
			devItems = append(devItems, newNpxItem(cmd))
		default:
			utilityItems = append(utilityItems, newNpxItem(cmd))
		}
	}

	categoryMap["Popular Tools"] = popularItems
	categoryMap["Development"] = devItems
	categoryMap["Testing"] = testingItems
	categoryMap["Utilities"] = utilityItems

	return &NpxPanel{
		title:         "npx",
		commandList:   commandList,
		npxRunner:     npxRunner,
		logsPanel:     logsPanel,
		input:         input,
		descInput:     descInput,
		mode:          "view",
		categories:    categories,
		categoryMap:   categoryMap,
		categoryIndex: 0,
	}
}

// Init initializes the panel
func (p *NpxPanel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (p *NpxPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input based on mode
		switch p.mode {
		case "view":
			// Handle view mode keys
			switch msg.String() {
			case "n":
				// New npx command
				p.mode = "input"
				p.input.Focus()
				p.input.SetValue("")
				p.descInput.SetValue("")
				return p, nil

			case "e":
				// Edit selected command
				if i, ok := p.commandList.SelectedItem().(npxItem); ok {
					p.mode = "edit"
					p.editingCommand = i.command
					p.input.SetValue(i.command.Command)
					p.descInput.SetValue(i.command.Description)
					p.input.Focus()
					return p, nil
				}

			case "d":
				// Delete selected command
				if i, ok := p.commandList.SelectedItem().(npxItem); ok {
					// Only allow deletion of custom commands, not built-in ones
					isBuiltIn := false
					for _, cmd := range npx.GetPopularCommands() {
						if cmd.Command == i.command.Command {
							isBuiltIn = true
							break
						}
					}

					if !isBuiltIn {
						// Remove from runner's recent commands
						// Find and remove from recent commands
						newRecentCommands := []npx.NpxCommand{}
						for _, cmd := range p.npxRunner.GetRecentCommands() {
							if cmd.Command != i.command.Command {
								newRecentCommands = append(newRecentCommands, cmd)
							}
						}

						// Save the updated recent commands by clearing and re-adding
						p.npxRunner.RecentCommands = newRecentCommands
						p.npxRunner.SaveCache()

						// Rebuild the command list
						var commandItems []list.Item
						for _, cmd := range p.npxRunner.GetRecentCommands() {
							commandItems = append(commandItems, newNpxItem(cmd))
						}

						// Add popular commands
						for _, cmd := range npx.GetPopularCommands() {
							commandItems = append(commandItems, newNpxItem(cmd))
						}

						// Update the list
						p.commandList.SetItems(commandItems)

						// Update category map
						p.updateCategoryMap()
					}
				}

			case "c":
				// Switch to category selection mode
				p.mode = "category"
				return p, nil

			case "enter":
				// Run the selected command
				if i, ok := p.commandList.SelectedItem().(npxItem); ok {
					p.loading = true

					// Run the command in the background
					go func() {
						p.logsPanel.AddLog(fmt.Sprintf("Running npx %s", i.command.Command))

						c, err := p.npxRunner.RunCommand(i.command.Command)
						if err != nil {
							p.error = fmt.Sprintf("Error running npx command: %v", err)
							p.logsPanel.AddLog(fmt.Sprintf("Error: %v", err))
						} else {
							// Wait for the command to finish
							err = c.Wait()
							if err != nil {
								p.logsPanel.AddLog(fmt.Sprintf("Command exited with error: %v", err))
							} else {
								p.logsPanel.AddLog(fmt.Sprintf("Command completed: npx %s", i.command.Command))
							}
						}

						p.loading = false
					}()
				}

			case "/":
				// Focus search
				// Set focus to the built-in filter that comes with list.Model
				p.commandList.SetFilteringEnabled(true) // Ensure filtering is enabled
				// Send a '/' character to start filtering
				filterKeyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
				p.commandList, cmd = p.commandList.Update(filterKeyMsg)
				return p, cmd
			}

		case "category":
			// Handle category selection mode
			switch msg.String() {
			case "up", "k":
				if p.cursor > 0 {
					p.cursor--
				}
			case "down", "j":
				if p.cursor < len(p.categories)-1 {
					p.cursor++
				}
			case "enter":
				// Select category
				p.categoryIndex = p.cursor
				p.commandList.SetItems(p.categoryMap[p.categories[p.categoryIndex]])
				p.mode = "view"
			case "esc":
				// Cancel category selection
				p.mode = "view"
			}
			return p, nil

		case "input", "edit":
			// Handle input/edit mode keys
			switch msg.String() {
			case "tab":
				// Switch between command and description inputs
				if p.input.Focused() {
					p.input.Blur()
					p.descInput.Focus()
				} else {
					p.descInput.Blur()
					p.input.Focus()
				}

			case "enter":
				if p.descInput.Focused() {
					// When pressing enter in the description field, treat it as submission
					commandValue := p.input.Value()
					descValue := p.descInput.Value()

					if commandValue != "" {
						p.loading = true

						// Cache the command with description
						p.npxRunner.CacheCommand(commandValue, descValue)

						if p.mode == "input" {
							// Run the command in the background
							go func() {
								p.logsPanel.AddLog(fmt.Sprintf("Running npx %s", commandValue))

								c, err := p.npxRunner.RunCommand(commandValue)
								if err != nil {
									p.error = fmt.Sprintf("Error running npx command: %v", err)
									p.logsPanel.AddLog(fmt.Sprintf("Error: %v", err))
								} else {
									// Wait for the command to finish
									err = c.Wait()
									if err != nil {
										p.logsPanel.AddLog(fmt.Sprintf("Command exited with error: %v", err))
									} else {
										p.logsPanel.AddLog(fmt.Sprintf("Command completed: npx %s", commandValue))
									}
								}

								// Reset
								p.loading = false
								p.mode = "view"
								p.input.SetValue("")
								p.descInput.SetValue("")

								// Refresh the command list
								p.refreshCommandList()
							}()
						} else {
							// Just update the command in edit mode
							go func() {
								// If we were editing, update the description
								if p.mode == "edit" {
									for i, item := range p.commandList.Items() {
										if npxItem, ok := item.(npxItem); ok {
											if npxItem.command.Command == p.editingCommand.Command {
												// Update the item
												items := p.commandList.Items()
												newCmd := npx.NpxCommand{
													Name:        filepath.Base(commandValue),
													Description: descValue,
													Command:     commandValue,
												}
												// Create a new item using constructor
												items[i] = newNpxItem(newCmd)
												p.commandList.SetItems(items)
												break
											}
										}
									}
								}

								p.loading = false
								p.mode = "view"
								p.input.SetValue("")
								p.descInput.SetValue("")

								// Refresh the command list
								p.refreshCommandList()
							}()
						}
					}
				} else if p.input.Focused() {
					// When pressing enter in the command field, move to description
					p.input.Blur()
					p.descInput.Focus()
				}

			case "esc":
				// Cancel the input/edit
				p.mode = "view"
				p.input.SetValue("")
				p.descInput.SetValue("")
			}

			// Update the inputs
			if p.input.Focused() {
				p.input, cmd = p.input.Update(msg)
			} else if p.descInput.Focused() {
				p.descInput, cmd = p.descInput.Update(msg)
			}

			if cmd != nil {
				return p, cmd
			}
		}
	}

	// Only update the list if in view mode
	if p.mode == "view" {
		p.commandList, cmd = p.commandList.Update(msg)
	}

	return p, cmd
}

// View renders the panel
func (p *NpxPanel) View() string {
	// Styles for better UI
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fabd2f"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#b8bb26"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Background(lipgloss.Color("#3c3836"))

	// Normal text style
	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ebdbb2"))

	// Loading view
	if p.loading {
		return titleStyle.Render("Running npx command...") + "\n\n" +
			textStyle.Render("Please wait while the command executes.\nCheck the logs panel for output.")
	}

	// Error view
	if p.error != "" {
		return ErrorStyle.Render(p.error)
	}

	switch p.mode {
	case "input":
		// Command input view
		return fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
			titleStyle.Render("New NPX Command"),
			labelStyle.Render("Command:"),
			p.input.View(),
			labelStyle.Render("Description:"),
			p.descInput.View(),
			textStyle.Render("[Tab]Switch fields [Enter]Submit [Esc]Cancel"))

	case "edit":
		// Command edit view
		return fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n%s\n\n%s",
			titleStyle.Render("Edit NPX Command"),
			labelStyle.Render("Command:"),
			p.input.View(),
			labelStyle.Render("Description:"),
			p.descInput.View(),
			textStyle.Render("[Tab]Switch fields [Enter]Submit [Esc]Cancel"))

	case "category":
		// Category selection view
		var sb strings.Builder
		sb.WriteString(titleStyle.Render("Select Category:") + "\n\n")

		for i, category := range p.categories {
			itemCount := len(p.categoryMap[category])
			itemText := fmt.Sprintf("%s (%d)", category, itemCount)

			if i == p.cursor {
				sb.WriteString(selectedStyle.Render(" • "+itemText) + "\n")
			} else {
				sb.WriteString(textStyle.Render(" • "+itemText) + "\n")
			}
		}

		sb.WriteString("\n" + textStyle.Render("[↵]Select [Esc]Cancel"))
		return sb.String()
	}

	// Update the list title to show the current category
	p.commandList.Title = p.categories[p.categoryIndex]

	// Update the list dimensions
	availableHeight := p.height - 2 // Reserve space for title and help text
	if availableHeight < 1 {
		availableHeight = 1
	}
	p.commandList.SetSize(p.width, availableHeight)

	// View mode with enhanced help text
	return fmt.Sprintf("%s\n%s",
		p.commandList.View(),
		textStyle.Render("[n]New [e]Edit [d]Delete [c]Categories [/]Search [↵]Run"))
}

// refreshCommandList refreshes the command list with the latest commands
func (p *NpxPanel) refreshCommandList() {
	var commandItems []list.Item

	// Add recent commands from runner
	for _, cmd := range p.npxRunner.GetRecentCommands() {
		commandItems = append(commandItems, newNpxItem(cmd))
	}

	// Add popular commands
	for _, cmd := range npx.GetPopularCommands() {
		commandItems = append(commandItems, newNpxItem(cmd))
	}

	// Update the list
	p.commandList.SetItems(commandItems)

	// Update category map
	p.updateCategoryMap()
}

// updateCategoryMap updates the category map
func (p *NpxPanel) updateCategoryMap() {
	// Clear existing categories
	for k := range p.categoryMap {
		p.categoryMap[k] = []list.Item{}
	}

	// Fill "All Commands" with all items
	p.categoryMap["All Commands"] = p.commandList.Items()

	// Separate recent commands
	for _, cmd := range p.npxRunner.GetRecentCommands() {
		p.categoryMap["Recent Commands"] = append(p.categoryMap["Recent Commands"], newNpxItem(cmd))
	}

	// Process other categories
	for _, item := range p.commandList.Items() {
		if npxItem, ok := item.(npxItem); ok {
			cmd := npxItem.command

			// Categorize by type
			switch {
			case strings.Contains(cmd.Name, "test") || strings.Contains(cmd.Name, "jest") || strings.Contains(cmd.Name, "mocha"):
				p.categoryMap["Testing"] = append(p.categoryMap["Testing"], newNpxItem(cmd))
			case strings.Contains(cmd.Name, "create-") || strings.Contains(cmd.Name, "init") || strings.Contains(cmd.Name, "vite") || strings.Contains(cmd.Name, "next"):
				p.categoryMap["Development"] = append(p.categoryMap["Development"], newNpxItem(cmd))
			}

			// Add popular tools
			for _, popularCmd := range npx.GetPopularCommands() {
				if popularCmd.Command == cmd.Command {
					p.categoryMap["Popular Tools"] = append(p.categoryMap["Popular Tools"], newNpxItem(cmd))
					break
				}
			}

			// Add to utilities if not categorized elsewhere
			if !strings.Contains(cmd.Name, "test") &&
				!strings.Contains(cmd.Name, "jest") &&
				!strings.Contains(cmd.Name, "mocha") &&
				!strings.Contains(cmd.Name, "create-") &&
				!strings.Contains(cmd.Name, "init") &&
				!strings.Contains(cmd.Name, "vite") &&
				!strings.Contains(cmd.Name, "next") {

				isPopular := false
				for _, popularCmd := range npx.GetPopularCommands() {
					if popularCmd.Command == cmd.Command {
						isPopular = true
						break
					}
				}

				if !isPopular {
					p.categoryMap["Utilities"] = append(p.categoryMap["Utilities"], newNpxItem(cmd))
				}
			}
		}
	}
}

// Width returns the panel width
func (p *NpxPanel) Width() int {
	return p.width
}

// Height returns the panel height
func (p *NpxPanel) Height() int {
	return p.height
}

// SetSize sets the panel size
func (p *NpxPanel) SetSize(width, height int) {
	// Ensure minimum dimensions
	if width < 10 {
		width = 10
	}
	if height < 5 {
		height = 5
	}

	p.width = width
	p.height = height

	// Set list size with proper constraints
	listWidth := width - 2   // Account for borders
	listHeight := height - 2 // Account for title and help text

	if listWidth < 5 {
		listWidth = 5
	}
	if listHeight < 1 {
		listHeight = 1
	}

	p.commandList.SetSize(listWidth, listHeight)
}

// Title returns the panel title
func (p *NpxPanel) Title() string {
	return p.title
}
