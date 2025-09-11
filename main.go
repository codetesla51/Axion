/*
Axion CLI Calculator - Main Entry Point
========================================
Author: Uthman
Year: 2025

This file implements the main command-line interface for the Axion calculator.
It provides a REPL (Read-Eval-Print-Loop) that processes user input and delegates
to appropriate modules for mathematical expression evaluation, unit conversion,
and history management.

The expression evaluation pipeline follows this sequence:
1. Input string -> Tokenizer (lexical analysis)
2. Tokens -> Parser (syntax analysis, AST construction)
3. AST -> Evaluator (expression evaluation)
4. Result -> History storage and display

Supported commands:
- Mathematical expressions: evaluated through the tokenizer->parser->evaluator pipeline
- Unit conversions: handled by the units module
- History display: managed by the history module
- Help and exit commands

Error handling is implemented at each stage to provide meaningful feedback
for invalid input, mathematical errors, and system failures.
*/

package main

import (
	"Axion/evaluator"
	"Axion/history"
	"Axion/parser"
	"Axion/tokenizer"
	"Axion/units"
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// printHelp displays available commands and usage examples to the user
func printHelp() {
	fmt.Println("Axion CLI Calculator")
	fmt.Println("--------------------")
	fmt.Println("Commands:")
	fmt.Println("  <expression>           Evaluate a math expression")
	fmt.Println("  convert <v> <from> to <to>   Convert units (e.g. convert 10 km to m)")
	fmt.Println("  history                Show calculation history")
	fmt.Println("  exit                   Exit the program")
	fmt.Println("  help                   Show this help message")
	fmt.Println()
}

func main() {
	// Initialize scanner for reading from standard input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Axion Calculator! Type 'help' for commands.")
	
	// Main REPL loop - continues until user exits or EOF
	for {
		fmt.Print(">> ")
		
		// Read next line of input
		if !scanner.Scan() {
			// EOF reached (Ctrl+D on Unix, Ctrl+Z on Windows)
			break
		}
		
		// Clean input by removing leading and trailing whitespace
		input := strings.TrimSpace(scanner.Text())

		// Skip processing empty lines
		if input == "" {
			continue
		}

		// Process input based on command type
		switch {
		
		case input == "exit":
			fmt.Println("Goodbye!")
			return

		case input == "help":
			printHelp()
			continue

		case input == "history":
			history.ShowHistory()
			continue

		// Handle unit conversion commands with format: "convert <value> <from> to <to>"
		case strings.HasPrefix(input, "convert "):
			// Parse conversion command into components
			parts := strings.Fields(input)
			if len(parts) != 5 || parts[3] != "to" {
				fmt.Println("Usage: convert <value> <from> to <to>")
				fmt.Println("Example: convert 10 km to m")
				continue
			}

			// Extract conversion parameters
			valueStr := parts[1]
			from := parts[2]
			to := parts[4]

			// Parse numeric value from string
			var value float64
			_, err := fmt.Sscanf(valueStr, "%f", &value)
			if err != nil {
				fmt.Println("Invalid number:", valueStr)
				continue
			}

			// Perform unit conversion
			result, err := units.Convert(value, from, to)
			if err != nil {
				fmt.Println("Conversion error:", err)
				continue
			}

			// Display conversion result
			fmt.Printf("%g %s = %g %s\n", value, from, result, to)
			continue

		// Default case: treat input as mathematical expression
		default:
			tokens, err := tokenizer.Tokenize(input)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			p := parser.Parser{Tokens: tokens}
			ast := p.ParseExpression()
			fmt.Printf("%+v\n", ast)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			
			if math.IsNaN(result) {
				fmt.Println("Result: undefined (NaN)")
			} else if math.IsInf(result, 1) {
				fmt.Println("Result: +∞")
			} else if math.IsInf(result, -1) {
				fmt.Println("Result: -∞")
			} else {
				fmt.Printf("Result: %g\n", result)
			}

			err = history.AddHistory(input, result)
			if err != nil {
				fmt.Println("Failed to save history:", err)
			}
		}
	}
}