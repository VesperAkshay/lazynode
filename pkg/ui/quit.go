package ui

import (
	"fmt"
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
}

// NewQuitModel creates a new quit screen model
func NewQuitModel() QuitModel {
	farewell := []string{
		"Thanks for using LazyNode!",
		"See you soon!",
		"Have a great day!",
	}

	return QuitModel{
		frame:       0,
		startTime:   time.Now(),
		displayTime: time.Second * 3, // Show quit screen for 3 seconds
		fadeEffect:  FadeIn,
		fadeSteps:   10,
		messages:    farewell,
	}
}

// Init initializes the quit screen
func (m QuitModel) Init() tea.Cmd {
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

	// Get a random farewell message
	selectedMessage := m.messages[int(time.Since(m.startTime).Seconds())%len(m.messages)]

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

	// Style the farewell message
	farewellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x",
			int(220*alpha),
			int(220*alpha),
			int(220*alpha)))).
		Bold(true).
		Italic(true)

	farewell := farewellStyle.Render(selectedMessage)

	// Center the content
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			coloredArt,
			"",
			farewell,
		),
	)

	return centered
}

// IsDone indicates if the quit screen animation has completed
func (m QuitModel) IsDone() bool {
	return m.done
}
