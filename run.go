package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type runModel struct {
	quitting bool
	current  int
	total    int
	group    int
	status   string
	summary  string
	done     bool
	percent  float64
	bar      progress.Model
}

func newRunModel() runModel {
	bar := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(48),
	)

	return runModel{
		status: "Starting tuner...",
		bar:    bar,
	}
}

func (m runModel) Init() tea.Cmd {
	return nil
}

func (m runModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width := msg.Width - 24
		if width < 20 {
			width = 20
		}
		if width > 80 {
			width = 80
		}
		m.bar.Width = width
	case tea.KeyMsg:
		// Handle keyboard input for a Vite-style picker UX.
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case runProgressMsg:
		m.current = msg.Current
		m.total = msg.Total
		m.group = msg.Group
		if m.total > 0 {
			m.percent = float64(m.current) / float64(m.total)
		}
		if m.percent < 0 {
			m.percent = 0
		}
		if m.percent > 1 {
			m.percent = 1
		}
		if msg.Args != "" {
			m.status = fmt.Sprintf("Group %d: %s", msg.Group, msg.Args)
		}
		return m, nil
	case runStatusMsg:
		m.status = string(msg)
	case runSummaryMsg:
		m.summary = string(msg)
		m.status = "Tuning finished"
	case runDoneMsg:
		m.done = true
		m.percent = 1
		if m.total > 0 {
			m.current = m.total
		}
		return m, tea.Quit
	}
	return m, nil
}

func (m runModel) View() string {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39")).Render("Running tuner")
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(m.bar.ViewAs(m.percent))
	b.WriteString("\n")

	if m.total > 0 {
		percent := int(m.percent * 100)
		b.WriteString(fmt.Sprintf("%d/%d complete (%d%%)", m.current, m.total, percent))
		b.WriteString("\n")
	}

	b.WriteString(m.status)

	if !m.done {
		b.WriteString("\nPress q to quit")
	}

	return b.String()
}

type runProgressMsg struct {
	Current int
	Total   int
	Group   int
	Args    string
}

type runStatusMsg string
type runSummaryMsg string
type runDoneMsg struct{}
