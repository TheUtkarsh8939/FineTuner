# Finetuner

Finetuner is a CLI tuning tool that runs a target program across all combinations of configured input ranges, evaluates configured result metrics, and returns the best input set per group.

The run UI uses Bubble Tea and shows a gradient progress bar with live status updates.

## Features

- Interactive `init` command to generate starter config files.
- Exhaustive search across each group's input ranges.
- Per-result optimization rules:
  - `maximizer`: larger value is better
  - `minimizer`: smaller value is better
- Bubble Tea run UI with progress bar and automatic exit when tuning completes.
- Final JSON summary of best inputs and results per group.

## Requirements

- Go 1.24.2 or newer

## Quick Start

### 1) Build the bundled demo target (optional)

```powershell
go build -o demo/demo.exe ./demo
```

### 2) Create or edit a config

Generate a starter config:

```powershell
go run . init
```

Or create `config.json` manually.

### 3) Run tuning

```powershell
go run . run ./config.json
```

### 4) Build the tool binary

```powershell
go build -o finetuner.exe .
```

## Commands

- `finetuner init`
  - Opens an interactive picker and writes `config.json`.
- `finetuner run <configPath>`
  - Runs the tuner using the given config file path.
- `finetuner help`
  - Prints command help.

## Config Format

```json
{
  "app": "demo/demo.exe",
  "groups": [
    {
      "input": {
        "group1var1": [1, 10],
        "group1var2": [5, 15]
      },
      "result": {
        "group1result1": "minimizer",
        "group1result2": "maximizer"
      }
    }
  ]
}
```

### Fields

- `app`: path to the program being tuned.
- `groups`: list of independent tuning groups.
- `group.input`: map of input name to inclusive range `[start, end]`.
- `group.result`: map of result key to optimization mode (`maximizer` or `minimizer`).

## How Tuning Works

For each group, Finetuner:

1. Generates the Cartesian product of all input ranges.
2. Invokes the target program with flags in `--name value` form.
3. Parses JSON output from the target program.
4. Extracts only result keys declared in that group.
5. Compares candidates using the configured maximize/minimize rules.
6. Tracks and returns the best input/result combination for the group.

Progress is sent to the UI on every program invocation.

## Target Program Contract

The tuned application is expected to:

- Accept input values from CLI flags.
- Print a JSON object to stdout.
- Include all result keys required by the current group.
- Emit numeric values for required result keys.

Example invocation from the tuner:

```text
demo.exe --group2var1 10 --group2var2 15
```

Example output from the tuned app:

```json
{
  "group1result1": 25,
  "group1result2": 75,
  "group2result1": 50,
  "group2result2": 50
}
```

## Output Summary

At completion, Finetuner prints JSON like:

```json
[
  {
    "groupIndex": 1,
    "inputs": {
      "group1var1": 1,
      "group1var2": 5
    },
    "results": {
      "group1result1": 6,
      "group1result2": 94
    }
  }
]
```

## Notes

- Relative config paths are resolved from the current working directory.
- Ensure `app` points to the actual executable location.
  - Example: if using the bundled demo binary from project root, use `demo/demo.exe`.
