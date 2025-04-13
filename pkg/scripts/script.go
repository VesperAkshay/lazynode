package scripts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Script represents an npm script
type Script struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}

// ScriptRunner handles npm script operations
type ScriptRunner struct {
	PackageJSONPath string
	Scripts         []Script
	RunningScripts  map[string]*exec.Cmd
}

// NewScriptRunner creates a new script runner for the given project
func NewScriptRunner(packageJSONPath string) (*ScriptRunner, error) {
	sr := &ScriptRunner{
		PackageJSONPath: packageJSONPath,
		Scripts:         []Script{},
		RunningScripts:  make(map[string]*exec.Cmd),
	}

	// Load the initial scripts
	if err := sr.LoadScripts(); err != nil {
		return nil, err
	}

	return sr, nil
}

// LoadScripts loads all scripts from package.json
func (sr *ScriptRunner) LoadScripts() error {
	// Clear the existing scripts
	sr.Scripts = []Script{}

	// Read package.json
	data, err := ioutil.ReadFile(sr.PackageJSONPath)
	if err != nil {
		return err
	}

	// Parse package.json
	var packageJSON struct {
		Scripts map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Extract the scripts
	for name, command := range packageJSON.Scripts {
		sr.Scripts = append(sr.Scripts, Script{
			Name:    name,
			Command: command,
		})
	}

	return nil
}

// RunScript runs a script by name
func (sr *ScriptRunner) RunScript(name string) (*exec.Cmd, error) {
	// Check if the script is already running
	if _, ok := sr.RunningScripts[name]; ok {
		return nil, nil // Already running
	}

	// Run the script using npm run
	cmd := exec.Command("npm", "run", name)

	// Set up the working directory to the directory containing package.json
	cmd.Dir = filepath.Dir(sr.PackageJSONPath)

	// Redirect output to pipes - we don't actually read from them in this function
	// but setting them up allows capturing output elsewhere
	_, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	_, err = cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// Store the running command
	sr.RunningScripts[name] = cmd

	return cmd, nil
}

// StopScript stops a running script
func (sr *ScriptRunner) StopScript(name string) error {
	cmd, ok := sr.RunningScripts[name]
	if !ok {
		return nil // Not running
	}

	// Kill the process
	if err := cmd.Process.Kill(); err != nil {
		return err
	}

	// Remove from running scripts
	delete(sr.RunningScripts, name)

	return nil
}

// AddScript adds a new script to package.json
func (sr *ScriptRunner) AddScript(name, command string) error {
	// Read package.json
	data, err := ioutil.ReadFile(sr.PackageJSONPath)
	if err != nil {
		return err
	}

	// Parse package.json
	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Get or create the scripts section
	scripts, ok := packageJSON["scripts"].(map[string]interface{})
	if !ok {
		scripts = make(map[string]interface{})
		packageJSON["scripts"] = scripts
	}

	// Add the script
	scripts[name] = command

	// Write back to package.json
	updatedData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(sr.PackageJSONPath, updatedData, os.ModePerm); err != nil {
		return err
	}

	// Reload scripts
	return sr.LoadScripts()
}

// RemoveScript removes a script from package.json
func (sr *ScriptRunner) RemoveScript(name string) error {
	// Read package.json
	data, err := ioutil.ReadFile(sr.PackageJSONPath)
	if err != nil {
		return err
	}

	// Parse package.json
	var packageJSON map[string]interface{}
	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return err
	}

	// Get the scripts section
	scripts, ok := packageJSON["scripts"].(map[string]interface{})
	if !ok {
		return nil // No scripts section
	}

	// Remove the script
	delete(scripts, name)

	// Write back to package.json
	updatedData, err := json.MarshalIndent(packageJSON, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(sr.PackageJSONPath, updatedData, os.ModePerm); err != nil {
		return err
	}

	// Reload scripts
	return sr.LoadScripts()
}
