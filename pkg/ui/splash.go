package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ASCII art logo for LazyNode
const asciiLogo = `
██╗      █████╗ ███████╗██╗   ██╗███╗   ██╗ ██████╗ ██████╗ ███████╗
██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝████╗  ██║██╔═══██╗██╔══██╗██╔════╝
██║     ███████║  ███╔╝  ╚████╔╝ ██╔██╗ ██║██║   ██║██║  ██║█████╗  
██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║╚██╗██║██║   ██║██║  ██║██╔══╝  
███████╗██║  ██║███████╗   ██║   ██║ ╚████║╚██████╔╝██████╔╝███████╗
╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝
`

// Particle represents a moving particle in the splash background
type Particle struct {
	x, y     float64
	vx, vy   float64
	char     string
	color    lipgloss.Color
	lifespan int
}

// Terminal style animation frames for the splash screen
var loadingFrames = []string{
	"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷",
}

// Matrix-style characters
var matrixChars = []string{
	"0", "1", "ツ", "ノ", "ス", "ミ", "サ", "㋚", "ム", "$", "#", "%", "&", "+",
	"ヌ", "フ", "ァ", "ニ", "ミ", "ㄥ", "父", "Ξ", "δ", "░", "▒", "▓", "⌬", "⟁",
}

// SplashModel represents the splash screen state
type SplashModel struct {
	width         int
	height        int
	done          bool
	frame         int
	startTime     time.Time
	displayTime   time.Duration
	loadingText   string
	completedText string
	particles     []Particle
	loadingSteps  []string
	currentStep   int
	stepTimes     []time.Time
}

// NewSplashModel creates a new splash screen model
func NewSplashModel() SplashModel {
	// Prepare a more engaging loading sequence
	loadingSteps := []string{
		"Initializing system...",
		"Parsing node environment...",
		"Loading package data...",
		"Rendering interface...",
		"Optimizing performance...",
		"Ready!",
	}

	return SplashModel{
		frame:         0,
		startTime:     time.Now(),
		displayTime:   time.Second * 3, // Show splash longer
		loadingText:   loadingSteps[0],
		completedText: "Ready! Press any key to continue...",
		particles:     make([]Particle, 0),
		loadingSteps:  loadingSteps,
		currentStep:   0,
		stepTimes:     []time.Time{time.Now()},
	}
}

// Init initializes the splash screen
func (m SplashModel) Init() tea.Cmd {
	// Seed the random number generator for particles
	rand.Seed(time.Now().UnixNano())

	return tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// generateParticles creates background matrix-style effects
func (m *SplashModel) generateParticles() {
	// Only generate particles if we have enough space
	if m.width < 20 || m.height < 10 {
		return
	}

	// Add new particles occasionally
	if rand.Intn(10) < 3 && len(m.particles) < 100 {
		// Create a new particle
		p := Particle{
			x:    float64(rand.Intn(m.width)),
			y:    0,
			vx:   (rand.Float64() - 0.5) * 2,
			vy:   rand.Float64() * 1.5,
			char: matrixChars[rand.Intn(len(matrixChars))],
			color: lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
				100+rand.Intn(155),
				180+rand.Intn(75),
				80+rand.Intn(100))),
			lifespan: 20 + rand.Intn(60),
		}
		m.particles = append(m.particles, p)
	}

	// Update existing particles
	for i := 0; i < len(m.particles); i++ {
		p := &m.particles[i]
		p.x += p.vx
		p.y += p.vy
		p.lifespan--

		// Remove particles that are out of bounds or expired
		if p.x < 0 || p.x >= float64(m.width) || p.y >= float64(m.height) || p.lifespan <= 0 {
			m.particles = append(m.particles[:i], m.particles[i+1:]...)
			i--
		}
	}
}

