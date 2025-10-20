package tokenizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize_Numbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "integer",
			input: "42",
			want:  []Token{{Type: NUMBER, Value: "42"}},
		},
		{
			name:  "decimal",
			input: "3.14",
			want:  []Token{{Type: NUMBER, Value: "3.14"}},
		},
		{
			name:  "scientific notation",
			input: "1.5e-10",
			want:  []Token{{Type: NUMBER, Value: "1.5e-10"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTokenize_Functions(t *testing.T) {
	tests := []struct {
		name, input string
		want        []Token
	}{
		{"sin", "sin", []Token{{Type: FUNCTION, Value: "sin"}}},
		{
			"cos", "cos", []Token{{Type: FUNCTION, Value: "cos"}},
		},
		{
			"tan", "tan", []Token{{Type: FUNCTION, Value: "tan"}},
		},
		{
			"log", "log", []Token{{Type: FUNCTION, Value: "log"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestTokenize_Ident(t *testing.T) {
	tests := []struct {
		name, input string
		want        []Token
	}{
		{"Idnet1", "var", []Token{{Type: IDENT, Value: "var"}}},
		{
			"Ident2", "r", []Token{{Type: IDENT, Value: "r"}},
		},
		{
			"Ident3", "time", []Token{{Type: IDENT, Value: "time"}},
		},
		{
			"Ident4", "root", []Token{{Type: IDENT, Value: "root"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTokenize_Operators(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "addition",
			input: "3+4",
			want: []Token{
				{Type: NUMBER, Value: "3"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "4"},
			},
		},
		{
			name:  "subtraction",
			input: "5-2",
			want: []Token{
				{Type: NUMBER, Value: "5"},
				{Type: OPERATOR, Value: "-"},
				{Type: NUMBER, Value: "2"},
			},
		},
		{
			name:  "multiplication",
			input: "2*3",
			want: []Token{
				{Type: NUMBER, Value: "2"},
				{Type: OPERATOR, Value: "*"},
				{Type: NUMBER, Value: "3"},
			},
		},
		{
			name:  "division",
			input: "8/4",
			want: []Token{
				{Type: NUMBER, Value: "8"},
				{Type: OPERATOR, Value: "/"},
				{Type: NUMBER, Value: "4"},
			},
		},
		{
			name:  "exponentiation",
			input: "2^3",
			want: []Token{
				{Type: NUMBER, Value: "2"},
				{Type: OPERATOR, Value: "^"},
				{Type: NUMBER, Value: "3"},
			},
		},
		{
			name:  "comma",
			input: "max(1,2)",
			want: []Token{
				{Type: FUNCTION, Value: "max"},
				{Type: PAREN, Value: "("},
				{Type: NUMBER, Value: "1"},
				{Type: OPERATOR, Value: ","},
				{Type: NUMBER, Value: "2"},
				{Type: PAREN, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTokenize_Parentheses(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "simple grouping",
			input: "(3+4)",
			want: []Token{
				{Type: PAREN, Value: "("},
				{Type: NUMBER, Value: "3"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "4"},
				{Type: PAREN, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTokenize_ImplicitMultiplication(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "number before function",
			input: "2sin(x)",
			want: []Token{
				{Type: NUMBER, Value: "2"},
				{Type: OPERATOR, Value: "*"}, // IMPLICIT
				{Type: FUNCTION, Value: "sin"},
				{Type: PAREN, Value: "("},
				{Type: IDENT, Value: "x"},
				{Type: PAREN, Value: ")"},
			},
		},
		{
			name:  "parentheses multiplication",
			input: "2(3+4)",
			want: []Token{
				{Type: NUMBER, Value: "2"},
				{Type: OPERATOR, Value: "*"}, // IMPLICIT
				{Type: PAREN, Value: "("},
				{Type: NUMBER, Value: "3"},
				{Type: OPERATOR, Value: "+"},
				{Type: NUMBER, Value: "4"},
				{Type: PAREN, Value: ")"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestTokenize_LogicalAndComparison(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "logical AND",
			input: "4 && 5",
			want: []Token{
				{Type: NUMBER, Value: "4"},
				{Type: LOGICAL, Value: "&&"},
				{Type: NUMBER, Value: "5"},
			},
		},
		{
			name:  "logical OR",
			input: "4 || 5",
			want: []Token{
				{Type: NUMBER, Value: "4"},
				{Type: LOGICAL, Value: "||"},
				{Type: NUMBER, Value: "5"},
			},
		},
		{
			name:  "greater than",
			input: "6 > 5",
			want: []Token{
				{Type: NUMBER, Value: "6"},
				{Type: COMPARISON, Value: ">"},
				{Type: NUMBER, Value: "5"},
			},
		},
		{
			name:  "greater than or equal to",
			input: "6 >= 5",
			want: []Token{
				{Type: NUMBER, Value: "6"},
				{Type: COMPARISON, Value: ">="},
				{Type: NUMBER, Value: "5"},
			},
		},
		{
			name:  "less than or equal to",
			input: "4 <= 9",
			want: []Token{
				{Type: NUMBER, Value: "4"},
				{Type: COMPARISON, Value: "<="},
				{Type: NUMBER, Value: "9"},
			},
		},
		{
			name:  "equality",
			input: "6 == 7",
			want: []Token{
				{Type: NUMBER, Value: "6"},
				{Type: COMPARISON, Value: "=="},
				{Type: NUMBER, Value: "7"},
			},
		},
		{
			name:  "not equal",
			input: "3 != 4",
			want: []Token{
				{Type: NUMBER, Value: "3"},
				{Type: COMPARISON, Value: "!="},
				{Type: NUMBER, Value: "4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestTokenize_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid character", "3 @ 4"},
		{"multiple decimals", "3.14.15"},
		{"bad scientific notation", "1.5e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Tokenize(tt.input)
			assert.Error(t, err) // MUST have error
		})
	}
}
func TestTokenize_containsBool(t *testing.T) {
	tests := []struct {
		name, input string
		want        bool
	}{
		{"bool1", "1.24", true},
		{"bool2", "0.5729", true},
		{"bool3", "5005", false},
		{"bool4", "4647", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsDot(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkTokenize_Simple(b *testing.B) {
	input := "3+4*5"
	for i := 0; i < b.N; i++ {
		Tokenize(input)
	}
}

func BenchmarkTokenize_Complex(b *testing.B) {
	input := "2*sin(3.14)+sqrt(16)/log(100)"
	for i := 0; i < b.N; i++ {
		Tokenize(input)
	}
}

func BenchmarkTokenize_Scientific(b *testing.B) {
	input := "1.5e-10+2.3E+5"
	for i := 0; i < b.N; i++ {
		Tokenize(input)
	}
}
