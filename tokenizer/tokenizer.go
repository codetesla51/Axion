package tokenizer

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

func Tokenize(input string) []Token {
	var tokens []Token
	var wordBuffer string
	var numberBuffer string

	for _, ch := range input {
		switch {
		case ch >= '0' && ch <= '9' || ch == '.':
			numberBuffer += string(ch)
		case ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z':
			wordBuffer += string(ch)
		case ch == '+' || ch == '-' || ch == '*' || ch == '^' || ch == '/':
			if numberBuffer != "" {
				tokens = append(tokens, Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			tokens = append(tokens, Token{Type: OPERATOR, Value: string(ch)})
		case ch == '(' || ch == ')':
			if numberBuffer != "" {
				tokens = append(tokens, Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				tokens = append(tokens, Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
			tokens = append(tokens, Token{Type: PAREN, Value: string(ch)})
		case ch == ' ':
			if numberBuffer != "" {
				tokens = append(tokens, Token{Type: NUMBER, Value: numberBuffer})
				numberBuffer = ""
			}
			if wordBuffer != "" {
				tokens = append(tokens, Token{Type: FUNCTION, Value: wordBuffer})
				wordBuffer = ""
			}
		}
	}
	if numberBuffer != "" {
		tokens = append(tokens, Token{Type: NUMBER, Value: numberBuffer})
	}
	if wordBuffer != "" {
		tokens = append(tokens, Token{Type: FUNCTION, Value: wordBuffer})
	}
	return tokens
}
