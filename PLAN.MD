# LazyNode - A TUI for Node.js, npm, and npx

LazyNode is an interactive terminal UI for managing Node.js projects — inspired by Lazygit. It lets you manage `npm`, `npx`, and `package.json` workflows without touching the command line. Fast, intuitive, and fully featured, it's designed to supercharge your Node development workflow.

---

## 🚀 Features

### 📦 npm Package Manager
- Install/uninstall packages (with autocomplete)
- View installed packages and versions
- Check and update outdated packages
- Show dependency tree (expandable/collapsible)
- Add/remove devDependencies, peerDependencies, etc.

### 🧪 Script Runner
- List and run npm scripts interactively
- Re-run, stop, and monitor logs in real-time
- Highlight failed runs with error output
- Hotkey for running common scripts (e.g., build/test/start)

### ⚙️ Project Explorer
- View package.json content in a structured panel
- Inspect metadata: name, version, author, engines, etc.
- Edit fields like name, version, description

### ⚡ npx Runner
- Search and run npx commands
- Autocomplete for popular tools (create-react-app, vite, etc.)
- Cache frequently used commands

### 🔁 Live Diff Viewer
- Inspect changes in package.json & package-lock.json
- Git-style diff viewer
- Undo changes (revert package versions, script edits, etc.)

### 🧩 Git Integration (Optional)
- Show git diff of changes to package files
- Commit, stash, and push directly

### 🖥️ UI
- Fully keyboard navigable (Vim-style or Lazygit-style)
- Mouse support (optional)
- Multi-panel layout (scripts, dependencies, logs, etc.)
- Light and dark themes

---

## 🏗 Development Roadmap

### ✅ Phase 1: MVP (Minimum Viable Product)
**Timeline: Week 1–2**
- Project detection (read package.json)
- Terminal-based UI layout with panels
- Script listing and execution (basic log output)
- Basic keyboard navigation (up/down/select)
- Install/uninstall npm packages

### 🚧 Phase 2: Core Features
**Timeline: Week 3–4**
- Dependency viewer (tree format)
- Autocomplete for installing packages
- View and update outdated packages
- Script runner with live logs and color support
- Editable package.json fields from UI
- Run npx with input or search

### 🔨 Phase 3: Advanced UI and Git Integration
**Timeline: Week 5–6**
- Diff viewer for package.json / lock files
- Git integration (diff, stage, commit, push)
- Tabbed or resizable panels
- Multi-theme support (light/dark/custom)
- Config file support (lazyconfig.json)

### 🌟 Phase 4: Final Polish & Packaging
**Timeline: Week 7–8**
- Polish and bug fixing
- Binary builds for Linux, macOS, Windows
- Publish to npm (for Node.js version)
- Documentation site with usage examples
- Community feature request system

---

## 💻 Usage

```bash
npx lazynode        # Run in existing Node project
lazynode            # If installed globally
```

### Keyboard Shortcuts (Lazygit-inspired)
- `j/k` or arrows: Navigate lists
- `enter`: Run selected script
- `i`: Install package
- `d`: Delete package
- `o`: Check outdated
- `u`: Update selected package
- `r`: Reload UI
- `e`: Edit package.json
- `q`: Quit

---

## 🧰 Tech Stack

- Language: Go
- TUI Library:
  - Go: Bubbletea + Lip Gloss + Bubbles

---

## 🤝 Contributing

Contributions welcome! Open issues, suggest features, or submit PRs. See the `CONTRIBUTING.md` for more info.

---

## 📄 License

MIT License. See `LICENSE` for details.

