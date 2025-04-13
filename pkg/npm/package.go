package npm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Package represents an npm package
type Package struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	LatestVersion   string `json:"latestVersion,omitempty"`
	Description     string `json:"description,omitempty"`
	Type            string `json:"type"` // "dependency", "devDependency", "peerDependency", etc.
	PackageJSONPath string `json:"-"`
}

// PackageManager handles npm package operations
type PackageManager struct {
	PackageJSONPath string
	Packages        map[string]Package
}

// NewPackageManager creates a new package manager for the given project
func NewPackageManager(packageJSONPath string) (*PackageManager, error) {
	pm := &PackageManager{
		PackageJSONPath: packageJSONPath,
		Packages:        make(map[string]Package),
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

	// First try to use npm list
	if err := pm.loadPackagesFromNpmList(); err != nil {
		// If that fails, fall back to reading package.json directly
		return pm.loadPackagesFromPackageJSON()
	}

	return nil
}

// loadPackagesFromNpmList uses the npm list command to get package information
func (pm *PackageManager) loadPackagesFromNpmList() error {
	// Read package.json
	cmd := exec.Command("npm", "list", "--json", "--depth=0")
	cmd.Dir = strings.TrimSuffix(pm.PackageJSONPath, "package.json")

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

	// Process dependencies
	for name, info := range result.Dependencies {
		pm.Packages[name] = Package{
			Name:            name,
			Version:         info.Version,
			Type:            "dependency",
			PackageJSONPath: pm.PackageJSONPath,
		}
	}

	// Process dev dependencies
	for name, info := range result.DevDependencies {
		pm.Packages[name] = Package{
			Name:            name,
			Version:         info.Version,
			Type:            "devDependency",
			PackageJSONPath: pm.PackageJSONPath,
		}
	}

	return nil
}

// loadPackagesFromPackageJSON reads package.json directly to get package information
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
		}
	}

	// Process dev dependencies
	for name, version := range packageJSON.DevDependencies {
		pm.Packages[name] = Package{
			Name:            name,
			Version:         version,
			Type:            "devDependency",
			PackageJSONPath: pm.PackageJSONPath,
		}
	}

	return nil
}

// InstallPackage installs a new package
func (pm *PackageManager) InstallPackage(name string, isDev bool) error {
	args := []string{"install", name}
	if isDev {
		args = append(args, "--save-dev")
	} else {
		args = append(args, "--save")
	}

	// Set the working directory to the same directory as package.json
	cmd := exec.Command("npm", args...)
	cmd.Dir = strings.TrimSuffix(pm.PackageJSONPath, "package.json")

	// Capture both stdout and stderr
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm install error: %v - %s", err, string(output))
	}

	// Reload packages after installing
	return pm.LoadPackages()
}

// UninstallPackage removes a package
func (pm *PackageManager) UninstallPackage(name string) error {
	cmd := exec.Command("npm", "uninstall", name)
	cmd.Dir = strings.TrimSuffix(pm.PackageJSONPath, "package.json")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("npm uninstall error: %v - %s", err, string(output))
	}

	// Reload packages after uninstalling
	return pm.LoadPackages()
}

// CheckOutdatedPackages checks for outdated packages
func (pm *PackageManager) CheckOutdatedPackages() (map[string]string, error) {
	cmd := exec.Command("npm", "outdated", "--json")
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

// UpdatePackage updates a package to the latest version
func (pm *PackageManager) UpdatePackage(name string) error {
	cmd := exec.Command("npm", "update", name)
	cmd.Dir = strings.TrimSuffix(pm.PackageJSONPath, "package.json")

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
