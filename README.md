# LazyNode

<p align="center">
  <img src="docs/logo.png" alt="LazyNode Logo" width="300" />
</p>

<p align="center">
  <pre align="center">
  ██╗      █████╗ ███████╗██╗   ██╗███╗   ██╗ ██████╗ ██████╗ ███████╗
  ██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝████╗  ██║██╔═══██╗██╔══██╗██╔════╝
  ██║     ███████║  ███╔╝  ╚████╔╝ ██╔██╗ ██║██║   ██║██║  ██║█████╗  
  ██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║╚██╗██║██║   ██║██║  ██║██╔══╝  
  ███████╗██║  ██║███████╗   ██║   ██║ ╚████║╚██████╔╝██████╔╝███████╗
  ╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝
  </pre>
  <br>
  <em>A powerful Terminal UI for managing Node.js projects with style</em>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#ui-enhancements">UI Enhancements</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#keyboard-shortcuts">Shortcuts</a> •
  <a href="#panels-overview">Panels</a> •
  <a href="#development">Development</a> •
  <a href="#contributing">Contributing</a> •
  <a href="#license">License</a>
</p>

---

LazyNode is a feature-rich terminal user interface inspired by Lazygit and LazyNPM that helps you manage Node.js projects without having to remember or type npm/npx commands. Boost your productivity with an intuitive, keyboard-driven interface for all your Node.js development needs.

![LazyNode Screenshot](docs/screenshot.png)

## Features

### 📦 Package Management
- Install regular and dev dependencies with dedicated keybindings (`i` and `I`)
- Uninstall packages with confirmation dialog
- Check for outdated packages and update them
- View detailed package information
- Link and unlink packages with simple keystrokes
- Manage all dependencies through a visually appealing interface

### 🧪 Script Management
- List and run npm scripts interactively
- Monitor script execution in real-time
- View script logs with automatic timestamps
- Stop running scripts with a simple keystroke

### ⚡ NPX Integration
- Execute npx commands directly from the UI
- Access frequently used commands quickly
- Get suggestions for popular npx tools
- Save and reuse custom npx commands

### 🖥️ Modern Terminal UI
- Side-by-side panel layout for efficient navigation
- Real-time animations and progress indicators
- Fully keyboard navigable with intuitive shortcuts
- Smart navigation between panels with Alt+Arrow keys
- Context-sensitive command display
- Advanced fullscreen mode for focused work

## UI Enhancements

LazyNode features several UI enhancements to provide a rich terminal experience:

### 🎨 Vibrant Color Scheme
- A carefully crafted terminal color palette with distinct colors for different UI elements
- Bright accent colors for important actions and status indicators
- Gradients and color transitions for an elegant look-and-feel
- Enhanced syntax highlighting for code and command output

### 🔠 ASCII Art & Text Styling
- Beautiful ASCII art logo on splash and quit screens
- Stylized borders, shadows, and panels
- Rich text formatting with bold, underline, and italic styles
- Progress indicators with circle notation

### ⚡ Animations & Visual Effects
- Dynamic loading animations during operations
- Smooth fade-in/fade-out transitions between screens
- Matrix-style particle effects on the splash screen
- Interactive progress bars and spinners with multi-color gradients

### 🎭 Terminal Enhancements
- Terminal-style command prompt in the logs panel
- Colorized output based on message type (error, warning, success)
- Styled ASCII borders for terminal-authentic feel
- Animated input prompts and confirmation dialogs
- Multi-step loading process visualization

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

LazyNode uses intuitive keyboard shortcuts for efficient navigation and control:

