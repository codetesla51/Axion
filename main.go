/*
Axion CLI Calculator
--------------------
Author: Uthman
Year: 2025
GitHub: https://github.com/codetesla51/Axion
License: MIT

This is a simple CLI calculator with support for:
- Evaluating math expressions
- Unit conversion (e.g., convert 10 km to m)
- Calculation history
- Built-in math functions: sin, cos, tan, log, sqrt, factorial, etc.

MIT License

Copyright (c) 2025 Uthman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
	"os"
	"strings"
)

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
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Axion Calculator! Type 'help' for commands.")
	for {
		fmt.Print(">> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

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

		case strings.HasPrefix(input, "convert "):
			parts := strings.Fields(input)
			if len(parts) != 5 || parts[3] != "to" {
				fmt.Println("Usage: convert <value> <from> to <to>")
				fmt.Println("Example: convert 10 km to m")
				continue
			}

			valueStr := parts[1]
			from := parts[2]
			to := parts[4]

			var value float64
			_, err := fmt.Sscanf(valueStr, "%f", &value)
			if err != nil {
				fmt.Println("Invalid number:", valueStr)
				continue
			}

			result, err := units.Convert(value, from, to)
			if err != nil {
				fmt.Println("Conversion error:", err)
				continue
			}

			fmt.Printf("%g %s = %g %s\n", value, from, result, to)
			continue

		default:
			tokens, err := tokenizer.Tokenize(input)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			p := parser.Parser{Tokens: tokens}
			ast := p.ParseExpression()

			result, err := evaluator.Eval(ast)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Printf("Result: %g\n", result)

			err = history.AddHistory(input, result)
			if err != nil {
				fmt.Println("Failed to save history:", err)
			}
		}
	}
}
