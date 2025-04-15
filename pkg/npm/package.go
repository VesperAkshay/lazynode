package npm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Package represents an npm package
type Package struct {
	Name            string     `json:"name"`
	Version         string     `json:"version"`
	LatestVersion   string     `json:"latestVersion,omitempty"`
	Description     string     `json:"description,omitempty"`
	Type            string     `json:"type"` // "dependency", "devDependency", "peerDependency", etc.
	PackageJSONPath string     `json:"-"`
	LinkStatus      LinkStatus `json:"linkStatus"`
	IsLinked        bool       `json:"isLinked"`
}

// PackageManager handles npm package operations
type PackageManager struct {
	PackageJSONPath string
	Packages        map[string]Package
	ProjectRoot     string
}

// NewPackageManager creates a new package manager for the given project
func NewPackageManager(packageJSONPath string) (*PackageManager, error) {
	pm := &PackageManager{
		PackageJSONPath: packageJSONPath,
		Packages:        make(map[string]Package),
		ProjectRoot:     strings.TrimSuffix(packageJSONPath, "package.json"),
	}

	// Load the initial packages
	if err := pm.LoadPackages(); err != nil {
		return nil, err
	}

	return pm, nil
}

// LoadPackages loads all packages from package.json
func (pm *PackageManager) LoadPackages() error {
	// Clear existing packages
	pm.Packages = make(map[string]Package)

	// First try reading directly from package.json as it's much faster
	err := pm.loadPackagesFromPackageJSON()
	if err != nil {
		return err
	}

	// Start a background goroutine to get more detailed package info
	go func() {
		// This will update packages with additional info like latest versions
		// but won't block the initial UI loading
		_ = pm.updatePackagesWithNpmList()
	}()

	return nil
}

