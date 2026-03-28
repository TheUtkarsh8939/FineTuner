package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	values, err := parseFlagValues(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Print in sorted order for stable output.
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	total := 0
	for _, key := range keys {
		value := values[key]
		total += value
	}
	fmt.Printf(`
{
	"group1result1":%d,
	"group1result2":%d,
	"group2result1":%d,
	"group2result2":%d
}
	`, total, 100-total, total*2, 100-(total*2))
}

// parseFlagValues converts CLI args into a map where each flag name maps
// to its integer value (e.g. --group2var1 9 or --group2var1=9).
func parseFlagValues(args []string) (map[string]int, error) {
	values := make(map[string]int)

	for i := 0; i < len(args); i++ {
		token := args[i]
		if !strings.HasPrefix(token, "--") {
			return nil, fmt.Errorf("expected flag at argument %d, got %q", i+1, token)
		}

		raw := strings.TrimPrefix(token, "--")
		if raw == "" {
			return nil, fmt.Errorf("empty flag name at argument %d", i+1)
		}

		// Support --name=value form.
		if strings.Contains(raw, "=") {
			parts := strings.SplitN(raw, "=", 2)
			if parts[0] == "" {
				return nil, fmt.Errorf("empty flag name in %q", token)
			}
			value, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("flag %q expects an integer value", parts[0])
			}
			values[parts[0]] = value
			continue
		}

		// Handle --name value form.
		if i+1 >= len(args) {
			return nil, fmt.Errorf("flag %q is missing a value", raw)
		}

		next := args[i+1]
		if strings.HasPrefix(next, "--") {
			return nil, fmt.Errorf("flag %q is missing a numeric value", raw)
		}

		value, err := strconv.Atoi(next)
		if err != nil {
			return nil, fmt.Errorf("flag %q expects an integer value", raw)
		}

		values[raw] = value
		i++
	}

	return values, nil
}
