package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// The first argument selects the command: init, run, or help.
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	cmd := os.Args[1]

	switch cmd {
	case "init":
		// init launches the interactive picker and writes config.json.
		if err := runInit(); err != nil {
			fmt.Fprintf(os.Stderr, "init error: %v\n", err)
			os.Exit(1)
		}
	case "run":
		// run expects a config path argument.
		if err := runRun(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "run error: %v\n", err)
			os.Exit(1)
		}
	case "help":
		printHelp()
	default:

		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
}

func runRun(args []string) error {
	if len(args) < 1 {
		color := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
		fmt.Println(color.Render("Error: missing config path argument"))
		fmt.Println("Use finetuner help for more info")
		return nil
	}

	// The second CLI argument after "run" is the config path.
	configPath := args[0]
	p := tea.NewProgram(newRunModel())
	if config, err := parseConfig(configPath); err != nil {
		return err
	} else {
		go func(cfg tuningConfig) {
			bestByGroup, runErr := tunerRun(cfg, func(output string) {
				if progress, ok := parseProgressUpdate(output); ok {
					p.Send(progress)
					return
				}
				p.Send(runStatusMsg(output))
			})
			if runErr != nil {
				p.Send(runStatusMsg(fmt.Sprintf("Tuner failed: %v", runErr)))
				p.Send(runDoneMsg{})
				return
			}

			summaryJSON, marshalErr := json.MarshalIndent(bestByGroup, "", "  ")
			if marshalErr != nil {
				p.Send(runStatusMsg(fmt.Sprintf("Failed to encode best results: %v", marshalErr)))
				p.Send(runDoneMsg{})
				return
			}

			p.Send(runSummaryMsg("Best results by group:\n" + string(summaryJSON)))
			p.Send(runDoneMsg{})
		}(config)
	}
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	m, ok := finalModel.(runModel)
	if !ok {
		return fmt.Errorf("unexpected run model type")
	}
	if m.quitting {
		return nil
	}

	if m.summary != "" {
		fmt.Println(m.summary)
	}
	return nil
}

func parseProgressUpdate(line string) (runProgressMsg, bool) {
	if !strings.HasPrefix(line, "Progress ") {
		return runProgressMsg{}, false
	}

	var current int
	var total int
	var group int
	if _, err := fmt.Sscanf(line, "Progress %d/%d: group %d ->", &current, &total, &group); err != nil {
		return runProgressMsg{}, false
	}

	args := ""
	if idx := strings.Index(line, "->"); idx >= 0 {
		args = strings.TrimSpace(line[idx+2:])
	}

	return runProgressMsg{
		Current: current,
		Total:   total,
		Group:   group,
		Args:    args,
	}, true
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  finetuner <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  init    interactive config setup (basic/demo)")
	fmt.Println("  run     runs tuner with config file path")
	fmt.Println("  help    shows this help")
}
