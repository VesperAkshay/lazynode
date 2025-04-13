# LazyNode

<p align="center">
  <img src="docs/logo.png" alt="LazyNode Logo" width="200" />
  <br>
  <em>A powerful Terminal UI for managing Node.js projects</em>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#keyboard-shortcuts">Shortcuts</a> •
  <a href="#development">Development</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

LazyNode is a feature-rich terminal user interface inspired by Lazygit that helps you manage Node.js projects without having to remember or type npm/npx commands. Boost your productivity with an intuitive, keyboard-driven interface for all your Node.js development needs.

![LazyNode Screenshot](docs/screenshot.png)

## Features

### 📦 Package Management
- Install regular and dev dependencies with a few keystrokes
- Uninstall packages with confirmation dialog
- Check for outdated packages and update them
- View detailed package information
- Manage all dependencies through a visually appealing interface

### 🧪 Script Management
- List and run npm scripts interactively
- Monitor script execution in real-time
- View script logs with automatic timestamps

### ⚡ NPX Integration
- Execute npx commands directly from the UI
- Access frequently used commands quickly
- Get suggestions for popular npx tools

### 🖥️ Modern Terminal UI
- Side-by-side panel layout for efficient navigation
- Real-time animations and progress indicators
- Fully keyboard navigable with intuitive shortcuts
- Smart navigation between panels (Alt+Up/Down) and within content (Up/Down)
- Context-sensitive command display

## Installation

### Quick Install (Linux/macOS)

```bash
# Install the latest version with one command
curl -sSL https://raw.githubusercontent.com/yourusername/lazynode/main/scripts/install.sh | bash
```

### From Pre-built Binaries (Recommended)

