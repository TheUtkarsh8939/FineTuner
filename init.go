package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type initModel struct {
	// cursor tracks the highlighted option in the list.
	cursor int
	// choices are the init templates the user can pick.
	choices []string
	// selected stores the currently chosen template index.
	selected int
	// done becomes true after Enter is pressed.
	done bool
	// quitting is set when the user exits without confirming.
	quitting bool
}

func newInitModel() initModel {
	return initModel{
		choices:  []string{"Minimal", "Demo"},
		selected: 0,
	}
}

func (m initModel) Init() tea.Cmd {
	return nil
}

func (m initModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle keyboard input for a Vite-style picker UX.
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			m.selected = m.cursor
		case "enter":
			m.selected = m.cursor
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}
func (m initModel) View() string {
	// Return empty view when the program is about to exit.
	if m.done || m.quitting {
		return ""
	}
	var b strings.Builder
	color := lipgloss.NewStyle().Foreground(lipgloss.Color("#51ff00")).Bold(true)
	b.WriteString(color.Render("Please choose a config type:"))

	for i, choice := range m.choices {
		// Show a caret for the current row and an x for the selected row.
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		selected := " "
		if i == m.selected {
			selected = "x"
		}
		var color lipgloss.Style
		switch choice {
		case "Minimal":
			color = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
		case "Demo":
			color = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
		default:
			color = lipgloss.NewStyle()
		}

		fmt.Fprintf(&b, "\n%s [%s] %s", cursor, selected, color.Render(choice))
	}

	b.WriteString("\nUse Space/Enter to select, and up/down to navigate:\n\n")
	return b.String()
}
func runInit() error {
	p := tea.NewProgram(newInitModel())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	m, ok := finalModel.(initModel)
	if !ok {
		return fmt.Errorf("unexpected init model type")
	}
	// If the user quit (q/Ctrl+C), do not create a config file.
	if m.quitting {
		return nil
	}

	// Persist the chosen template to config.json.
	chosen := m.choices[m.selected]
	if err := os.WriteFile("config.json", []byte(configFor(chosen)), 0644); err != nil {
		return err
	}

	color := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ffff"))
	fmt.Println(color.Render(fmt.Sprintf("Initialised %s config", chosen)))
	return nil
}

func configFor(configType string) string {
	// Return minimal starter templates for each supported profile.
	switch configType {
	case "Demo":
		return `
{
    "app":"demo.exe",
    "groups":[
        {
            "input":{
                "group1var1":[1,10],
                "group1var2":[5,15]
            },
            "result":{
                "group1result1":"minimizer",
                "group1result2":"maximizer"
            }
        },
        {
            "input":{
                "group2var1":[1,10],
                "group2var2":[5,15]
            },
            "result":{
                "group2result1":"minimizer",
                "group2result2":"maximizer"
            }
        }
    ]
}
		`
	case "Minimal":
		fallthrough
	default:
		return `
{
    "app":"",
    "groups":[]
}
		`
	}
}
