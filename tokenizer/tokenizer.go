/*
Tokenizer Module - Lexical Analysis
====================================

This module implements lexical analysis for mathematical expressions, converting
input strings into sequences of tokens that can be processed by the parser.

The tokenizer handles:
- Numeric literals including scientific notation (1.5e-10)
- Mathematical operators (+, -, *, /, ^)
- Function names (sin, cos, log, etc.)
- Parentheses for grouping
- Implicit multiplication insertion (2sin(x) -> 2 * sin(x))

Key implementation details:
- Scientific notation is parsed as single tokens to maintain numeric integrity
- Implicit multiplication is inserted between adjacent operands
- Decimal validation prevents malformed numbers (3.14.15)
- Buffer management ensures complete token extraction
*/

// File: tokenizer/tokenizer.go
// Fixed version of your existing tokenizer

package tokenizer

import (
	"fmt"
	"unicode"
)

// TokenType defines the categories of tokens
type TokenType int

const (
	NUMBER   TokenType = iota // 0
	OPERATOR                  // 1
	PAREN                     // 2
	FUNCTION                  // 3
	IDENT                     // 4
	ASSIGN                    // 5
)

// Token represents a lexical unit
type Token struct {
	Type  TokenType
	Value string
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
			// Flush buffers
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}
			addToken(Token{Type: OPERATOR, Value: string(ch)})

		// Handle assignment operator
		case ch == '=':
			// Flush buffers
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}
			addToken(Token{Type: ASSIGN, Value: "="})

		// Handle parentheses
		case ch == '(' || ch == ')':
			// Flush buffers
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}
			addToken(Token{Type: PAREN, Value: string(ch)})

		// Handle factorial
		case ch == '!':
			// Flush buffers
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}
			addToken(Token{Type: FUNCTION, Value: "!"})

		// Handle whitespace
		case unicode.IsSpace(ch):
			// Flush buffers on whitespace
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}

		// Handle comma (function argument separator)
		case ch == ',':
			// Flush buffers
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				if isMathFunction(wordBuffer) {
					addToken(Token{Type: FUNCTION, Value: wordBuffer})
				} else {
					addToken(Token{Type: IDENT, Value: wordBuffer})
				}
				wordBuffer = ""
			}
			addToken(Token{Type: OPERATOR, Value: ","})

		// Handle invalid characters
		default:
			return nil, fmt.Errorf("invalid character: %q", ch)
		}
	}

	// Process any remaining buffered content
	if numberBuffer != "" {
		addToken(Token{Type: NUMBER, Value: numberBuffer})
	}
	if wordBuffer != "" {
		if isMathFunction(wordBuffer) {
			addToken(Token{Type: FUNCTION, Value: wordBuffer})
		} else {
			addToken(Token{Type: IDENT, Value: wordBuffer})
		}
	}

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
		"asin": true, "acos": true, "atan": true,
		"sqrt": true, "exp": true, "abs": true,
		"ceil": true, "floor": true, "log": true,
		"log10": true, "pow": true, "max": true,
		"min": true, "mean": true, "median": true,
		"mode": true,
	}
	return functions[word]
}
