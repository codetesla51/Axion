/*
Axion CLI Calculator - Main Entry Point
========================================
Author: Uthman
Year: 2025

This file implements the main command-line interface for the Axion calculator.
It provides a comprehensive REPL (Read-Eval-Print-Loop) that processes user input
and delegates to appropriate modules for mathematical expression evaluation, unit
conversion, variable management, and history operations.

The expression evaluation pipeline follows this sequence:
1. Input string -> Tokenizer (lexical analysis)
2. Tokens -> Parser (syntax analysis, AST construction)
3. AST -> Evaluator (expression evaluation with variable support)
4. Result -> History storage and display

Supported command categories:
- Mathematical expressions: Full expression evaluation with variables and functions
- Variable operations: Assignment, retrieval, and management
- Unit conversions: Multi-category conversion system
- History management: Persistent calculation storage and retrieval
- Settings: Precision control and display customization
- System commands: Help, clear screen, and exit functionality

Error handling is implemented at each stage to provide meaningful feedback
for invalid input, mathematical errors, domain violations, and system failures.
The REPL maintains session state including variables and settings.
*/

package main

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
)

// printHelp displays comprehensive command reference and usage examples
func printHelp() {
	fmt.Println("Axion CLI Calculator")
	fmt.Println("====================")
	fmt.Println()

	fmt.Println("BASIC COMMANDS")
	fmt.Println("--------------")
	fmt.Printf("  %-20s %s\n", "<expression>", "Evaluate mathematical expression")
	fmt.Printf("  %-20s %s\n", "help", "Show this help message")
	fmt.Printf("  %-20s %s\n", "exit", "Exit the calculator")
	fmt.Printf("  %-20s %s\n", "clear", "Clear terminal screen")
	fmt.Printf("  %-20s %s\n", "variables", "Show all stored variables")
	fmt.Printf("  %-20s %s\n", "history", "Display calculation history")
	fmt.Println()

	fmt.Println("MATHEMATICAL FUNCTIONS")
	fmt.Println("----------------------")
	fmt.Printf("  %-20s %s\n", "Trigonometric:", "sin, cos, tan, asin, acos, atan")
	fmt.Printf("  %-20s %s\n", "Logarithmic:", "ln, log, log10, log2")
	fmt.Printf("  %-20s %s\n", "Exponential:", "exp, pow, sqrt")
	fmt.Printf("  %-20s %s\n", "Utility:", "abs, ceil, floor, round, sign")
	fmt.Printf("  %-20s %s\n", "Statistical:", "mean, median, mode, sum, product")
	fmt.Printf("  %-20s %s\n", "Other:", "max, min, mod, ! (factorial)")
	fmt.Println()

	fmt.Println("VARIABLES & CONSTANTS")
	fmt.Println("---------------------")
	fmt.Printf("  %-20s %s\n", "Assignment:", "x = 5, area = pi * r^2")
	fmt.Printf("  %-20s %s\n", "Constants:", "pi, e, phi, c, G, h")
	fmt.Println()

	fmt.Println("UNIT CONVERSION")
	fmt.Println("---------------")
	fmt.Printf("  %-20s %s\n", "Syntax:", "convert <value> <from> to <to>")
	fmt.Printf("  %-20s %s\n", "Length:", "m, cm, mm, km, in, ft, yd, mi")
	fmt.Printf("  %-20s %s\n", "Weight:", "kg, g, mg, lb, oz, ton")
	fmt.Printf("  %-20s %s\n", "Time:", "s, ms, min, h, d")
	fmt.Printf("  %-20s %s\n", "Example:", "convert 100 cm to m")
	fmt.Println()

	fmt.Println("SETTINGS")
	fmt.Println("--------")
	fmt.Printf("  %-20s %s\n", "precision <n>", "Set decimal precision (0-20)")
	fmt.Println()

	fmt.Println("EXAMPLES")
	fmt.Println("--------")
	fmt.Printf("  %-20s %s\n", "Basic:", "2 + 3 * 4, (10 - 5) / 2")
	fmt.Printf("  %-20s %s\n", "Functions:", "sin(30), sqrt(16), log(100)")
	fmt.Printf("  %-20s %s\n", "Variables:", "x = 10, y = x * 2")
	fmt.Printf("  %-20s %s\n", "Scientific:", "2e-10, 3.14E+5")
	fmt.Printf("  %-20s %s\n", "Statistics:", "mean(1,2,3,4,5)")
	fmt.Println()
}

