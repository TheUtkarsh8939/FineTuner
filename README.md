# Finetuner

**DEMO VIDEO:** [Demo Video](https://youtu.be/qqJxVPGtGPI)
I built this CLI tool to help me find the best parameters for programs by running them across a bunch of different input combinations. It tests every possibility, tracks which one works best, and gives you back the winner. The UI has a nice progress bar and shows live updates while it's running.

## Quick Setup

Just grab the precompiled `Finetuner.exe` from the releases and you're good to go. No need to build anything.

## Usage

### 1) Grab a config

Either generate a new one:

```powershell
./Finetuner.exe init
```

Or make your own `config.json` manually.

### 2) Run it

```powershell
./Finetuner.exe run ./config.json
```

Done. It'll spit out the best inputs and results when it finishes.

## Config

Here's what your `config.json` should look like:

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

- `app`: path to your program
- `groups`: different testing groups (they run independently)
- `group.input`: variables and their min/max ranges
- `group.result`: what you want to measure, and whether you want bigger numbers (`maximizer`) or smaller numbers (`minimizer`)

## How It Works

For each group, here's what happens:

1. I generate every possible combo of your input ranges
2. Run your program with each combo
3. Parse the JSON output
4. Compare the results and figure out which combo was best
5. Tell you the winner

Pretty straightforward stuff.

## Your Program

Your program just needs to:

- Accept inputs as command-line flags (like `--myvar 42`)
- Print a JSON object to stdout with your results
- Include all the result keys you defined in the config

Example of how I'll call your program:

```
demo.exe --group2var1 10 --group2var2 15
```

And it should output something like:

```json
{
  "group1result1": 25,
  "group1result2": 75,
  "group2result1": 50,
  "group2result2": 50
}
```

## Output

When everything's done, you get back a JSON file with the best results for each group:

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

That's it. Use those inputs and you're golden.

## Notes

- Config paths are relative to where you run the command from
- Make sure your `app` path points to the right place (e.g., `demo/demo.exe` if you're running from the project root)
