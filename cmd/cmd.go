/*
Axion CLI Calculator - Cobra Command Structure
===============================================
Author: Uthman
Year: 2025

This file implements the Cobra-based command structure for Axion calculator.
The root command launches the interactive REPL, while subcommands provide
direct access to specific features (conversion, history, etc.).
*/

package cmd

import (
	"Axion/constants"
	"Axion/evaluator"
	"Axion/history"
	"Axion/parser"
	"Axion/settings"
	"Axion/tokenizer"
	"Axion/units"
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const banner = `
  ╔═╗─┐ ┬┬┌─┐┌┐┌
  ╠═╣┌┴┬┘││ ││││
  ╩ ╩┴ └─┴└─┘┘└┘
`

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

var rootCmd = &cobra.Command{
	Use:   "axion",
	Short: "Axion - A powerful CLI calculator",
	Long: colorCyan + banner + colorReset + `
` + colorBold + `Axion` + colorReset + ` is a feature-rich command-line calculator supporting:
  ` + colorGreen + `✓` + colorReset + ` Mathematical expressions with variables
  ` + colorGreen + `✓` + colorReset + ` Unit conversions across multiple categories
  ` + colorGreen + `✓` + colorReset + ` Built-in mathematical functions and constants
  ` + colorGreen + `✓` + colorReset + ` Calculation history and session management
  ` + colorGreen + `✓` + colorReset + ` Customizable precision and settings`,
	Run: startREPL,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize constants system
	err := constants.Load("constants.json")
	if err != nil {
		fmt.Printf(colorYellow+"Warning: Failed to load constants: %v\n"+colorReset, err)
	}
}

// startREPL launches the interactive calculator session
func startREPL(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	printWelcome()

	for {
		fmt.Print(colorCyan + "» " + colorReset)

		if !scanner.Scan() {
			fmt.Println(colorYellow + "\nGoodbye!" + colorReset)
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		switch {
		case input == "exit" || input == "quit":
			fmt.Println(colorYellow + "Goodbye!" + colorReset)
			return

		case input == "clear" || input == "cls":
			clearScreen()
			printWelcome()
			continue

		case input == "help":
			printHelp()
			continue

		case input == "variables" || input == "vars":
			showVariables()
			continue

		case input == "history":
			err := history.ShowHistory()
			if err != nil {
				fmt.Printf(colorRed+"Error displaying history: %v\n"+colorReset, err)
			}
			continue

		case strings.HasPrefix(input, "precision "):
			handlePrecision(input)
			continue

		case strings.HasPrefix(input, "convert "):
			handleConversion(input)
			continue

		default:
			handleExpression(input)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf(colorRed+"Input error: %v\n"+colorReset, err)
	}
}

// printWelcome displays the welcome banner
func printWelcome() {
	fmt.Println(colorCyan + banner + colorReset)
	fmt.Println(colorBold + "  A Powerful CLI Calculator" + colorReset)
	fmt.Println(colorDim + "  Type 'help' for commands or 'exit' to quit\n" + colorReset)
}

// printHelp displays comprehensive command reference
func printHelp() {
	fmt.Println(colorCyan + "╔════════════════════════════════════════════════════════════╗" + colorReset)
	fmt.Println(colorCyan + "║" + colorBold + "                    AXION CALCULATOR                       " + colorReset + colorCyan + "║" + colorReset)
	fmt.Println(colorCyan + "╚════════════════════════════════════════════════════════════╝" + colorReset)
	fmt.Println()

	fmt.Println(colorYellow + "┌─ BASIC COMMANDS ─────────────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorGreen+"<expression>"+colorReset, "Evaluate mathematical expression")
	fmt.Printf("│ %-25s %s\n", colorGreen+"help"+colorReset, "Show this help message")
	fmt.Printf("│ %-25s %s\n", colorGreen+"exit"+colorReset, "Exit the calculator")
	fmt.Printf("│ %-25s %s\n", colorGreen+"clear"+colorReset, "Clear terminal screen")
	fmt.Printf("│ %-25s %s\n", colorGreen+"variables"+colorReset, "Show all stored variables")
	fmt.Printf("│ %-25s %s\n", colorGreen+"history"+colorReset, "Display calculation history")
	fmt.Println(colorYellow + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()

	fmt.Println(colorPurple + "┌─ MATHEMATICAL FUNCTIONS ─────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorBold+"Trigonometric:"+colorReset, "sin, cos, tan, asin, acos, atan")
	fmt.Printf("│ %-25s %s\n", colorBold+"Logarithmic:"+colorReset, "ln, log, log10, log2")
	fmt.Printf("│ %-25s %s\n", colorBold+"Exponential:"+colorReset, "exp, pow, sqrt")
	fmt.Printf("│ %-25s %s\n", colorBold+"Utility:"+colorReset, "abs, ceil, floor, round, sign")
	fmt.Printf("│ %-25s %s\n", colorBold+"Statistical:"+colorReset, "mean, median, mode, sum, product")
	fmt.Printf("│ %-25s %s\n", colorBold+"Other:"+colorReset, "max, min, mod, ! (factorial)")
	fmt.Println(colorPurple + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()

	fmt.Println(colorBlue + "┌─ VARIABLES & CONSTANTS ──────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorBold+"Assignment:"+colorReset, "x = 5, area = pi * r^2")
	fmt.Printf("│ %-25s %s\n", colorBold+"Constants:"+colorReset, "pi, e, phi, c, G, h")
	fmt.Println(colorBlue + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()

	fmt.Println(colorGreen + "┌─ UNIT CONVERSION ────────────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorBold+"Syntax:"+colorReset, "convert <value> <from> to <to>")
	fmt.Printf("│ %-25s %s\n", colorBold+"Length:"+colorReset, "m, cm, mm, km, in, ft, yd, mi")
	fmt.Printf("│ %-25s %s\n", colorBold+"Weight:"+colorReset, "kg, g, mg, lb, oz, ton")
	fmt.Printf("│ %-25s %s\n", colorBold+"Time:"+colorReset, "s, ms, min, h, d")
	fmt.Printf("│ %-25s %s\n", colorBold+"Example:"+colorReset, colorCyan+"convert 100 cm to m"+colorReset)
	fmt.Println(colorGreen + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()

	fmt.Println(colorYellow + "┌─ SETTINGS ───────────────────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorGreen+"precision <n>"+colorReset, "Set decimal precision (0-20)")
	fmt.Println(colorYellow + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()

	fmt.Println(colorCyan + "┌─ EXAMPLES ───────────────────────────────────────────────┐" + colorReset)
	fmt.Printf("│ %-25s %s\n", colorBold+"Basic:"+colorReset, "2 + 3 * 4, (10 - 5) / 2")
	fmt.Printf("│ %-25s %s\n", colorBold+"Functions:"+colorReset, "sin(30), sqrt(16), log(100)")
	fmt.Printf("│ %-25s %s\n", colorBold+"Variables:"+colorReset, "x = 10, y = x * 2")
	fmt.Printf("│ %-25s %s\n", colorBold+"Scientific:"+colorReset, "2e-10, 3.14E+5")
	fmt.Printf("│ %-25s %s\n", colorBold+"Statistics:"+colorReset, "mean(1,2,3,4,5)")
	fmt.Println(colorCyan + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()
}

// clearScreen clears the terminal display
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// formatResult formats numerical results with proper precision
func formatResult(result float64) string {
	if math.IsNaN(result) {
		return colorRed + "undefined (NaN)" + colorReset
	} else if math.IsInf(result, 1) {
		return colorYellow + "+∞" + colorReset
	} else if math.IsInf(result, -1) {
		return colorYellow + "-∞" + colorReset
	} else {
		format := fmt.Sprintf("%%.%dg", settings.Precision)
		return colorGreen + fmt.Sprintf(format, result) + colorReset
	}
}

// showVariables displays all currently stored variables
func showVariables() {
	if len(evaluator.Vars) == 0 {
		fmt.Println(colorYellow + "No variables defined." + colorReset)
		return
	}

	fmt.Println(colorCyan + "┌─ Stored Variables ───────────────────────────────────────┐" + colorReset)
	for name, value := range evaluator.Vars {
		fmt.Printf(colorCyan+"│ "+colorReset+colorBold+"%-15s"+colorReset+" = %s\n", name, formatResult(value))
	}
	fmt.Println(colorCyan + "└──────────────────────────────────────────────────────────┘" + colorReset)
	fmt.Println()
}

// handlePrecision processes precision setting commands
func handlePrecision(input string) {
	parts := strings.Fields(input)
	if len(parts) != 2 {
		fmt.Println(colorRed + "Usage: " + colorReset + "precision <number>")
		fmt.Println(colorDim + "   Example: precision 10" + colorReset)
		return
	}

	precision, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Printf(colorRed+"Invalid number: %s\n"+colorReset, parts[1])
		return
	}

	if err := settings.Set(precision); err != nil {
		fmt.Printf(colorRed+"Error: %v\n"+colorReset, err)
		return
	}

	fmt.Printf(colorGreen+"Precision set to %d decimal places\n"+colorReset, settings.Precision)
}

// handleConversion processes unit conversion commands
func handleConversion(input string) {
	parts := strings.Fields(input)
	if len(parts) != 5 || parts[3] != "to" {
		fmt.Println(colorRed + "Usage: " + colorReset + "convert <value> <from> to <to>")
		fmt.Println(colorDim + "   Example: convert 10 km to m" + colorReset)
		return
	}

	valueStr := parts[1]
	fromUnit := parts[2]
	toUnit := parts[4]

	var value float64
	_, err := fmt.Sscanf(valueStr, "%f", &value)
	if err != nil {
		fmt.Printf(colorRed+"Invalid number: %s\n"+colorReset, valueStr)
		return
	}

	result, err := units.Convert(value, fromUnit, toUnit)
	if err != nil {
		fmt.Printf(colorRed+"Conversion error: %v\n"+colorReset, err)
		return
	}

	fmt.Printf(colorBold+"%s %s"+colorReset+" = "+colorGreen+"%s %s\n"+colorReset,
		formatResult(value), fromUnit,
		formatResult(result), toUnit)
}

// handleExpression processes mathematical expressions
func handleExpression(input string) {
	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		fmt.Printf(colorRed+"Error: %v\n"+colorReset, err)
		return
	}
	p := parser.Parser{Tokens: tokens}
	ast, err := p.ParseExpression()
	if err != nil {
		fmt.Printf(colorRed+"Error: %v\n"+colorReset, err)
		return
	}

	result, err := evaluator.Eval(ast)
	if err != nil {
		fmt.Printf(colorRed+"Error: %v\n"+colorReset, err)
		return
	}

	fmt.Printf(colorBold+"Result: "+colorReset+"%s\n", formatResult(result))

	if err := history.AddHistory(input, result); err != nil {
		fmt.Printf(colorYellow+"Warning: Failed to save to history: %v\n"+colorReset, err)
	}
}