// clearScreen clears the terminal display
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// formatResult formats numerical results with proper precision
func formatResult(result float64) string {
	if math.IsNaN(result) {
		return "undefined (NaN)"
	} else if math.IsInf(result, 1) {
		return "+∞"
	} else if math.IsInf(result, -1) {
		return "-∞"
	} else {
		format := fmt.Sprintf("%%.%dg", settings.Precision)
		return fmt.Sprintf(format, result)
	}
}

// showVariables displays all currently stored variables
func showVariables() {
	if len(evaluator.Vars) == 0 {
		fmt.Println("No variables defined.")
		return
	}

	fmt.Println("Stored Variables:")
	fmt.Println("-----------------")
	for name, value := range evaluator.Vars {
		fmt.Printf("  %-10s = %s\n", name, formatResult(value))
	}
	fmt.Println()
}

func main() {
	// Initialize constants system
	err := constants.Load("constants.json")
	if err != nil {
		fmt.Printf("Warning: Failed to load constants: %v\n", err)
	}

	// Initialize input scanner
	scanner := bufio.NewScanner(os.Stdin)

	// Display welcome message
	fmt.Println("Welcome to Axion Calculator! Type 'help' for commands.")

	// Main REPL loop
	for {
		fmt.Print(">> ")

		if !scanner.Scan() {
			fmt.Println("\nGoodbye!")
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		switch {
		case input == "exit" || input == "quit":
			fmt.Println("Goodbye!")
			return

		case input == "clear" || input == "cls":
			clearScreen()
			fmt.Println("Welcome to Axion Calculator! Type 'help' for commands.")
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
				fmt.Printf("Error displaying history: %v\n", err)
			}
			continue

		case strings.HasPrefix(input, "precision "):
			parts := strings.Fields(input)
			if len(parts) != 2 {
				fmt.Println("Usage: precision <number>")
				fmt.Println("Example: precision 10")
				continue
			}

			precision, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Invalid number: %s\n", parts[1])
				continue
			}

			if err := settings.Set(precision); err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Printf("Precision set to %d decimal places\n", settings.Precision)
			continue

		case strings.HasPrefix(input, "convert "):
			parts := strings.Fields(input)
			if len(parts) != 5 || parts[3] != "to" {
				fmt.Println("Usage: convert <value> <from> to <to>")
				fmt.Println("Example: convert 10 km to m")
				continue
			}

			valueStr := parts[1]
			fromUnit := parts[2]
			toUnit := parts[4]

			var value float64
			_, err := fmt.Sscanf(valueStr, "%f", &value)
			if err != nil {
				fmt.Printf("Invalid number: %s\n", valueStr)
				continue
			}

			result, err := units.Convert(value, fromUnit, toUnit)
			if err != nil {
				fmt.Printf("Conversion error: %v\n", err)
				continue
			}

			fmt.Printf("%s %s = %s %s\n",
				formatResult(value), fromUnit,
				formatResult(result), toUnit)
			continue

		default:
			// Mathematical expression evaluation
			tokens, err := tokenizer.Tokenize(input)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			p := parser.Parser{Tokens: tokens}
			ast := p.ParseExpression()
			if ast == nil {
				fmt.Println("Error: Invalid expression")
				continue
			}

			result, err := evaluator.Eval(ast)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Printf("Result: %s\n", formatResult(result))

			if err := history.AddHistory(input, result); err != nil {
				fmt.Printf("Warning: Failed to save to history: %v\n", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Input error: %v\n", err)
	}
}
