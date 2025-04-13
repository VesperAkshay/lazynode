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

### Using Go (from source)

```bash
# Install directly using Go
go install github.com/VesperAkshay/lazynode/cmd/lazynode@latest

# Or clone and build manually
git clone https://github.com/VesperAkshay/lazynode.git
cd lazynode
go build -o lazynode cmd/lazynode/main.go

# Move the binary to your PATH (optional)
sudo mv lazynode /usr/local/bin/
```

### Pre-built Binaries

Download the pre-built binary for your platform from the [Releases](https://github.com/VesperAkshay/lazynode/releases) page.

#### Quick Install Scripts

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/VesperAkshay/lazynode/main/scripts/install.sh | bash
```

**Windows PowerShell (Run as Administrator):**
```powershell
iwr -useb https://raw.githubusercontent.com/VesperAkshay/lazynode/main/scripts/install.ps1 | iex
```

#### Manual Install

**Linux/macOS:**
```bash
# Download the latest release (example for Linux amd64)
curl -L https://github.com/VesperAkshay/lazynode/releases/latest/download/lazynode_linux_amd64 -o lazynode

# Make it executable
chmod +x lazynode

# Move to a directory in your PATH
sudo mv lazynode /usr/local/bin/
```

**Windows:**
Download the appropriate `.exe` file from the Releases page and add it to your PATH.

### Package Managers

#### Homebrew (macOS and Linux)
```bash
brew tap lazynode/lazynode
brew install lazynode
```

#### Scoop (Windows)
```powershell
scoop bucket add lazynode https://github.com/lazynode/scoop-bucket.git
scoop install lazynode
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