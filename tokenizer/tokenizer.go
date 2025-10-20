/*
Tokenizer Module - Lexical Analysis Engine
===========================================
Part of Axion CLI Calculator
Author: Uthman
Year: 2025

This module implements comprehensive lexical analysis for mathematical expressions,
transforming input strings into structured token sequences for parser consumption.
The tokenizer serves as the first stage in the expression evaluation pipeline.

Core Capabilities:
- Numeric literal parsing with full scientific notation support (1.5e-10, 2E+5)
- Mathematical operator recognition (+, -, *, /, ^, !, =)
- Function name identification and classification
- Parentheses and comma handling for grouping and function arguments
- Intelligent implicit multiplication insertion (2sin(x) → 2 * sin(x))
- Variable and identifier recognition
- Assignment operator support

Key Features:
- Scientific Notation: Complete support for exponential format including optional signs
- Input Validation: Prevents malformed numbers (e.g., 3.14.15) and invalid characters
- Implicit Multiplication: Automatically inserts multiplication between adjacent operands
- Buffer Management: Ensures complete token extraction with proper boundary handling
- Error Reporting: Provides detailed error messages with context for invalid input

Token Categories:
- NUMBER: Numeric literals including decimals and scientific notation
- OPERATOR: Mathematical operators and separators
- FUNCTION: Built-in mathematical functions (sin, cos, log, etc.)
- IDENT: User-defined variables and identifiers
- PAREN: Grouping operators for precedence control
- ASSIGN: Variable assignment operator

The tokenizer maintains mathematical expression integrity while providing
flexibility for complex calculations and function compositions.
*/

package tokenizer

import (
	"fmt"
	"unicode"
)

// TokenType defines the categories of tokens
type TokenType int

const (
	NUMBER     TokenType = iota // 0 - numeric literals
	OPERATOR                    // 1 - arithmetic operators (+, -, *, /, ^)
	PAREN                       // 2 - parentheses
	FUNCTION                    // 3 - functions (sin, cos, log, etc.)
	IDENT                       // 4 - identifiers/variables
	ASSIGN                      // 5 - assignment (=)
	COMPARISON                  // 6 - comparison operators (>, <, >=, <=, ==, !=)
	LOGICAL                     // 7 - logical operators (&&, ||, !)
)

// Token represents a lexical unit
type Token struct {
	Type  TokenType
	Value string
}

// flushBuffers adds buffered content as tokens
func flushBuffers(numberBuffer, wordBuffer *string, addToken func(Token)) {
	if *numberBuffer != "" {
		addToken(Token{Type: NUMBER, Value: *numberBuffer})
		*numberBuffer = ""
	}
	if *wordBuffer != "" {
		if isMathFunction(*wordBuffer) {
			addToken(Token{Type: FUNCTION, Value: *wordBuffer})
		} else {
			addToken(Token{Type: IDENT, Value: *wordBuffer})
		}
		*wordBuffer = ""
	}
}

