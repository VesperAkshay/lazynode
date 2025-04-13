package npx

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

// NpxCommand represents an npx command
type NpxCommand struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Command     string `json:"command"`
}

// Runner handles npx command operations
type Runner struct {
	ProjectDir     string
	CacheFile      string
	RecentCommands []NpxCommand
}

// NewRunner creates a new npx runner
func NewRunner(projectDir string) (*Runner, error) {
	// Create the cache file in the project directory
	cacheFile := filepath.Join(projectDir, ".lazynode", "npx-cache.json")

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(cacheFile), os.ModePerm); err != nil {
		return nil, err
	}

	runner := &Runner{
		ProjectDir: projectDir,
		CacheFile:  cacheFile,
	}

	// Load the cache file if it exists
	if _, err := os.Stat(cacheFile); err == nil {
		if err := runner.LoadCache(); err != nil {
			return nil, err
		}
	}

	return runner, nil
}

// RunCommand runs an npx command
func (r *Runner) RunCommand(command string) (*exec.Cmd, error) {
	// Split the command into parts
	args := []string{}

	// Add the command and any arguments
	args = append(args, command)

	// Create the command
	cmd := exec.Command("npx", args...)

	// Set up the working directory
	cmd.Dir = r.ProjectDir

	// Redirect output to pipes
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

	// Cache the command
	r.CacheCommand(command, "")

	return cmd, nil
}

// GetRecentCommands returns the recent commands
func (r *Runner) GetRecentCommands() []NpxCommand {
	return r.RecentCommands
}

// CacheCommand adds a command to the cache
func (r *Runner) CacheCommand(command, description string) error {
	// Check if the command already exists in the cache
	for i, cmd := range r.RecentCommands {
		if cmd.Command == command {
			// Move the command to the top
			r.RecentCommands = append(r.RecentCommands[:i], r.RecentCommands[i+1:]...)
			r.RecentCommands = append([]NpxCommand{{
				Name:        filepath.Base(command),
				Description: description,
				Command:     command,
			}}, r.RecentCommands...)

			return r.SaveCache()
		}
	}

	// Add the command to the top
	r.RecentCommands = append([]NpxCommand{{
		Name:        filepath.Base(command),
		Description: description,
		Command:     command,
	}}, r.RecentCommands...)

	// Limit the cache to 20 commands
	if len(r.RecentCommands) > 20 {
		r.RecentCommands = r.RecentCommands[:20]
	}

	return r.SaveCache()
}

// LoadCache loads the cache from the cache file
func (r *Runner) LoadCache() error {
	// Read the cache file
	data, err := os.ReadFile(r.CacheFile)
	if err != nil {
		return err
	}

	// Parse the cache file
	return json.Unmarshal(data, &r.RecentCommands)
}

// SaveCache saves the cache to the cache file
func (r *Runner) SaveCache() error {
	// Create the cache file
	data, err := json.MarshalIndent(r.RecentCommands, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.CacheFile, data, os.ModePerm)
}

// GetPopularCommands returns a list of popular npx commands
func GetPopularCommands() []NpxCommand {
	return []NpxCommand{
		{
			Name:        "create-react-app",
			Description: "Create React applications with no build configuration",
			Command:     "create-react-app",
		},
		{
			Name:        "vite",
			Description: "Next generation frontend tooling",
			Command:     "vite",
		},
		{
			Name:        "next",
			Description: "Next.js framework",
			Command:     "next",
		},
		{
			Name:        "tsc",
			Description: "TypeScript compiler",
			Command:     "tsc",
		},
		{
			Name:        "eslint",
			Description: "Find and fix problems in your JavaScript code",
			Command:     "eslint",
		},
		{
			Name:        "prettier",
			Description: "Code formatter",
			Command:     "prettier",
		},
		{
			Name:        "jest",
			Description: "JavaScript testing framework",
			Command:     "jest",
		},
	}
}
