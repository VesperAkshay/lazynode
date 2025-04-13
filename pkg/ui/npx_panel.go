package ui

import (
	"fmt"

	"github.com/VesperAkshay/lazynode/pkg/npx"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// NpxPanel displays and manages npx commands
type NpxPanel struct {
	title       string
	width       int
	height      int
	commandList list.Model
	npxRunner   *npx.Runner
	logsPanel   *LogsPanel
	loading     bool
	error       string
	input       textinput.Model
	showInput   bool
	history     []string
}

// npxItem represents an npx command item in the list
type npxItem struct {
	command npx.NpxCommand
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
		commandItems = append(commandItems, npxItem{cmd})
	}

	// Add recent commands
	for _, cmd := range npxRunner.GetRecentCommands() {
		commandItems = append(commandItems, npxItem{cmd})
	}

	// Create the list model
	commandList := list.New(commandItems, list.NewDefaultDelegate(), 0, 0)
	commandList.Title = "npx Commands"

	// Create the input model
	input := textinput.New()
	input.Placeholder = "npx command"
	input.Focus()

	return &NpxPanel{
		title:       "npx",
		commandList: commandList,
		npxRunner:   npxRunner,
		logsPanel:   logsPanel,
		input:       input,
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
		// Handle keyboard input
		switch {
		case p.showInput:
			// Handle input mode
			switch msg.String() {
			case "enter":
				// Process the input
				value := p.input.Value()

				if value != "" {
					p.loading = true

					// Run the command in the background
					go func() {
						p.logsPanel.AddLog(fmt.Sprintf("Running npx %s", value))

						c, err := p.npxRunner.RunCommand(value)
						if err != nil {
							p.error = fmt.Sprintf("Error running npx command: %v", err)
							p.logsPanel.AddLog(fmt.Sprintf("Error: %v", err))
						} else {
							// Wait for the command to finish
							err = c.Wait()
							if err != nil {
								p.logsPanel.AddLog(fmt.Sprintf("Command exited with error: %v", err))
							} else {
								p.logsPanel.AddLog(fmt.Sprintf("Command completed: npx %s", value))
							}
						}

						// Reset
						p.loading = false
						p.showInput = false
						p.input.SetValue("")
					}()
				}

				p.showInput = false

			case "esc":
				// Cancel the input
				p.showInput = false
				p.input.SetValue("")
			}

			// Update the input
			p.input, cmd = p.input.Update(msg)

		case !p.showInput:
			// Handle normal mode
			switch msg.String() {
			case "n":
				// New npx command
				p.showInput = true
				p.input.Focus()

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
			}
		}
	}

	// Update the list model
	p.commandList, cmd = p.commandList.Update(msg)

	return p, cmd
}

// View renders the panel
func (p *NpxPanel) View() string {
	// In a 4-panel grid, we need to be more economical with space
	if p.loading {
		return "Running npx command..."
	}

	if p.error != "" {
		return ErrorStyle.Render(p.error)
	}

	if p.showInput {
		return fmt.Sprintf("NPX Command:\n%s", p.input.View())
	}

	// Update the list dimensions
	availableHeight := p.height - 2 // Reserve space for title and help text
	if availableHeight < 1 {
		availableHeight = 1
	}
	p.commandList.SetSize(p.width, availableHeight)

	// Show a compact version of the command list
	return fmt.Sprintf("%s\n[n]New [↵]Run",
		p.commandList.View())
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
