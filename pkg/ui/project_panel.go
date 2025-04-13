package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lazynode/lazynode/pkg/project"
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
	// Show loading or error
	if p.loading {
		return PanelStyle.Width(p.width).Height(p.height).Render(
			fmt.Sprintf("%s\n\nSaving...", TitleStyle.Render(p.title)),
		)
	}

	if p.error != "" {
		return PanelStyle.Width(p.width).Height(p.height).Render(
			fmt.Sprintf("%s\n\n%s", TitleStyle.Render(p.title), ErrorStyle.Render(p.error)),
		)
	}

	// Create the package.json view
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#83a598")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ebdbb2"))

	// Show edit or view
	if p.mode == "edit" {
		return PanelStyle.Width(p.width).Height(p.height).Render(
			fmt.Sprintf("%s\n\nEditing %s\n\n%s",
				TitleStyle.Render(p.title),
				keyStyle.Render(p.editKey),
				p.input.View(),
			),
		)
	}

	// Build the package.json view
	var content strings.Builder

	content.WriteString(fmt.Sprintf("%s: %s\n",
		keyStyle.Render("name"),
		valueStyle.Render(p.project.Name),
	))

	content.WriteString(fmt.Sprintf("%s: %s\n",
		keyStyle.Render("version"),
		valueStyle.Render(p.project.Version),
	))

	if p.project.Description != "" {
		content.WriteString(fmt.Sprintf("%s: %s\n",
			keyStyle.Render("description"),
			valueStyle.Render(p.project.Description),
		))
	}

	if p.project.Author != "" {
		content.WriteString(fmt.Sprintf("%s: %s\n",
			keyStyle.Render("author"),
			valueStyle.Render(p.project.Author),
		))
	}

	if p.project.License != "" {
		content.WriteString(fmt.Sprintf("%s: %s\n",
			keyStyle.Render("license"),
			valueStyle.Render(p.project.License),
		))
	}

	// Add help text
	content.WriteString("\nPress the key to edit a field:\n")
	content.WriteString("(e) name, (v) version, (d) description, (a) author, (l) license")

	return PanelStyle.Width(p.width).Height(p.height).Render(
		fmt.Sprintf("%s\n\n%s", TitleStyle.Render(p.title), content.String()),
	)
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
