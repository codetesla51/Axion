package tokenizer

import (
	"fmt"
	"unicode"
)

type TokenType int

const (
	NUMBER TokenType = iota
	OPERATOR
	PAREN
	FUNCTION
)

type Token struct {
	Type  TokenType
	Value string
}
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	var numberBuffer string
	var wordBuffer string

	addToken := func(t Token) {
		// Check for implicit multiplication
		if len(tokens) > 0 {
			last := tokens[len(tokens)-1]
			if (last.Type == NUMBER || (last.Type == PAREN && last.Value == ")")) &&
				(t.Type == NUMBER || t.Type == FUNCTION || (t.Type == PAREN && t.Value == "(")) {
				// Insert multiplication
				tokens = append(tokens, Token{Type: OPERATOR, Value: "*"})
			}
		}
		tokens = append(tokens, t)
	}

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		switch {
		case unicode.IsDigit(ch) || ch == '.':
			if ch == '.' && containsDot(numberBuffer) {
				return nil, fmt.Errorf("invalid number: multiple decimals in %q", numberBuffer+string(ch))
			}
			numberBuffer += string(ch)

			if i+1 < len(input) && (input[i+1] == 'e' || input[i+1] == 'E') {
				i++
				numberBuffer += string(input[i])
				if i+1 < len(input) && (input[i+1] == '+' || input[i+1] == '-') {
					i++
					numberBuffer += string(input[i])
				}
				digitsFound := false
				for i+1 < len(input) && unicode.IsDigit(rune(input[i+1])) {
					i++
					numberBuffer += string(input[i])
					digitsFound = true
				}
				if !digitsFound {
					return nil, fmt.Errorf("invalid exponent in number: %q", numberBuffer)
				}
			}

		case unicode.IsLetter(ch):
			wordBuffer += string(ch)

		case ch == '+' || ch == '-' || ch == '*' || ch == '/' || ch == '^':
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: OPERATOR, Value: string(ch)})

		case ch == '(' || ch == ')':
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: PAREN, Value: string(ch)})

		case ch == '!':
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: FUNCTION, Value: "!"})

		case unicode.IsSpace(ch):
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}

		case ch == ',':
			if numberBuffer != "" {
				addToken(Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				addToken(Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			addToken(Token{Type: OPERATOR, Value: ","})

		default:
			return nil, fmt.Errorf("invalid character: %q", ch)
		}
	}

	if numberBuffer != "" {
		addToken(Token{Type: NUMBER, Value: numberBuffer})
	}
	if wordBuffer != "" {
		addToken(Token{Type: FUNCTION, Value: wordBuffer})
	}

	return tokens, nil
}

func containsDot(s string) bool {
	for _, ch := range s {
		if ch == '.' {
			return true
		}
	}
	return false
}