// Tokenize performs lexical analysis on the input expression
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	var numberBuffer string
	var wordBuffer string

	// Helper function to add tokens with implicit multiplication
	addToken := func(t Token) {
		// Insert implicit multiplication operators where needed
		if len(tokens) > 0 {
			last := tokens[len(tokens)-1]

			// Special case: factorial operators do not require preceding multiplication
			if t.Type == FUNCTION && t.Value == "!" {
				// No implicit multiplication before factorial
			} else if (last.Type == NUMBER || (last.Type == PAREN && last.Value == ")")) &&
				(t.Type == NUMBER || t.Type == FUNCTION || t.Type == IDENT || (t.Type == PAREN && t.Value == "(")) {
				// Insert multiplication between appropriate tokens
				tokens = append(tokens, Token{Type: OPERATOR, Value: "*"})
			}
		}
		tokens = append(tokens, t)
	}

	// Character-by-character processing
	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		switch {
		// Handle digits and decimal points
		case unicode.IsDigit(ch) || ch == '.':
			// Validate decimal point usage
			if ch == '.' && containsDot(numberBuffer) {
				return nil, fmt.Errorf("invalid number: multiple decimal points in %q", numberBuffer+string(ch))
			}
			if wordBuffer != "" {
				wordBuffer += string(ch)
				continue
			}
			numberBuffer += string(ch)

			// Handle scientific notation
			if i+1 < len(input) && (input[i+1] == 'e' || input[i+1] == 'E') {
				i++
				numberBuffer += string(input[i])

				// Handle optional sign
				if i+1 < len(input) && (input[i+1] == '+' || input[i+1] == '-') {
					i++
					numberBuffer += string(input[i])
				}

				// Require digits after exponent
				digitsFound := false
				for i+1 < len(input) && unicode.IsDigit(rune(input[i+1])) {
					i++
					numberBuffer += string(input[i])
					digitsFound = true
				}

				if !digitsFound {
					return nil, fmt.Errorf("invalid scientific notation in %q", numberBuffer)
				}
			}

		// Handle alphabetic characters (functions and identifiers)
		case unicode.IsLetter(ch):
			wordBuffer += string(ch)

		// Handle mathematical operators
		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			addToken(Token{Type: OPERATOR, Value: string(ch)})

		// Handle assignment operator
		case ch == '=':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			if i+1 < len(input) {
				next := rune(input[i+1])
				if next == '=' {
					addToken(Token{Type: COMPARISON, Value: string(ch) + string(next)})
					i++
					continue
				}
			}
			addToken(Token{Type: ASSIGN, Value: "="})

		case ch == '>' || ch == '<':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			if i+1 < len(input) {
				next := rune(input[i+1])
				if next == '=' {
					addToken(Token{Type: COMPARISON, Value: string(ch) + string(next)})
					i++
					continue
				}
			}
			addToken(Token{Type: COMPARISON, Value: string(ch)})

		case ch == '&' || ch == '|':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			if i+1 < len(input) {
				next := rune(input[i+1])
				if next == ch {
					addToken(Token{Type: LOGICAL, Value: string(ch) + string(next)})
					i++
					continue
				}
			}
			return nil, fmt.Errorf("invalid logical operator: %q", ch)

		case ch == '(' || ch == ')':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			addToken(Token{Type: PAREN, Value: string(ch)})

		// Handle factorial
		case ch == '!':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			if i+1 < len(input) {
				next := rune(input[i+1])
				if next == '=' {
					addToken(Token{Type: COMPARISON, Value: string(ch) + string(next)})
					i++
					continue
				}
			}

			addToken(Token{Type: FUNCTION, Value: "!"})

		// Handle whitespace
		case unicode.IsSpace(ch):
			flushBuffers(&numberBuffer, &wordBuffer, addToken)

		// Handle comma (function argument separator)
		case ch == ',':
			flushBuffers(&numberBuffer, &wordBuffer, addToken)
			addToken(Token{Type: OPERATOR, Value: ","})

		// Handle invalid characters
		default:
			return nil, fmt.Errorf("invalid character: %q", ch)
		}
	}

	// Process any remaining buffered content
	flushBuffers(&numberBuffer, &wordBuffer, addToken)

	return tokens, nil
}

// containsDot checks for decimal point in number buffer
func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}
	return false
}

// isMathFunction checks if a word is a mathematical function
func isMathFunction(word string) bool {
	functions := map[string]bool{
		"sin": true, "cos": true, "tan": true,
		"asin": true, "acos": true, "atan": true, "atan2": true,

		// Logarithmic
		"log": true, "log10": true, "log2": true, "ln": true,

		// Power/Root
		"sqrt": true, "exp": true, "pow": true,

		// Rounding
		"abs": true, "ceil": true, "floor": true, "round": true, "trunc": true,

		// Utility
		"sign": true, "mod": true,

		// Conversion
		"deg2rad": true, "rad2deg": true,

		// Statistical
		"max": true, "min": true, "mean": true, "median": true,
		"mode": true, "sum": true, "product": true,

		//reserved
		"print": true,
	}
	return functions[word]
}