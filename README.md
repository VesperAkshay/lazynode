# LazyNode - A TUI for Node.js, npm, and npx

LazyNode is an interactive terminal UI for managing Node.js projects — inspired by Lazygit. It lets you manage `npm`, `npx`, and `package.json` workflows without touching the command line. Fast, intuitive, and fully featured, it's designed to supercharge your Node development workflow.

![LazyNode Screenshot](screenshot.png)

## Features

- 📦 **npm Package Manager**
  - Install/uninstall packages
  - View installed packages and versions
  - Check and update outdated packages

- 🧪 **Script Runner**
  - List and run npm scripts interactively
  - View script logs in real-time

- ⚙️ **Project Explorer**
  - View package.json content in a structured panel
  - Edit fields like name, version, description

- ⚡ **npx Runner**
  - Run npx commands
  - Quick access to popular tools and recently used commands

- 🖥️ **TUI**
  - Fully keyboard navigable
  - Tab-based interface
  - Color-coded panels

## Installation

### Prerequisites

- Go 1.19 or higher
- Node.js and npm installed

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/lazynode.git
cd lazynode

# Build the project
go build -o lazynode ./cmd/lazynode

# Run LazyNode
./lazynode
```

## Usage

Run LazyNode inside any Node.js project directory:

```bash
lazynode
```

### Keyboard Shortcuts

Navigation:
- `up`/`k`: Move cursor up
- `down`/`j`: Move cursor down
- `1-5`: Switch between panels
- `?`: Toggle help

Package Management:
- `i`: Install a package
- `d`: Delete a package
- `o`: Check for outdated packages
- `u`: Update a package

Scripts:
- `enter`: Run the selected script

Project:
- `e`: Edit name
- `v`: Edit version
- `d`: Edit description
- `a`: Edit author
- `l`: Edit license

npx Commands:
- `n`: Run a new npx command
- `enter`: Run the selected command

General:
- `r`: Reload UI
- `q`: Quit

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

## License

MIT License 