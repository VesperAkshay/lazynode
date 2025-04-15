package npm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// LinkStatus represents the linking status of a package
type LinkStatus int

const (
	NotLinked LinkStatus = iota
	LinkedGlobally
	LinkedLocally
)

// LinkPackage links a package either globally or to the current project
func (p *PackageManager) LinkPackage(packageName string, global bool) error {
	var cmd *exec.Cmd

	if global {
		// Link the current package globally
		cmd = exec.Command("npm", "link")
		cmd.Dir = filepath.Dir(p.PackageJSONPath)
	} else {
		// Link an existing global package to the current project
		cmd = exec.Command("npm", "link", packageName)
		cmd.Dir = filepath.Dir(p.PackageJSONPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// UnlinkPackage unlinks a package either globally or from the current project
func (p *PackageManager) UnlinkPackage(packageName string, global bool) error {
	var cmd *exec.Cmd

	if global {
		// Unlink the current package globally
		cmd = exec.Command("npm", "unlink", "-g")
		cmd.Dir = filepath.Dir(p.PackageJSONPath)
	} else {
		// Unlink a package from the current project
		cmd = exec.Command("npm", "unlink", packageName)
		cmd.Dir = filepath.Dir(p.PackageJSONPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// GetLinkStatus checks if a package is linked
func (p *PackageManager) GetLinkStatus(packageName string) (LinkStatus, error) {
	// If we're in a background loading process and don't need immediate accuracy,
	// do a quick check by looking at the file system directly rather than running commands
	if len(packageName) > 0 {
		// Check node_modules for a symlink - much faster than running npm commands
		nodePath := filepath.Join(filepath.Dir(p.PackageJSONPath), "node_modules", packageName)
		if fileInfo, err := os.Lstat(nodePath); err == nil && fileInfo.Mode()&os.ModeSymlink != 0 {
			return LinkedLocally, nil
		}
	}

	// Only run these commands when needed for detailed information
	// (they're slow but provide accurate information)

	// Check if package is linked globally
	globalCmd := exec.Command("npm", "ls", "-g", "--depth=0", packageName)
	globalResult, _ := globalCmd.Output()

	isGlobal := len(globalResult) > 0 && globalCmd.ProcessState.ExitCode() == 0

	// Check if package is linked locally
	localCmd := exec.Command("npm", "ls", "--depth=0", "--link", packageName)
	localCmd.Dir = filepath.Dir(p.PackageJSONPath)
	localResult, _ := localCmd.Output()

	isLocal := len(localResult) > 0 && localCmd.ProcessState.ExitCode() == 0

	if isLocal {
		return LinkedLocally, nil
	} else if isGlobal {
		return LinkedGlobally, nil
	}

	return NotLinked, nil
}

// ListLinkedPackages returns a list of packages linked in the current project
// This is an expensive operation so use it sparingly
func (p *PackageManager) ListLinkedPackages() ([]Package, error) {
	linkedPackages := []Package{}

	// Quick file system check for symlinks - much faster than running npm commands
	nodeModulesPath := filepath.Join(filepath.Dir(p.PackageJSONPath), "node_modules")

	// Check if node_modules directory exists
	if info, err := os.Stat(nodeModulesPath); err == nil && info.IsDir() {
		// Read directory entries
		entries, err := os.ReadDir(nodeModulesPath)
		if err == nil {
			for _, entry := range entries {
				// Check if it's a symlink
				if entry.Type()&os.ModeSymlink != 0 {
					name := entry.Name()
					if pkg, exists := p.Packages[name]; exists {
						pkg.LinkStatus = LinkedLocally
						pkg.IsLinked = true
						linkedPackages = append(linkedPackages, pkg)
					}
				}
			}
			return linkedPackages, nil
		}
	}

	// Fallback to npm command if the file system check didn't work
	cmd := exec.Command("npm", "ls", "--depth=0", "--link", "--json")
	cmd.Dir = filepath.Dir(p.PackageJSONPath)

	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list linked packages: %v", err)
	}

	// TODO: Parse JSON output to get linked packages
	// For now, we'll just reload all packages and check link status

	err = p.LoadPackages()
	if err != nil {
		return nil, err
	}

	for _, pkg := range p.Packages {
		status, _ := p.GetLinkStatus(pkg.Name)
		if status == LinkedLocally {
			linkedPackages = append(linkedPackages, pkg)
		}
	}

	return linkedPackages, nil
}
