package ui

import (
	"fmt"

	"github.com/VesperAkshay/lazynode/pkg/project"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// ProjectPanel displays and manages package.json
type ProjectPanel struct {
	title     string
	width     int
	height    int
	project   *project.Project
	mode      string // "view", "edit"
	editKey   string
	editValue string
	input     textinput.Model
	error     string
	loading   bool
}

// NewProjectPanel creates a new project panel
func NewProjectPanel(project *project.Project) *ProjectPanel {
	// Create the input model
	input := textinput.New()
	input.Focus()

	return &ProjectPanel{
		title:   "Project",
		project: project,
		mode:    "view",
		input:   input,
	}
}

// Init initializes the panel
func (p *ProjectPanel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (p *ProjectPanel) Update(msg tea.Msg) (Panel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch p.mode {
		case "view":
			// Handle view mode keys
			switch msg.String() {
			case "e":
				// Edit name
				p.mode = "edit"
				p.editKey = "name"
				p.editValue = p.project.Name
				p.input.SetValue(p.editValue)
				p.input.Focus()

			case "v":
				// Edit version
				p.mode = "edit"
				p.editKey = "version"
				p.editValue = p.project.Version
				p.input.SetValue(p.editValue)
				p.input.Focus()

			case "d":
				// Edit description
				p.mode = "edit"
				p.editKey = "description"
				p.editValue = p.project.Description
				p.input.SetValue(p.editValue)
				p.input.Focus()

			case "a":
				// Edit author
				p.mode = "edit"
				p.editKey = "author"
				p.editValue = p.project.Author
				p.input.SetValue(p.editValue)
				p.input.Focus()

			case "l":
				// Edit license
				p.mode = "edit"
				p.editKey = "license"
				p.editValue = p.project.License
				p.input.SetValue(p.editValue)
				p.input.Focus()
			}

		case "edit":
			// Handle edit mode keys
			switch msg.String() {
			case "enter":
				// Save changes
				p.loading = true

				// Update the project
				switch p.editKey {
				case "name":
					p.project.Name = p.input.Value()
				case "version":
					p.project.Version = p.input.Value()
				case "description":
					p.project.Description = p.input.Value()
				case "author":
					p.project.Author = p.input.Value()
				case "license":
					p.project.License = p.input.Value()
				}

				// Save changes
				go func() {
					if err := p.project.SavePackageJSON(); err != nil {
						p.error = fmt.Sprintf("Error saving package.json: %v", err)
					}
					p.mode = "view"
					p.loading = false
				}()

			case "esc":
				// Cancel edit
				p.mode = "view"
				p.input.SetValue("")
			}

			// Update the input
			p.input, cmd = p.input.Update(msg)
		}
	}

	return p, cmd
}

// View renders the panel
func (p *ProjectPanel) View() string {
	// In a 4-panel grid, we need to be more economical with space
	if p.mode == "edit" {
		// Simple edit view
		return fmt.Sprintf("%s:\n%s\n[↵]Save [esc]Cancel",
			p.editKey,
			p.input.View())
	}

	// Normal view
	if p.error != "" {
		return ErrorStyle.Render(p.error)
	}

	if p.loading {
		return "Loading..."
	}

	// Show compact project info
	var details string
	details += fmt.Sprintf("Name: %s\n", p.project.Name)
	details += fmt.Sprintf("Ver: %s\n", p.project.Version)

	// Show just a few more fields to avoid overflowing
	if p.project.Description != "" {
		desc := p.project.Description
		if len(desc) > p.width-10 {
			desc = desc[:p.width-13] + "..."
		}
		details += fmt.Sprintf("Desc: %s", desc)
	}

	return fmt.Sprintf("%s\n\n[e]Edit", details)
}

// Width returns the panel width
func (p *ProjectPanel) Width() int {
	return p.width
}

// Height returns the panel height
func (p *ProjectPanel) Height() int {
	return p.height
}

// SetSize sets the panel size
func (p *ProjectPanel) SetSize(width, height int) {
	p.width = width
	p.height = height
}

// Title returns the panel title
func (p *ProjectPanel) Title() string {
	return p.title
}