// updatePackagesWithNpmList uses the npm list command to get detailed package information
// This runs asynchronously to avoid blocking the UI
func (pm *PackageManager) updatePackagesWithNpmList() error {
	// Use a faster npm command with limited depth
	cmd := exec.Command("npm", "list", "--json", "--depth=0")
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.Output()
	if err != nil {
		// npm list might return non-zero exit code but still provide useful output
		if len(output) == 0 {
			return fmt.Errorf("npm list error: %v", err)
		}
	}

	// Parse the output
	var result struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
			From    string `json:"from"`
		} `json:"dependencies"`
		DevDependencies map[string]struct {
			Version string `json:"version"`
			From    string `json:"from"`
		} `json:"devDependencies"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return err
	}

	// Update existing packages with more details
	for name, info := range result.Dependencies {
		if pkg, exists := pm.Packages[name]; exists {
			pkg.Version = info.Version
			pm.Packages[name] = pkg
		}
	}

	for name, info := range result.DevDependencies {
		if pkg, exists := pm.Packages[name]; exists {
			pkg.Version = info.Version
			pm.Packages[name] = pkg
		}
	}

	// Update link status for all packages
	for name, pkg := range pm.Packages {
		linkStatus, _ := pm.GetLinkStatus(name)
		pkg.LinkStatus = linkStatus
		pkg.IsLinked = (linkStatus != NotLinked)
		pm.Packages[name] = pkg
	}

	return nil
}

// loadPackagesFromNpmList is now a private method used by updatePackagesWithNpmList

// loadPackagesFromPackageJSON reads package.json directly to get package information
// This is much faster than running npm commands
func (pm *PackageManager) loadPackagesFromPackageJSON() error {
	// Read package.json content
	packageJSONContent, err := os.ReadFile(pm.PackageJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %v", err)
	}

	// Parse the JSON content
	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(packageJSONContent, &packageJSON); err != nil {
		return fmt.Errorf("failed to parse package.json: %v", err)
	}

	// Process dependencies
	for name, version := range packageJSON.Dependencies {
		pm.Packages[name] = Package{
			Name:            name,
			Version:         version,
			Type:            "dependency",
			PackageJSONPath: pm.PackageJSONPath,
			LinkStatus:      NotLinked,
			IsLinked:        false,
		}
	}

	// Process dev dependencies
	for name, version := range packageJSON.DevDependencies {
		pm.Packages[name] = Package{
			Name:            name,
			Version:         version,
			Type:            "devDependency",
			PackageJSONPath: pm.PackageJSONPath,
			LinkStatus:      NotLinked,
			IsLinked:        false,
		}
	}

	// If no packages were found in package.json, check if there are any installed in node_modules
	// This handles cases where packages are installed but not yet saved to package.json
	if len(pm.Packages) == 0 {
		nodeModulesPath := filepath.Join(pm.ProjectRoot, "node_modules")
		if _, err := os.Stat(nodeModulesPath); err == nil {
			// Read the directories in node_modules
			entries, err := os.ReadDir(nodeModulesPath)
			if err == nil {
				installedCount := 0
				for _, entry := range entries {
					// Skip hidden directories and files (like .bin)
					if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
						pm.Packages[entry.Name()] = Package{
							Name:            entry.Name(),
							Version:         "installed",  // We don't know the exact version
							Type:            "dependency", // Default to regular dependency
							PackageJSONPath: pm.PackageJSONPath,
							LinkStatus:      NotLinked,
							IsLinked:        false,
							Description:     "Installed but not in package.json",
						}
						installedCount++
					}
				}

				// If we found packages, log a message
				if installedCount > 0 {
					fmt.Printf("Found %d packages in node_modules that are not listed in package.json\n", installedCount)
				}
			}
		}
	}

	return nil
}

// InstallPackage installs a new package
func (pm *PackageManager) InstallPackage(name string, isDev bool) error {
	// Build the installation command with optimized flags
	args := []string{"install", name, "--no-fund", "--no-audit", "--prefer-offline"}
	if isDev {
		args = append(args, "--save-dev")
	} else {
		args = append(args, "--save")
	}

	// Set the working directory to the same directory as package.json
	cmd := exec.Command("npm", args...)
	cmd.Dir = pm.ProjectRoot

	// Setup pipes to capture output in real-time
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("npm install error: %v", err)
	}

	// Create a channel to signal when we should read the result
	done := make(chan struct{})

	// Collect output in background
	var output strings.Builder

	// Process stdout in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Process stderr in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Wait for both stdout and stderr to be fully read
	<-done
	<-done

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("npm install error: %v - %s", err, output.String())
	}

	// Create a lightweight initial entry in the Packages map
	// before the full LoadPackages operation
	if _, exists := pm.Packages[name]; !exists {
		// Add a basic entry for immediate UI feedback
		pkgType := "dependency"
		if isDev {
			pkgType = "devDependency"
		}

		pm.Packages[name] = Package{
			Name:            name,
			Version:         "latest", // Will be updated by LoadPackages
			Type:            pkgType,
			PackageJSONPath: pm.PackageJSONPath,
			LinkStatus:      NotLinked,
			IsLinked:        false,
		}
	}

	// Reload packages asynchronously to avoid blocking the UI
	go func() {
		// This performs a lightweight update without blocking
		_ = pm.loadPackagesFromPackageJSON()

		// Full detailed update in background
		_ = pm.updatePackagesWithNpmList()
	}()

	return nil
}

// UninstallPackage removes a package
func (pm *PackageManager) UninstallPackage(name string) error {
	// Use optimized flags for faster uninstallation
	cmd := exec.Command("npm", "uninstall", name, "--no-fund", "--no-audit")
	cmd.Dir = pm.ProjectRoot

	// Setup pipes to capture output in real-time
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("npm uninstall error: %v", err)
	}

	// Create a channel to signal when we should read the result
	done := make(chan struct{})

	// Collect output in background
	var output strings.Builder

	// Process stdout in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Process stderr in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Wait for both stdout and stderr to be fully read
	<-done
	<-done

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("npm uninstall error: %v - %s", err, output.String())
	}

	// Remove the package from our cache immediately for fast UI feedback
	delete(pm.Packages, name)

	// Reload packages asynchronously to avoid blocking the UI
	go func() {
		// This performs a lightweight update without blocking
		_ = pm.loadPackagesFromPackageJSON()
	}()

	return nil
}

// CheckOutdatedPackages checks for outdated packages
func (pm *PackageManager) CheckOutdatedPackages() (map[string]string, error) {
	cmd := exec.Command("npm", "outdated", "--json")
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.Output()
	if err != nil {
		// npm outdated returns a non-zero exit code if outdated packages are found
		// so we need to check if we got any output
		if len(output) == 0 {
			return nil, err
		}
	}

	// Parse the output
	var outdated map[string]struct {
		Current  string `json:"current"`
		Latest   string `json:"latest"`
		Wanted   string `json:"wanted"`
		Location string `json:"location"`
	}

	if err := json.Unmarshal(output, &outdated); err != nil {
		return nil, err
	}

	// Create a map of package name to latest version
	result := make(map[string]string)
	for name, info := range outdated {
		result[name] = info.Latest
		// Update the package in our cache
		if pkg, ok := pm.Packages[name]; ok {
			pkg.LatestVersion = info.Latest
			pm.Packages[name] = pkg
		}
	}

	return result, nil
}

// InstallAllDependencies installs all dependencies listed in package.json
func (pm *PackageManager) InstallAllDependencies() error {
	// Build the installation command with optimized flags
	args := []string{"install", "--no-fund", "--no-audit", "--prefer-offline"}

	// Set the working directory to the same directory as package.json
	cmd := exec.Command("npm", args...)
	cmd.Dir = pm.ProjectRoot

	// Setup pipes to capture output in real-time
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("npm install error: %v", err)
	}

	// Create a channel to signal when we should read the result
	done := make(chan struct{})

	// Collect output in background
	var output strings.Builder

	// Process stdout in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Process stderr in a goroutine
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- struct{}{}
	}()

	// Wait for both stdout and stderr to be fully read
	<-done
	<-done

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("npm install error: %v - %s", err, output.String())
	}

	// Reload packages asynchronously to avoid blocking the UI
	go func() {
		// This performs a lightweight update without blocking
		_ = pm.loadPackagesFromPackageJSON()

		// Full detailed update in background
		_ = pm.updatePackagesWithNpmList()
	}()

	return nil
}

// CheckMissingDependencies checks if there are dependencies in package.json that aren't installed
func (pm *PackageManager) CheckMissingDependencies() (bool, []string, error) {
	// Read package.json content
	packageJSONContent, err := os.ReadFile(pm.PackageJSONPath)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read package.json: %v", err)
	}

	// Parse the JSON content
	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(packageJSONContent, &packageJSON); err != nil {
		return false, nil, fmt.Errorf("failed to parse package.json: %v", err)
	}

	// Get the list of installed packages
	cmd := exec.Command("npm", "list", "--json", "--depth=0")
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.Output()
	// npm list might return non-zero exit code but still provide useful output
	if err != nil && len(output) == 0 {
		return false, nil, fmt.Errorf("npm list error: %v", err)
	}

	// Parse the output
	var installedResult struct {
		Dependencies map[string]interface{} `json:"dependencies"`
	}

	if err := json.Unmarshal(output, &installedResult); err != nil {
		return false, nil, err
	}

	// Check if all dependencies from package.json are installed
	missingDeps := []string{}

	// If there are no installed dependencies but package.json has some,
	// we know they're all missing
	if installedResult.Dependencies == nil &&
		(len(packageJSON.Dependencies) > 0 || len(packageJSON.DevDependencies) > 0) {

		// Add all dependencies to the missing list
		for name := range packageJSON.Dependencies {
			missingDeps = append(missingDeps, name)
		}
		for name := range packageJSON.DevDependencies {
			missingDeps = append(missingDeps, name)
		}

		return true, missingDeps, nil
	}

	// Check regular dependencies
	for name := range packageJSON.Dependencies {
		if _, exists := installedResult.Dependencies[name]; !exists {
			missingDeps = append(missingDeps, name)
		}
	}

	// Check dev dependencies
	for name := range packageJSON.DevDependencies {
		if _, exists := installedResult.Dependencies[name]; !exists {
			missingDeps = append(missingDeps, name)
		}
	}

	return len(missingDeps) > 0, missingDeps, nil
}

// InstallMissingDependencies checks for missing dependencies and installs them if needed
// Returns true if dependencies were installed, false if nothing needed to be installed
func (pm *PackageManager) InstallMissingDependencies() (bool, error) {
	hasMissing, _, err := pm.CheckMissingDependencies()
	if err != nil {
		return false, fmt.Errorf("failed to check for missing dependencies: %v", err)
	}

	// If no missing dependencies, return early
	if !hasMissing {
		return false, nil
	}

	// Install all dependencies from package.json
	if err := pm.InstallAllDependencies(); err != nil {
		return false, fmt.Errorf("failed to install missing dependencies: %v", err)
	}

	return true, nil
}

// UpdatePackage updates a package to the latest version
func (pm *PackageManager) UpdatePackage(name string) error {
	cmd := exec.Command("npm", "update", name)
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm update error: %v - %s", err, string(output))
	}

	// Reload packages after updating
	return pm.LoadPackages()
}

// SearchPackage searches for packages on npm registry
func SearchPackage(query string) ([]Package, error) {
	cmd := exec.Command("npm", "search", query, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse the output
	var searchResults []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(output, &searchResults); err != nil {
		return nil, err
	}

	// Convert to our Package type
	var packages []Package
	for _, result := range searchResults {
		packages = append(packages, Package{
			Name:        result.Name,
			Version:     result.Version,
			Description: result.Description,
		})
	}

	return packages, nil
}

// GetDependencyTree returns the dependency tree
func (pm *PackageManager) GetDependencyTree() (string, error) {
	cmd := exec.Command("npm", "list", "--depth=1")
	cmd.Dir = strings.TrimSuffix(pm.PackageJSONPath, "package.json")

	output, err := cmd.Output()
	if err != nil {
		// npm list may return a non-zero exit code for some issues
		// but we still want to display the output
		if len(output) == 0 {
			return "", err
		}
	}

	// Return the tree as a string
	return strings.TrimSpace(string(output)), nil
}

// BuildPackage runs the build script for a package
func (pm *PackageManager) BuildPackage() error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm build error: %v - %s", err, string(output))
	}

	return nil
}

// TestPackage runs the test script for a package
func (pm *PackageManager) TestPackage() error {
	cmd := exec.Command("npm", "test")
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm test error: %v - %s", err, string(output))
	}

	return nil
}

// PublishPackage publishes a package to the npm registry
func (pm *PackageManager) PublishPackage(dryRun bool) error {
	args := []string{"publish"}
	if dryRun {
		args = append(args, "--dry-run")
	}

	cmd := exec.Command("npm", args...)
	cmd.Dir = pm.ProjectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm publish error: %v - %s", err, string(output))
	}

	return nil
}

// OpenEditor opens the package.json in the default editor
func (pm *PackageManager) OpenEditor() error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad", pm.PackageJSONPath)
	case "darwin":
		cmd = exec.Command("open", pm.PackageJSONPath)
	default: // Linux and others
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vi" // Default to vi if EDITOR is not set
		}
		cmd = exec.Command(editor, pm.PackageJSONPath)
	}

	// Start the editor in a separate process so we don't wait for it to close
	return cmd.Start()
}
