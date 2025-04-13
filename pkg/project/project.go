package project

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Project represents a Node.js project
type Project struct {
	PackageJSONPath string
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	Main            string            `json:"main,omitempty"`
	Author          string            `json:"author,omitempty"`
	License         string            `json:"license,omitempty"`
	Private         bool              `json:"private,omitempty"`
	Repository      map[string]string `json:"repository,omitempty"`
	Homepage        string            `json:"homepage,omitempty"`
	Bugs            map[string]string `json:"bugs,omitempty"`
	Keywords        []string          `json:"keywords,omitempty"`
	Engines         map[string]string `json:"engines,omitempty"`
	packageJSON     map[string]interface{}
}

// NewProject creates a new project from a package.json file
func NewProject(packageJSONPath string) (*Project, error) {
	project := &Project{
		PackageJSONPath: packageJSONPath,
	}

	// Load the package.json
	if err := project.LoadPackageJSON(); err != nil {
		return nil, err
	}

	return project, nil
}

// LoadPackageJSON loads the package.json file into the project
func (p *Project) LoadPackageJSON() error {
	// Read package.json
	data, err := ioutil.ReadFile(p.PackageJSONPath)
	if err != nil {
		return err
	}

	// Parse package.json into a raw map first
	if err := json.Unmarshal(data, &p.packageJSON); err != nil {
		return err
	}

	// Then parse into the Project struct to get the main fields
	if err := json.Unmarshal(data, p); err != nil {
		return err
	}

	return nil
}

// SavePackageJSON saves the project to package.json
func (p *Project) SavePackageJSON() error {
	// Apply changes from the Project struct to the raw map
	p.packageJSON["name"] = p.Name
	p.packageJSON["version"] = p.Version

	if p.Description != "" {
		p.packageJSON["description"] = p.Description
	}

	if p.Main != "" {
		p.packageJSON["main"] = p.Main
	}

	if p.Author != "" {
		p.packageJSON["author"] = p.Author
	}

	if p.License != "" {
		p.packageJSON["license"] = p.License
	}

	if p.Private {
		p.packageJSON["private"] = p.Private
	}

	if p.Repository != nil {
		p.packageJSON["repository"] = p.Repository
	}

	if p.Homepage != "" {
		p.packageJSON["homepage"] = p.Homepage
	}

	if p.Bugs != nil {
		p.packageJSON["bugs"] = p.Bugs
	}

	if p.Keywords != nil {
		p.packageJSON["keywords"] = p.Keywords
	}

	if p.Engines != nil {
		p.packageJSON["engines"] = p.Engines
	}

	// Write to package.json
	data, err := json.MarshalIndent(p.packageJSON, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p.PackageJSONPath, data, os.ModePerm)
}

// GetPackageJSON returns the raw package.json data
func (p *Project) GetPackageJSON() map[string]interface{} {
	return p.packageJSON
}

// UpdateField updates a field in the package.json
func (p *Project) UpdateField(key string, value interface{}) error {
	// Update the raw package.json
	p.packageJSON[key] = value

	// Save the changes
	return p.SavePackageJSON()
}

// GetField gets a field from the package.json
func (p *Project) GetField(key string) (interface{}, bool) {
	value, ok := p.packageJSON[key]
	return value, ok
}

// Detect attempts to find a package.json in the current or parent directories
func Detect() (string, error) {
	// Start with the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if package.json exists in the current directory
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		return packageJSONPath, nil
	}

	// Look in parent directories (up to 5 levels)
	for i := 0; i < 5; i++ {
		dir = filepath.Dir(dir)
		packageJSONPath = filepath.Join(dir, "package.json")

		if _, err := os.Stat(packageJSONPath); err == nil {
			return packageJSONPath, nil
		}
	}

	return "", os.ErrNotExist
}
