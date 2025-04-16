package ui

import (
	"fmt"
	"strings"

	"github.com/VesperAkshay/lazynode/pkg/project"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectPanel displays and manages package.json
type ProjectPanel struct {
	title        string
	width        int
	height       int
	project      *project.Project
	mode         string // "view", "edit", "menu"
	editKey      string
	editValue    string
	input        textinput.Model
	error        string
	loading      bool
	cursor       int // Menu cursor position
	menuItems    []string
	toggleValues map[string]bool // For boolean toggles like "private"
}

// NewProjectPanel creates a new project panel
func NewProjectPanel(project *project.Project) *ProjectPanel {
	// Create the input model
	input := textinput.New()
	input.Focus()

	menuItems := []string{
		"name",
		"version",
		"description",
		"main",
		"author",
		"license",
		"private",
		"homepage",
		"repository",
		"keywords",
		"engines.node",
	}

	toggleValues := make(map[string]bool)
	toggleValues["private"] = project.Private

	return &ProjectPanel{
		title:        "Project",
		project:      project,
		mode:         "view",
		input:        input,
		menuItems:    menuItems,
		toggleValues: toggleValues,
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
				// Show edit menu
				p.mode = "menu"
				p.cursor = 0
				return p, nil
			}

		case "menu":
			// Handle menu mode keys
			switch msg.String() {
			case "up", "k":
				if p.cursor > 0 {
					p.cursor--
				}
			case "down", "j":
				if p.cursor < len(p.menuItems)-1 {
					p.cursor++
				}
			case "enter":
				// Select field to edit
				p.editKey = p.menuItems[p.cursor]

				// Special case for toggle fields
				if p.editKey == "private" {
					p.toggleValues["private"] = !p.toggleValues["private"]
					p.project.Private = p.toggleValues["private"]

					// Save changes immediately for toggle
					go func() {
						if err := p.project.SavePackageJSON(); err != nil {
							p.error = fmt.Sprintf("Error saving package.json: %v", err)
						}
					}()

					return p, nil
				}

				// For normal text fields
				p.mode = "edit"

				// Set initial value based on field
				switch p.editKey {
				case "name":
					p.editValue = p.project.Name
				case "version":
					p.editValue = p.project.Version
				case "description":
					p.editValue = p.project.Description
				case "main":
					p.editValue = p.project.Main
				case "author":
					p.editValue = p.project.Author
				case "license":
					p.editValue = p.project.License
				case "homepage":
					p.editValue = p.project.Homepage
				case "repository":
					// Get repository URL if it exists
					if p.project.Repository != nil {
						if url, ok := p.project.Repository["url"]; ok {
							p.editValue = url
						}
					}
				case "keywords":
					// Join keywords with commas
					if p.project.Keywords != nil {
						p.editValue = strings.Join(p.project.Keywords, ", ")
					}
				case "engines.node":
					// Get Node.js engine version if it exists
					if p.project.Engines != nil {
						if nodeVersion, ok := p.project.Engines["node"]; ok {
							p.editValue = nodeVersion
						}
					}
				}

				p.input.SetValue(p.editValue)
				p.input.Focus()

			case "esc":
				// Return to view mode
				p.mode = "view"
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
				case "main":
					p.project.Main = p.input.Value()
				case "author":
					p.project.Author = p.input.Value()
				case "license":
					p.project.License = p.input.Value()
				case "homepage":
					p.project.Homepage = p.input.Value()
				case "repository":
					// Create repository map if it doesn't exist
					if p.project.Repository == nil {
						p.project.Repository = make(map[string]string)
					}
					p.project.Repository["url"] = p.input.Value()
					// Also set type if it doesn't exist
					if _, ok := p.project.Repository["type"]; !ok {
						p.project.Repository["type"] = "git"
					}
				case "keywords":
					// Split keywords by comma and trim spaces
					keywordsStr := p.input.Value()
					if keywordsStr != "" {
						keywords := strings.Split(keywordsStr, ",")
						// Trim spaces
						for i, keyword := range keywords {
							keywords[i] = strings.TrimSpace(keyword)
						}
						p.project.Keywords = keywords
					} else {
						p.project.Keywords = nil
					}
				case "engines.node":
					// Create engines map if it doesn't exist
					if p.project.Engines == nil {
						p.project.Engines = make(map[string]string)
					}
					p.project.Engines["node"] = p.input.Value()
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
	// Styles for better UI
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fabd2f"))

	fieldStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#b8bb26"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ebdbb2"))

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#fe8019")).
		Background(lipgloss.Color("#3c3836"))

	switch p.mode {
	case "menu":
		// Menu view for selecting fields to edit
		var sb strings.Builder
		sb.WriteString(titleStyle.Render("Select field to edit:") + "\n\n")

		for i, item := range p.menuItems {
			// Show toggle status for boolean fields
			if item == "private" {
				status := "off"
				if p.toggleValues["private"] {
					status = "on"
				}

				itemText := fmt.Sprintf("%s: [%s]", item, status)

				if i == p.cursor {
					sb.WriteString(selectedStyle.Render(" • "+itemText) + "\n")
				} else {
					sb.WriteString(" • " + itemText + "\n")
				}
			} else {
				if i == p.cursor {
					sb.WriteString(selectedStyle.Render(" • "+item) + "\n")
				} else {
					sb.WriteString(" • " + item + "\n")
				}
			}
		}

		sb.WriteString("\n[↵]Select [esc]Cancel")
		return sb.String()

	case "edit":
		// Simple edit view
		return fmt.Sprintf("%s:\n%s\n\n[↵]Save [esc]Cancel",
			titleStyle.Render(p.editKey),
			p.input.View())
	}

	// Normal view
	if p.error != "" {
		return ErrorStyle.Render(p.error)
	}

	if p.loading {
		return "Loading..."
	}

	// Build detailed project info
	var sb strings.Builder

	// Function to add a field to the display
	addField := func(name, value string) {
		if value != "" {
			sb.WriteString(fieldStyle.Render(name+": ") + valueStyle.Render(value) + "\n")
		}
	}

	addField("Name", p.project.Name)
	addField("Version", p.project.Version)

	// Description with wrap
	if p.project.Description != "" {
		desc := p.project.Description
		if len(desc) > p.width-15 {
			desc = desc[:p.width-15] + "..."
		}
		addField("Description", desc)
	}

	addField("Main", p.project.Main)
	addField("Author", p.project.Author)
	addField("License", p.project.License)

	// Private status
	privateStatus := "No"
	if p.project.Private {
		privateStatus = "Yes"
	}
	addField("Private", privateStatus)

	// Repository info
	if p.project.Repository != nil {
		if url, ok := p.project.Repository["url"]; ok {
			addField("Repository", url)
		}
	}

	addField("Homepage", p.project.Homepage)

	// Keywords
	if p.project.Keywords != nil && len(p.project.Keywords) > 0 {
		keywordsStr := strings.Join(p.project.Keywords, ", ")
		if len(keywordsStr) > p.width-15 {
			keywordsStr = keywordsStr[:p.width-15] + "..."
		}
		addField("Keywords", keywordsStr)
	}

	// Engine versions
	if p.project.Engines != nil {
		if nodeVersion, ok := p.project.Engines["node"]; ok {
			addField("Node Engine", nodeVersion)
		}
	}

	sb.WriteString("\n[e]Edit Fields")

	return sb.String()
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