1. Visit the [Releases](https://github.com/yourusername/lazynode/releases) page
2. Download the appropriate package for your operating system:
   - Windows: `lazynode_<version>_windows_amd64.zip`
   - macOS (Intel): `lazynode_<version>_darwin_amd64.tar.gz`
   - macOS (Apple Silicon): `lazynode_<version>_darwin_arm64.tar.gz`
   - Linux (x86_64): `lazynode_<version>_linux_amd64.tar.gz`
   - Linux (ARM64): `lazynode_<version>_linux_arm64.tar.gz`
3. Extract the archive and make the binary executable

#### Windows
```powershell
# Extract the zip file using Windows Explorer or:
Expand-Archive lazynode_<version>_windows_amd64.zip -DestinationPath C:\some\directory

# Add the directory to your PATH or move the executable somewhere in your PATH
```

#### macOS/Linux
```bash
# Extract
tar xzvf lazynode_<version>_<os>_<arch>.tar.gz

# Make executable
chmod +x lazynode_<version>_<os>_<arch>

# Move to a directory in your PATH
sudo mv lazynode_<version>_<os>_<arch> /usr/local/bin/lazynode
```

### Using npm (Node.js)

```bash
# Install globally via npm
npm install -g lazynode

# Or use npx to run without installing
npx lazynode
```

### Using Go Modules

```bash
# Install using go install
go install github.com/yourusername/lazynode/cmd/lazynode@latest
```

### Arch Linux (AUR)

```bash
# Using an AUR helper like yay
yay -S lazynode

# Or manually from the AUR
git clone https://aur.archlinux.org/lazynode.git
cd lazynode
makepkg -si
```

### Windows Package Managers

#### Winget
```powershell
# Install using Windows Package Manager
winget install lazynode
```

#### PowerShell
```powershell
# Install using PowerShell module
Install-Module -Name LazyNode
Import-Module LazyNode
```

#### Chocolatey
```powershell
# Install using Chocolatey
choco install lazynode
```

#### Scoop
```powershell
# Install using Scoop
scoop install lazynode
```

### Using Go (from source)

Prerequisites:
- Go 1.19 or higher
- Git

```bash
# Clone the repository
git clone https://github.com/yourusername/lazynode.git

# Build the project
cd lazynode
make build

# Install (optional)
sudo make install
```

### Verify Installation

Verify the installation by running:

```bash
lazynode --version
```

This should display the version information for LazyNode.

## Publishing to npm

Even though LazyNode is built with Go, you can publish it to npm to make it installable via npm/npx. Follow these steps:

### Prerequisites

1. Create an npm account at https://www.npmjs.com/signup
2. Build your Go binaries for all platforms (Windows, macOS, Linux)
3. Create GitHub releases with the binaries

### Publishing Steps

1. **Login to npm**:
   ```bash
   npm login
   ```

2. **Verify package.json settings**:
   - Update the version in package.json
   - Make sure your GitHub repository is correctly referenced
   - Update author and other metadata as needed

3. **Build Release Artifacts**:
   ```bash
   # Build for all platforms and create release archives
   make release
   
   # Or manually:
   GOOS=darwin GOARCH=amd64 go build -o ./dist/lazynode_0.1.0_darwin_amd64/lazynode
   GOOS=darwin GOARCH=arm64 go build -o ./dist/lazynode_0.1.0_darwin_arm64/lazynode
   GOOS=linux GOARCH=amd64 go build -o ./dist/lazynode_0.1.0_linux_amd64/lazynode
   GOOS=windows GOARCH=amd64 go build -o ./dist/lazynode_0.1.0_windows_amd64/lazynode.exe
   
   # Create archives
   tar -czf ./dist/lazynode_0.1.0_darwin_amd64.tar.gz -C ./dist/lazynode_0.1.0_darwin_amd64 lazynode
   tar -czf ./dist/lazynode_0.1.0_darwin_arm64.tar.gz -C ./dist/lazynode_0.1.0_darwin_arm64 lazynode
   tar -czf ./dist/lazynode_0.1.0_linux_amd64.tar.gz -C ./dist/lazynode_0.1.0_linux_amd64 lazynode
   zip -j ./dist/lazynode_0.1.0_windows_amd64.zip ./dist/lazynode_0.1.0_windows_amd64/lazynode.exe
   ```

4. **Create a GitHub Release**:
   - Create a new GitHub release for version v0.1.0
   - Upload all archive files (.tar.gz and .zip)

5. **Test the package locally**:
   ```bash
   # Pack without publishing
   npm pack
   
   # Install the local package
   npm install -g ./lazynode-0.1.0.tgz
   ```

6. **Publish to npm**:
   ```bash
   # Publish the package
   npm publish
   
   # For scoped packages (optional)
   npm publish --access public
   ```

7. **To update the package**:
   ```bash
   # Update version in package.json
   npm version patch  # or minor or major
   
   # Build new binaries and create a new GitHub release
   
   # Publish the new version
   npm publish
   ```

Once published, users can install LazyNode using:
```bash
npm install -g lazynode
# or
npx lazynode
```

## Usage

Navigate to your Node.js project directory and run:

```bash
lazynode
```

This will launch the LazyNode interface, automatically detecting your project's package.json file.

## Keyboard Shortcuts

LazyNode uses intuitive keyboard shortcuts inspired by Lazygit:

### Navigation
- `↑/↓` or `j/k`: Navigate within a panel
- `alt+↑/↓`: Switch between panels
- `tab`: Cycle through panels
- `enter`: Select/activate item

### Package Management
- `a`: Show all package actions
- `i`: Install a package
- `d`: Uninstall selected package
- `o`: Check for outdated packages
- `u`: Update selected package
- `/`: Search packages

### General
- `?`: Show help
- `r`: Refresh/reload
- `q`: Quit LazyNode

## Configuration

LazyNode stores its configuration in a `.lazynode` directory in your project root. This includes:

- NPX command history
- UI preferences and state

## Development

LazyNode is built with Go and uses the following libraries:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea): Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss): Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles): UI components

To set up the development environment:

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/lazynode/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Lazygit](https://github.com/jesseduffield/lazygit)
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Thanks to all contributors and the Node.js community! 