package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ASCII art for quit screen
const quitArt = `
 ████████╗██╗  ██╗ █████╗ ███╗   ██╗██╗  ██╗███████╗    ███████╗ ██████╗ ██████╗ 
 ╚══██╔══╝██║  ██║██╔══██╗████╗  ██║██║ ██╔╝██╔════╝    ██╔════╝██╔═══██╗██╔══██╗
    ██║   ███████║███████║██╔██╗ ██║█████╔╝ ███████╗    █████╗  ██║   ██║██████╔╝
    ██║   ██╔══██║██╔══██║██║╚██╗██║██╔═██╗ ╚════██║    ██╔══╝  ██║   ██║██╔══██╗
    ██║   ██║  ██║██║  ██║██║ ╚████║██║  ██╗███████║    ██║     ╚██████╔╝██║  ██║
    ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝  ╚═╝╚══════╝    ╚═╝      ╚═════╝ ╚═╝  ╚═╝
                                                                               
                        ██╗   ██╗███████╗██╗███╗   ██╗ ██████╗                  
                        ██║   ██║██╔════╝██║████╗  ██║██╔════╝                  
                        ██║   ██║███████╗██║██╔██╗ ██║██║  ███╗                 
                        ██║   ██║╚════██║██║██║╚██╗██║██║   ██║                 
                        ╚██████╔╝███████║██║██║ ╚████║╚██████╔╝                 
                         ╚═════╝ ╚══════╝╚═╝╚═╝  ╚═══╝ ╚═════╝                  
                                                                               
                    ██╗      █████╗ ███████╗██╗   ██╗███╗   ██╗ ██████╗ ██████╗ ███████╗
                    ██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝████╗  ██║██╔═══██╗██╔══██╗██╔════╝
                    ██║     ███████║  ███╔╝  ╚████╔╝ ██╔██╗ ██║██║   ██║██║  ██║█████╗  
                    ██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║╚██╗██║██║   ██║██║  ██║██╔══╝  
                    ███████╗██║  ██║███████╗   ██║   ██║ ╚████║╚██████╔╝██████╔╝███████╗
                    ╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝
`

// FadeEffect represents the fade direction
type FadeEffect int

const (
	FadeIn FadeEffect = iota
	FadeOut
)

// QuitModel represents the quit screen state
type QuitModel struct {
	width       int
	height      int
	done        bool
	frame       int
	startTime   time.Time
	displayTime time.Duration
	fadeEffect  FadeEffect
	fadeSteps   int
	currentStep int
	messages    []string
	particles   []Particle
}

// NewQuitModel creates a new quit screen model
func NewQuitModel() QuitModel {
	farewell := []string{
		"Thanks for using LazyNode!",
		"See you soon!",
		"Have a great day!",
		"Happy coding!",
		"Node.js projects made easier!",
		"Until next time!",
		"Keep building awesome things!",
	}

	return QuitModel{
		frame:       0,
		startTime:   time.Now(),
		displayTime: time.Second * 3, // Show quit screen for 3 seconds
		fadeEffect:  FadeIn,
		fadeSteps:   15,
		messages:    farewell,
		particles:   make([]Particle, 0),
	}
}

// generateParticles creates sparkling particles for the quit screen
func (m *QuitModel) generateParticles() {
	// Only generate particles if we have enough space
	if m.width < 20 || m.height < 10 {
		return
	}

	// Add new particles occasionally
	if m.frame%2 == 0 && len(m.particles) < 100 {
		// Create new particle at random position
		x := float64(m.width/4 + rand.Intn(m.width/2))
		y := float64(m.height/4 + rand.Intn(m.height/2))

		// Random movement and appearance
		vx := (rand.Float64() - 0.5) * 2
		vy := (rand.Float64() - 0.5) * 2

		// Sparkle characters
		chars := []string{"✨", "•", "⋆", "✦", "✳", "✵", "*", "✺", "✹"}
		char := chars[rand.Intn(len(chars))]

		// Create glowing color with alpha based on fade effect
		r, g, b := 220+rand.Intn(35), 220+rand.Intn(35), 160+rand.Intn(95)
		color := lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))

		p := Particle{
			x:        x,
			y:        y,
			vx:       vx,
			vy:       vy,
			char:     char,
			color:    color,
			lifespan: 20 + rand.Intn(30),
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
		if p.x < 0 || p.x >= float64(m.width) || p.y < 0 || p.y >= float64(m.height) || p.lifespan <= 0 {
			m.particles = append(m.particles[:i], m.particles[i+1:]...)
			i--
		}
	}
}

