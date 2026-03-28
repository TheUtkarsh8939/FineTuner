package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type groupBestResult struct {
	GroupIndex int                `json:"groupIndex"`
	Inputs     map[string]int     `json:"inputs"`
	Results    map[string]float64 `json:"results"`
}

// tunerRun evaluates every input combination in each group and returns the best
// input/result set per group based on maximize/minimize preferences.
func tunerRun(config tuningConfig, sendOutput func(string)) ([]groupBestResult, error) {
	appPath, err := cleanPath(config.appPath)
	if err != nil {
		return nil, err
	}
	if len(config.groups) == 0 {
		return nil, fmt.Errorf("config has no groups to tune")
	}

	totalRuns := totalRunCount(config.groups)
	runsDone := 0
	bestByGroup := make([]groupBestResult, 0, len(config.groups))

	for groupIndex, group := range config.groups {
		currentInputs := make(map[string]int, len(group.inputs))
		bestInputs := make(map[string]int)
		bestResults := make(map[string]float64)
		foundBest := false

		var evaluate func(inputIndex int) error
		evaluate = func(inputIndex int) error {
			if inputIndex == len(group.inputs) {
				runsDone++
				args := buildArgs(group.inputs, currentInputs)
				sendOutput(fmt.Sprintf("Progress %d/%d: group %d -> %s", runsDone, totalRuns, groupIndex+1, strings.Join(args, " ")))

				extractedResults, execErr := executeAndExtract(appPath, args, group.results)
				if execErr != nil {
					sendOutput(fmt.Sprintf("Run failed for group %d (%s): %v", groupIndex+1, strings.Join(args, " "), execErr))
					return nil
				}

				if !foundBest || isBetter(extractedResults, bestResults, group.results) {
					bestInputs = cloneInputs(currentInputs)
					bestResults = cloneResults(extractedResults)
					foundBest = true

					bestResultsJSON, marshalErr := json.Marshal(bestResults)
					if marshalErr == nil {
						sendOutput(fmt.Sprintf("New best for group %d: inputs=%s results=%s", groupIndex+1, formatInputs(bestInputs), string(bestResultsJSON)))
					}
				}

				return nil
			}

			input := group.inputs[inputIndex]
			for value := input.start; value <= input.end; value++ {
				currentInputs[input.name] = value
				if err := evaluate(inputIndex + 1); err != nil {
					return err
				}
			}

			return nil
		}

		if err := evaluate(0); err != nil {
			return nil, err
		}
		if !foundBest {
			return nil, fmt.Errorf("no successful runs for group %d", groupIndex+1)
		}

		bestByGroup = append(bestByGroup, groupBestResult{
			GroupIndex: groupIndex + 1,
			Inputs:     bestInputs,
			Results:    bestResults,
		})
	}

	return bestByGroup, nil
}

func totalRunCount(groups []tuningGroup) int {
	total := 0
	for _, group := range groups {
		count := 1
		for _, input := range group.inputs {
			count *= (input.end - input.start + 1)
		}
		total += count
	}
	if total == 0 {
		return 1
	}
	return total
}

func buildArgs(inputs []tuninginput, values map[string]int) []string {
	args := make([]string, 0, len(inputs)*2)
	for _, input := range inputs {
		args = append(args, fmt.Sprintf("--%s", input.name), strconv.Itoa(values[input.name]))
	}
	return args
}

func executeAndExtract(appPath string, args []string, expected []tuningresult) (map[string]float64, error) {
	cmd := exec.Command(appPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		trimmed := strings.TrimSpace(string(output))
		if trimmed == "" {
			return nil, fmt.Errorf("exec error: %w", err)
		}
		return nil, fmt.Errorf("exec error: %w (output: %s)", err, trimmed)
	}

	var payload map[string]json.RawMessage
	if err := json.Unmarshal(output, &payload); err != nil {
		return nil, fmt.Errorf("invalid JSON output: %w", err)
	}

	missing := make([]string, 0)
	extracted := make(map[string]float64, len(expected))
	for _, result := range expected {
		raw, ok := payload[result.name]
		if !ok {
			missing = append(missing, result.name)
			continue
		}

		var numeric float64
		if err := json.Unmarshal(raw, &numeric); err != nil {
			return nil, fmt.Errorf("result %q must be numeric", result.name)
		}
		extracted[result.name] = numeric
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing result keys %v", missing)
	}

	return extracted, nil
}

func isBetter(candidate map[string]float64, currentBest map[string]float64, priorities []tuningresult) bool {
	const epsilon = 1e-9

	for _, result := range priorities {
		candidateValue := candidate[result.name]
		bestValue := currentBest[result.name]

		if result.isMaximizer {
			if candidateValue > bestValue+epsilon {
				return true
			}
			if candidateValue < bestValue-epsilon {
				return false
			}
		} else {
			if candidateValue < bestValue-epsilon {
				return true
			}
			if candidateValue > bestValue+epsilon {
				return false
			}
		}
	}

	return false
}

func cloneInputs(values map[string]int) map[string]int {
	cloned := make(map[string]int, len(values))
	for k, v := range values {
		cloned[k] = v
	}
	return cloned
}

func cloneResults(values map[string]float64) map[string]float64 {
	cloned := make(map[string]float64, len(values))
	for k, v := range values {
		cloned[k] = v
	}
	return cloned
}

func formatInputs(values map[string]int) string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", key, values[key]))
	}

	return strings.Join(parts, ", ")
}
