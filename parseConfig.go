package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type tuningConfig struct {
	appPath string
	groups  []tuningGroup
}

type tuningGroup struct {
	inputs  []tuninginput
	results []tuningresult
}
type tuninginput struct {
	name  string
	start int
	end   int
}
type tuningresult struct {
	name        string
	isMaximizer bool
}
type ConfigJSON struct {
	App    string  `json:"app"`
	Groups []Group `json:"groups"`
}

type Group struct {
	Input  map[string][]int  `json:"input"`
	Result map[string]string `json:"result"`
}

//Opens the config file at the given path, parses it, and returns a tuningConfig struct. It also performs validation on the config structure and values, returning errors for any issues found.

func parseConfig(configPath string) (tuningConfig, error) {
	trimmedPath := strings.TrimSpace(configPath)
	if trimmedPath == "" {
		return tuningConfig{}, fmt.Errorf("config path is empty")
	}

	// Convert relative paths like ./config.json into an absolute path
	// based on the current working directory.
	absPath, err := filepath.Abs(trimmedPath)
	if err != nil {
		return tuningConfig{}, fmt.Errorf("failed to resolve config path %q: %w", trimmedPath, err)
	}
	absPath = filepath.Clean(absPath)

	rawJSONBytes, err := os.ReadFile(absPath)
	if err != nil {
		return tuningConfig{}, fmt.Errorf("error reading config file %q: %w", absPath, err)
	}

	var jsonObj ConfigJSON
	if err := json.Unmarshal(rawJSONBytes, &jsonObj); err != nil {
		return tuningConfig{}, fmt.Errorf("error parsing config file %q: %w", absPath, err)
	}
	config := tuningConfig{}
	config.appPath = jsonObj.App
	groups := jsonObj.Groups
	for groupIndex, group := range groups {
		parsedGroup := tuningGroup{
			inputs:  make([]tuninginput, 0, len(group.Input)),
			results: make([]tuningresult, 0, len(group.Result)),
		}

		// Sort keys so parsed output is deterministic regardless of map iteration order.
		inputNames := make([]string, 0, len(group.Input))
		for inputName := range group.Input {
			inputNames = append(inputNames, inputName)
		}
		sort.Strings(inputNames)

		for _, inputName := range inputNames {
			inputRange := group.Input[inputName]
			if len(inputRange) != 2 {
				return tuningConfig{}, fmt.Errorf("input range for %q must have exactly 2 integers", inputName)
			}
			if inputRange[0] > inputRange[1] {
				return tuningConfig{}, fmt.Errorf("input range for %q is invalid: start must be <= end", inputName)
			}

			parsedGroup.inputs = append(parsedGroup.inputs, tuninginput{
				name:  inputName,
				start: inputRange[0],
				end:   inputRange[1],
			})
		}

		resultNames := make([]string, 0, len(group.Result))
		for resultName := range group.Result {
			resultNames = append(resultNames, resultName)
		}
		sort.Strings(resultNames)

		for _, resultName := range resultNames {
			resultType := group.Result[resultName]
			isMaximizer := false
			switch resultType {
			case "maximizer":
				isMaximizer = true
			case "minimizer":
				isMaximizer = false
			default:
				return tuningConfig{}, fmt.Errorf("result type for %q must be either \"maximizer\" or \"minimizer\"", resultName)
			}

			parsedGroup.results = append(parsedGroup.results, tuningresult{
				name:        resultName,
				isMaximizer: isMaximizer,
			})
		}

		if len(parsedGroup.inputs) == 0 {
			return tuningConfig{}, fmt.Errorf("group %d must contain at least one input", groupIndex)
		}
		if len(parsedGroup.results) == 0 {
			return tuningConfig{}, fmt.Errorf("group %d must contain at least one result", groupIndex)
		}

		config.groups = append(config.groups, parsedGroup)
	}

	return config, nil
}