### Global Navigation
| Key | Action |
|-----|--------|
| `↑` / `k` | Move up in current panel |
| `↓` / `j` | Move down in current panel |
| `←` / `h` | Move left (when applicable) |
| `→` / `l` | Move right (when applicable) |
| `Alt+↑` / `Alt+k` | Go to previous panel |
| `Alt+↓` / `Alt+j` | Go to next panel |
| `Alt+←` / `Alt+h` | Go to left panel |
| `Alt+→` / `Alt+l` | Go to right panel |
| `Tab` | Cycle through panels |
| `Shift+Tab` | Cycle backwards through panels |
| `1` | Switch to Scripts panel |
| `2` | Switch to Packages panel |
| `3` | Switch to Project panel |
| `4` | Switch to NPX panel |
| `5` | Switch to Logs panel |
| `?` | Toggle help screen |
| `q` | Quit with elegant exit animation |
| `F` | Toggle fullscreen mode |

### Package Management (Packages Panel)
| Key | Action |
|-----|--------|
| `i` | Install a package |
| `I` | Install as dev dependency |
| `d` | Uninstall selected package |
| `o` | Check for outdated packages |
| `u` | Update selected package |
| `l` | Link package |
| `L` | Unlink package |
| `a` | Install all dependencies |
| `m` | Check missing dependencies and install them |
| `x` | Show actions menu |
| `/` | Search packages |
| `Enter` | Select/activate package |
| `Space` | Toggle detailed view |
| `Esc` | Cancel current action |

### Script Management (Scripts Panel)
| Key | Action |
|-----|--------|
| `Enter` | Run selected script |
| `Ctrl+C` | Stop running script |
| `r` | Reload scripts list |
| `Space` | View script details |

### NPX Commands (NPX Panel)
| Key | Action |
|-----|--------|
| `n` | Create new NPX command |
| `Enter` | Run selected NPX command |
| `Esc` | Cancel current action |
| `Space` | View command details |

### Project Actions
| Key | Action |
|-----|--------|
| `e` | Edit package.json |
| `o` | Open in editor |
| `b` | Build project |
| `t` | Run tests |
| `p` | Publish package |

## Panels Overview

LazyNode's interface is divided into multiple panels, each with a specific purpose:

### 📜 Scripts Panel
Displays all available npm scripts from your package.json file. Select and run scripts with a single keystroke.
- View script commands without having to remember them
- Run and monitor script execution in real-time
- Get colorized output based on exit status

### 📦 Packages Panel
Shows all dependencies (regular and development) installed in your project. Install, update, and remove packages with ease.
- Visual indicators for outdated packages
- Quick access to version information
- Fast actions for common operations
- Link and unlink functionality
- Smart confirmation dialogs for destructive actions

### 🔍 Project Panel
Provides an overview of your project, including package.json details, Node.js version, and environment information.
- View project metadata at a glance
- Quick access to important project properties
- Environment variable management
- Project health indicators and statistics

### ⚡ NPX Panel
Execute NPX commands without leaving the terminal UI. Includes history and suggestions for popular commands.
- Smart command history
- Popular command suggestions
- Custom command creation
- Real-time execution feedback

### 🖥️ Terminal Panel
Displays real-time output from running scripts, package operations, and system messages with color-coded formatting.
- ANSI color support
- Command timestamping
- Scrollable history
- Automatic error highlighting

## Special Screens

### 🚀 Splash Screen
An animated ASCII art welcome screen that appears when LazyNode starts up.
- Features a dynamic progress bar with stage indicators
- Matrix-style falling characters animation
- Gradient color logo effect
- Multi-step loading process visualization
- Smooth transition to the main interface

### 👋 Quit Screen
A farewell animation when you exit LazyNode.
- Fade-in/fade-out effect with thank you message
- Randomized farewell messages
- Gradient color effects
- Particle animation effects

## Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/VesperAkshay/lazynode.git
cd lazynode

# Install dependencies
go mod download

# Build the binary
go build -o lazynode cmd/lazynode/main.go

# Run tests
go test ./...
```

### Project Structure

- `cmd/` - Main application entry point
- `pkg/` - Library code organized by domain
  - `ui/` - Terminal UI components
  - `npm/` - npm package management
  - `npx/` - npx command execution
  - `scripts/` - Script management
  - `project/` - Project configuration
  - `version/` - Version information

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the project
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add some amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<p align="center">Made with ❤️ by <a href="https://github.com/VesperAkshay">VesperAkshay</a></p> 