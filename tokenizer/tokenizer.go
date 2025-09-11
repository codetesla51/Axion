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

package tokenizer

import (
	"fmt"
	"unicode"
)

// TokenType defines the categories of tokens recognized by the lexer
type TokenType int

const (
	NUMBER   TokenType = iota // Numeric literals: 42, 3.14, 1.5e-10
	OPERATOR                  // Mathematical operators: +, -, *, /, ^
	PAREN                     // Parentheses: (, )
	FUNCTION                  // Function identifiers: sin, cos, tan, log, sqrt, !
	IDENT 
   ASSIGN 
)

// Token represents a lexical unit with its classification and value
type Token struct {
	Type  TokenType // Token classification
	Value string    // Literal token content
}

// Tokenize performs lexical analysis on the input expression
// Returns a sequence of tokens or an error for invalid input
func Tokenize(input string) ([]Token, error) {
	var tokens []Token      // Accumulated token sequence
	var numberBuffer string // Buffer for constructing numeric tokens
	var wordBuffer string   // Buffer for constructing function name tokens

	// addToken handles token insertion with implicit multiplication logic
	addToken := func(t Token) {
		// Insert implicit multiplication operators where mathematically appropriate
		if len(tokens) > 0 {
			last := tokens[len(tokens)-1]

			// Special case: factorial operators do not require preceding multiplication
			if t.Type == FUNCTION && t.Value == "!" {
				// No implicit multiplication before factorial
			} else if (last.Type == NUMBER || (last.Type == PAREN && last.Value == ")")) &&
				(t.Type == NUMBER || t.Type == FUNCTION || (t.Type == PAREN && t.Value == "(")) {
				// Insert multiplication between:
				// NUMBER + NUMBER: "2 3" -> "2 * 3"
				// NUMBER + FUNCTION: "2sin" -> "2 * sin"
				// NUMBER + OPEN_PAREN: "2(" -> "2 * ("
				// CLOSE_PAREN + NUMBER: ")3" -> ") * 3"
				// CLOSE_PAREN + FUNCTION: ")sin" -> ") * sin"
				// CLOSE_PAREN + OPEN_PAREN: ")(" -> ") * ("
				tokens = append(tokens, Token{Type: OPERATOR, Value: "*"})
			}
		}
		tokens = append(tokens, t)
	}

	// Character-by-character lexical analysis
	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		switch {
		// Process numeric characters and decimal points
		case unicode.IsDigit(ch) || ch == '.':
			// Validate decimal point usage - only one per numeric literal
			if ch == '.' && containsDot(numberBuffer) {
				return nil, fmt.Errorf("invalid number: multiple decimals in %q", numberBuffer+string(ch))
			}
			numberBuffer += string(ch)

			// Handle scientific notation (exponential form)
			if i+1 < len(input) && (input[i+1] == 'e' || input[i+1] == 'E') {
				i++
				numberBuffer += string(input[i])

				// Process optional sign in exponent
				if i+1 < len(input) && (input[i+1] == '+' || input[i+1] == '-') {
					i++
					numberBuffer += string(input[i])
				}

				// Require digits in exponent
				digitsFound := false
				for i+1 < len(input) && unicode.IsDigit(rune(input[i+1])) {
					i++
					numberBuffer += string(input[i])
					digitsFound = true
				}

				// Validate complete exponent format
				if !digitsFound {
					return nil, fmt.Errorf("invalid exponent in number: %q", numberBuffer)
				}
			}

		// Process alphabetic characters (function names)
		case unicode.IsLetter(ch):
    wordBuffer := string(ch) 
    i++
    for i < len(input) && (unicode.IsLetter(rune(input[i])) || unicode.IsDigit(rune(input[i]))) {
        wordBuffer += string(input[i])
        i++
    }
    tokens = append(tokens, Token{Type: IDENT, Value: wordBuffer})
		// Process mathematical operators
		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^':
			// Flush accumulated buffers before processing operator
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: OPERATOR, Value: string(ch)})
case ch == '=':
    tokens = append(tokens, Token{Type: ASSIGN, Value: "="})
		// Process parentheses
		case ch == '(' || ch == ')':
			// Flush accumulated buffers before processing parenthesis
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: PAREN, Value: string(ch)})

		// Process factorial operator
		case ch == '!':
			// Flush accumulated buffers before processing factorial
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: FUNCTION, Value: "!"})

		// Process whitespace (token separator, not preserved)
		case unicode.IsSpace(ch):
			// Flush buffers on whitespace but do not create whitespace tokens
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}

		// Process comma (function argument separator)
		case ch == ',':
			// Flush accumulated buffers before processing comma
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			// Treat comma as operator for parsing purposes
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
		addToken(Token{Type: FUNCTION, Value: wordBuffer})
	}

	return tokens, nil
}

// containsDot checks for existing decimal point in numeric buffer
// Used to validate numeric literal format during tokenization
func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}
	return false
}
