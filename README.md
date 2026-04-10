# Finetuner

**DEMO VIDEO:** [Demo Video](https://youtu.be/qqJxVPGtGPI)

Finetuner is a really simple utility to help find the best input configurations for a program. You define ranges for your input variables and what results you want to optimize for, and Finetuner runs through all the combos to find the best one.

## Why?

I was working on my other project *Barracuda* which is a chess engine, and I wanted a quick way to tune the parameters for my search heurastics and evaluation function. I didn't want to set up a whole machine learning pipeline or anything, just a simple brute-force search over a few variables. So I built Finetuner to do exactly that.

## Where?

Finetuner works wherewever you have a CLI program and the you can score the performance of that program with some numeric results. It's not limited to chess engines or anything, you can use it for any kind of optimization problem where you can define input variables and measurable results.

## How?
Finetuner starts loads up the config files and starts the program in config file as a subprocess for each combination of input variables. It captures the output, parses the JSON results, and keeps track of which input combo gives the best results based on your criteria (minimizer or maximizer). Once it's done, it outputs the best inputs and their corresponding results.

The input and results are defined in groups, so you can have multiple sets of variables and results that are optimized independently. This is useful if you want to tune different aspects of your program separately.

The input is passed via CLI flags and the program must output a JSON object with the results. Finetuner will handle all the combinations and comparisons for you.

Example of how It'll call your program:
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

## Quick Setup

Just grab the precompiled `Finetuner.exe` from the releases and you're good to go. No need to build anything.

Then:
- Run `./finetuner.exe init" choose minimal config and it will create a "config.json" file for you(Note it will rewrite any existing config file, so make a backup if you have one)
- Edit the `config.json` to define your input variables, result variables, and optimization criteria
- Run `./finetuner.exe run {path to config.json}` to start the optimization process

## Demo Setup

Download the release and run the following commands in the project root:

- `./finetuner.exe init` (choose demo config)
- `./finetuner.exe run config.json`

## Notes

- Config paths are relative to where you run the command from
- Make sure your `app` path points to the right place (e.g., `demo/demo.exe` if you're running from the project root)