// Init initializes the quit screen
func (m QuitModel) Init() tea.Cmd {
	// Seed the random number generator for particles
	rand.Seed(time.Now().UnixNano())

	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update updates the quit screen
func (m QuitModel) Update(msg tea.Msg) (QuitModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		// Update frame counter
		m.frame++

		// Generate and update particles
		m.generateParticles()

		// Update fade effect
		if m.fadeEffect == FadeIn && m.currentStep < m.fadeSteps {
			m.currentStep++
		} else if m.fadeEffect == FadeIn && m.currentStep >= m.fadeSteps {
			// Switch to fade out after showing for a while
			if time.Since(m.startTime) > m.displayTime {
				m.fadeEffect = FadeOut
			}
		} else if m.fadeEffect == FadeOut {
			m.currentStep--
			if m.currentStep <= 0 {
				m.done = true
				return m, tea.Quit
			}
		}

		// Continue animation
		return m, tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}

	return m, nil
}

// View renders the quit screen
func (m QuitModel) View() string {
	// If window size is not set, show a simple but colorful goodbye message
	if m.width < 80 || m.height < 20 {
		// Create a simple colorful farewell
		selectedMessage := m.messages[int(time.Since(m.startTime).Seconds())%len(m.messages)]

		// Simple goodbye with progress animation
		spinnerChars := "⣾⣽⣻⢿⡿⣟⣯⣷"
		spinnerChar := string(spinnerChars[m.frame%8])

		// Calculate fade alpha
		var alpha float64
		if m.fadeEffect == FadeIn {
			alpha = float64(m.currentStep) / float64(m.fadeSteps)
		} else {
			alpha = float64(m.currentStep) / float64(m.fadeSteps)
		}

		// Simple but colorful goodbye box
		box := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
				int(150*alpha), int(180*alpha), int(210*alpha)))).
			Padding(1, 3).
			Align(lipgloss.Center).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Center,
					lipgloss.NewStyle().
						Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
							int(230*alpha), int(210*alpha), int(120*alpha)))).
						Bold(true).
						Render("Goodbye from LazyNode!"),
					"",
					lipgloss.NewStyle().
						Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
							int(200*alpha), int(220*alpha), int(200*alpha)))).
						Render(selectedMessage),
					"",
					lipgloss.NewStyle().
						Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
							int(150*alpha), int(150*alpha), int(220*alpha)))).
						Render(spinnerChar),
				),
			)

		return box
	}

	// Calculate fade alpha based on current step
	var alpha float64
	if m.fadeEffect == FadeIn {
		alpha = float64(m.currentStep) / float64(m.fadeSteps)
	} else {
		alpha = float64(m.currentStep) / float64(m.fadeSteps)
	}

	// Get a random farewell message that changes periodically
	selectedMessage := m.messages[int(time.Since(m.startTime).Seconds()*0.5)%len(m.messages)]

	// Style the ASCII art with gradient colors and fading effect
	artLines := strings.Split(quitArt, "\n")
	coloredArtLines := make([]string, 0, len(artLines))

	for i, line := range artLines {
		if len(line) == 0 {
			coloredArtLines = append(coloredArtLines, "")
			continue
		}

		// Create a gradient effect with fading
		opacity := int(240 * alpha)
		var style lipgloss.Style

		// Different colors for different parts of the text
		if i < 6 {
			// "THANKS FOR"
			color := fmt.Sprintf("#%02x%02x%02x", opacity, opacity/2, opacity/4)
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
		} else if i < 12 {
			// "USING"
			color := fmt.Sprintf("#%02x%02x%02x", opacity/4, opacity, opacity/2)
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
		} else {
			// "LAZYNODE"
			color := fmt.Sprintf("#%02x%02x%02x", opacity/4, opacity/2, opacity)
			style = lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
		}

		coloredArtLines = append(coloredArtLines, style.Render(line))
	}

	coloredArt := strings.Join(coloredArtLines, "\n")

	// Style the farewell message with shimmer effect
	shimmerPhase := m.frame % 20
	shimmerPos := shimmerPhase * len(selectedMessage) / 20

	var farewellParts []string
	for i, char := range selectedMessage {
		brightness := 220.0
		// Add brightness to characters near the shimmer position
		distFromShimmer := abs(i - shimmerPos)
		if distFromShimmer < 5 {
			brightness = 255.0 - float64(distFromShimmer)*10.0
		}

		charStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
				int(brightness*alpha),
				int(brightness*alpha),
				int((brightness+20)*alpha)))).
			Bold(true).
			Italic(true)

		farewellParts = append(farewellParts, charStyle.Render(string(char)))
	}

	farewell := strings.Join(farewellParts, "")

	// Center the content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		coloredArt,
		"",
		farewell,
	)

	// Place particles in a 2D grid
	particleGrid := make([][]string, m.height)
	for y := range particleGrid {
		particleGrid[y] = make([]string, m.width)
		for x := range particleGrid[y] {
			particleGrid[y][x] = " "
		}
	}

	// Add particles to the grid
	for _, p := range m.particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < m.width && y >= 0 && y < m.height {
			particleGrid[y][x] = lipgloss.NewStyle().
				Foreground(p.color).
				Render(p.char)
		}
	}

	// Render the content in the center
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	// For simplicity in terminal rendering, we'll just show the centered content with particles
	return centered
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IsDone indicates if the quit screen animation has completed
func (m QuitModel) IsDone() bool {
	return m.done
}