// Update updates the splash screen
func (m SplashModel) Update(msg tea.Msg) (SplashModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Any key press skips the splash screen after a certain time
		if time.Since(m.startTime) > (m.displayTime / 2) {
			m.done = true
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		// Advance animation frame
		m.frame = (m.frame + 1) % len(loadingFrames)

		// Generate and update particles
		m.generateParticles()

		// Progress through loading steps
		elapsed := time.Since(m.startTime)
		stepDuration := m.displayTime / time.Duration(len(m.loadingSteps))
		currentStep := int(elapsed / stepDuration)

		if currentStep < len(m.loadingSteps) && currentStep > m.currentStep {
			m.currentStep = currentStep
			m.loadingText = m.loadingSteps[currentStep]
			m.stepTimes = append(m.stepTimes, time.Now())
		}

		// Check if we've displayed long enough
		if time.Since(m.startTime) > m.displayTime {
			m.done = true
			return m, nil
		}

		// Continue animation
		return m, tea.Tick(time.Millisecond*80, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// renderParticles draws the matrix-like effect in the background
func (m SplashModel) renderParticles() []string {
	// Create a grid for the particles
	grid := make([][]string, m.height)
	for y := range grid {
		grid[y] = make([]string, m.width)
		for x := range grid[y] {
			grid[y][x] = " "
		}
	}

	// Place particles on the grid
	for _, p := range m.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			grid[y][x] = lipgloss.NewStyle().
				Foreground(p.color).
				Render(p.char)
		}
	}

	// Convert grid to rows of text
	rows := make([]string, m.height)
	for y, row := range grid {
		rows[y] = strings.Join(row, "")
	}

	return rows
}

// View renders the splash screen
func (m SplashModel) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Get background particles if there's enough space
	background := ""
	if m.width >= 80 && m.height >= 20 {
		particleRows := m.renderParticles()
		background = strings.Join(particleRows, "\n")
	}

	// Style the ASCII logo with gradient colors
	logoLines := strings.Split(asciiLogo, "\n")
	coloredLogoLines := make([]string, len(logoLines))

	for i, line := range logoLines {
		// Create a gradient effect
		style := lipgloss.NewStyle().Bold(true)
		if i < len(logoLines)/2 {
			style = style.Foreground(terminalBrightBlue)
		} else {
			style = style.Foreground(terminalBrightCyan)
		}
		coloredLogoLines[i] = style.Render(line)
	}

	coloredLogo := strings.Join(coloredLogoLines, "\n")

	// Get progress indicator with dynamic color
	spinnerChar := loadingFrames[m.frame]

	// Determine text based on progress
	progress := float64(time.Since(m.startTime)) / float64(m.displayTime)

	// Generate the loading bar
	loadingWidth := 20
	filledWidth := int(progress * float64(loadingWidth))
	emptyWidth := loadingWidth - filledWidth

	loadingBar := "["
	loadingBar += strings.Repeat("=", filledWidth)
	if emptyWidth > 0 {
		loadingBar += ">"
		loadingBar += strings.Repeat(" ", emptyWidth-1)
	}
	loadingBar += "]"

	// Format the loading status
	var statusText string
	if m.currentStep == len(m.loadingSteps)-1 {
		statusText = lipgloss.NewStyle().
			Foreground(terminalBrightGreen).
			Bold(true).
			Render(m.completedText)
	} else {
		statusText = lipgloss.NewStyle().
			Foreground(terminalBrightYellow).
			Bold(true).
			Render(fmt.Sprintf("%s %s %s", m.loadingText, spinnerChar, loadingBar))
	}

	// Create version info display
	versionInfo := lipgloss.NewStyle().
		Foreground(terminalWhite).
		Render("TUI for Node.js, npm, and npx")

	// Add a border around the content
	mainContent := lipgloss.JoinVertical(
		lipgloss.Center,
		coloredLogo,
		"",
		statusText,
		"",
		versionInfo,
	)

	borderedContent := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(terminalBrightBlue).
		Padding(1, 2).
		Render(mainContent)

	// Center the content
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		borderedContent,
	)

	// If we have particles, overlay them
	if background != "" {
		// This is a simple way to overlay, but in a real implementation
		// you would need more sophisticated handling to properly blend
		// the particles with the centered content
		return centered
	}

	return centered
}

// IsDone indicates if the splash screen has completed
func (m SplashModel) IsDone() bool {
	return m.done
}
