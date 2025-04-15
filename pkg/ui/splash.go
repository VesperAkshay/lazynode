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
func (m SplashModel) renderParticles() {
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

	// Update the particles state
	m.particles = m.particles[:0]
	for _, row := range rows {
		for _, char := range row {
			if char != ' ' {
				m.particles = append(m.particles, Particle{
					x:    float64(len(row) / 2),
					y:    float64(len(rows) / 2),
					vx:   0,
					vy:   0,
					char: string(char),
					color: lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
						100+rand.Intn(155),
						180+rand.Intn(75),
						80+rand.Intn(100))),
					lifespan: 20 + rand.Intn(60),
				})
			}
		}
	}
}

// View renders the splash screen
func (m SplashModel) View() string {
	if m.width < 20 || m.height < 10 {
		return "Loading LazyNode..."
	}

	// Call renderParticles but don't use the return value for now
	// This still updates the state of particles
	m.renderParticles()

	// Style the ASCII art logo with gradient colors
	var logoLines []string
	asciiLines := strings.Split(strings.TrimSpace(asciiLogo), "\n")

	// Create a gradient effect for the logo
	for i, line := range asciiLines {
		// Calculate gradient color based on line number
		r := 180 + (i*10)%75
		g := 100 + (i*15)%155
		b := 220 + (i*5)%35

		styledLine := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))).
			Bold(true).
			Render(line)

		logoLines = append(logoLines, styledLine)
	}

	// Add tagline below the logo with glowing effect
	tagline := lipgloss.NewStyle().
		Foreground(terminalBrightCyan).
		Italic(true).
		Bold(true).
		Render("A powerful Terminal UI for managing Node.js projects with style")

	// Build the loading animation
	loadingAnimation := loadingFrames[m.frame]

	// Calculate progress through the splash screen
	progress := float64(time.Since(m.startTime)) / float64(m.displayTime)
	if progress > 1.0 {
		progress = 1.0
	}

	// Create a progress bar
	progressBarWidth := 40
	completedWidth := int(float64(progressBarWidth) * progress)
	remainingWidth := progressBarWidth - completedWidth

	progressBar := lipgloss.NewStyle().Foreground(terminalBrightGreen).Render(strings.Repeat("█", completedWidth))
	progressBar += lipgloss.NewStyle().Foreground(terminalBrightBlack).Render(strings.Repeat("█", remainingWidth))

	// Create fancy progress indicator with current step
	stepIndicator := ""
	for i := 0; i < len(m.loadingSteps); i++ {
		if i < m.currentStep {
			// Completed step
			stepIndicator += lipgloss.NewStyle().
				Foreground(terminalBrightGreen).
				Render("● ")
		} else if i == m.currentStep {
			// Current step
			stepIndicator += lipgloss.NewStyle().
				Foreground(terminalBrightYellow).
				Render(loadingAnimation + " ")
		} else {
			// Future step
			stepIndicator += lipgloss.NewStyle().
				Foreground(terminalBrightBlack).
				Render("○ ")
		}
	}

	// Create loading message with terminal style
	loadingMsg := lipgloss.NewStyle().
		Foreground(terminalBrightWhite).
		Bold(true).
		Render(m.loadingText)

	// Combine all elements with proper spacing and centering
	logoStr := strings.Join(logoLines, "\n")

	// Version info with monospace styling
	versionInfo := lipgloss.NewStyle().
		Foreground(terminalBrightBlack).
		Render("v1.0.0")

	// Center everything based on the terminal width
	logoWidth := len(asciiLines[0])
	leftPadding := (m.width - logoWidth) / 2
	if leftPadding < 0 {
		leftPadding = 0
	}

	// Apply padding to center the content
	paddedLogo := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(logoStr)
	paddedTagline := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(tagline)
	paddedProgressBar := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(progressBar)
	paddedStepIndicator := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(stepIndicator)
	paddedLoadingMsg := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(loadingMsg)
	paddedVersionInfo := lipgloss.NewStyle().PaddingLeft(leftPadding).Render(versionInfo)

	// Combine all components
	result := "\n\n" + paddedLogo + "\n"
	result += paddedTagline + "\n\n"
	result += paddedProgressBar + "\n"
	result += paddedStepIndicator + "\n"
	result += paddedLoadingMsg + "\n\n"
	result += paddedVersionInfo

	// Add skip message if we've shown enough
	if time.Since(m.startTime) > (m.displayTime / 2) {
		skipMsg := lipgloss.NewStyle().
			Foreground(terminalBrightWhite).
			Italic(true).
			Render("Press any key to skip...")

		result += "\n\n" + lipgloss.NewStyle().PaddingLeft(leftPadding).Render(skipMsg)
	}

	return result
}

// IsDone indicates if the splash screen has completed
func (m SplashModel) IsDone() bool {
	return m.done
}
