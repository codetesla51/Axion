package parser

import (
	"Axion/tokenizer"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParser_Expression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Node
	}{
		// Basic arithmetic operations
		{
			name:  "addition",
			input: "2+3",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{Type: NODE_NUMBER, Value: "3"},
			},
		},
		{
			name:  "subtraction",
			input: "10-5",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "-",
				Left:  &Node{Type: NODE_NUMBER, Value: "10"},
				Right: &Node{Type: NODE_NUMBER, Value: "5"},
			},
		},
		{
			name:  "multiplication",
			input: "3*4",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "*",
				Left:  &Node{Type: NODE_NUMBER, Value: "3"},
				Right: &Node{Type: NODE_NUMBER, Value: "4"},
			},
		},
		{
			name:  "division",
			input: "20/4",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "/",
				Left:  &Node{Type: NODE_NUMBER, Value: "20"},
				Right: &Node{Type: NODE_NUMBER, Value: "4"},
			},
		},
		{
			name:  "exponentiation",
			input: "2^3",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "^",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{Type: NODE_NUMBER, Value: "3"},
			},
		},

		// Operator precedence
		{
			name:  "multiplication precedence over addition",
			input: "2+3*4",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "*",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
					Right: &Node{Type: NODE_NUMBER, Value: "4"},
				},
			},
		},
		{
			name:  "division precedence over subtraction",
			input: "10-6/2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "-",
				Left:  &Node{Type: NODE_NUMBER, Value: "10"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "/",
					Left:  &Node{Type: NODE_NUMBER, Value: "6"},
					Right: &Node{Type: NODE_NUMBER, Value: "2"},
				},
			},
		},
		{
			name:  "exponentiation precedence over multiplication",
			input: "2*3^2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "*",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "^",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
					Right: &Node{Type: NODE_NUMBER, Value: "2"},
				},
			},
		},

		// Right associativity
		{
			name:  "exponentiation right associative",
			input: "2^3^2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "^",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "^",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
					Right: &Node{Type: NODE_NUMBER, Value: "2"},
				},
			},
		},
		{
			name:  "triple exponentiation",
			input: "2^2^2^2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "^",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "^",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
					Right: &Node{
						Type:  NODE_OPERATOR,
						Value: "^",
						Left:  &Node{Type: NODE_NUMBER, Value: "2"},
						Right: &Node{Type: NODE_NUMBER, Value: "2"},
					},
				},
			},
		},

		// Parentheses
		{
			name:  "parentheses override precedence",
			input: "(2+3)*4",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "*",
				Left: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
					Right: &Node{Type: NODE_NUMBER, Value: "3"},
				},
				Right: &Node{Type: NODE_NUMBER, Value: "4"},
			},
		},
		{
			name:  "nested parentheses",
			input: "2+(3*4)",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "*",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
					Right: &Node{Type: NODE_NUMBER, Value: "4"},
				},
			},
		},
		{
			name:  "deeply nested parentheses",
			input: "((((5))))",
			expected: &Node{Type: NODE_NUMBER, Value: "5"},
		},
		{
			name:  "multiple parenthesized groups",
			input: "(2+3)*(4+5)",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "*",
				Left: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
					Right: &Node{Type: NODE_NUMBER, Value: "3"},
				},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "4"},
					Right: &Node{Type: NODE_NUMBER, Value: "5"},
				},
			},
		},

		// Unary operators
		{
			name:  "unary minus",
			input: "-5",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "neg",
				Left:  &Node{Type: NODE_NUMBER, Value: "5"},
			},
		},
		{
			name:  "unary plus (ignored)",
			input: "+5",
			expected: &Node{Type: NODE_NUMBER, Value: "5"},
		},
		{
			name:  "unary minus in expression",
			input: "2+-3",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "neg",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
				},
			},
		},
		{
			name:  "unary minus with multiplication",
			input: "-2*3",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "*",
				Left: &Node{
					Type:  NODE_OPERATOR,
					Value: "neg",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				},
				Right: &Node{Type: NODE_NUMBER, Value: "3"},
			},
		},
		{
			name:  "unary minus with parentheses",
			input: "-(2+3)",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "neg",
				Left: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
					Right: &Node{Type: NODE_NUMBER, Value: "3"},
				},
			},
		},

		// Factorial
		{
			name:  "factorial",
			input: "5!",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "!",
				Children: []*Node{
					{Type: NODE_NUMBER, Value: "5"},
				},
			},
		},
		{
			name:  "factorial with parentheses",
			input: "(2+3)!",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "!",
				Children: []*Node{
					{
						Type:  NODE_OPERATOR,
						Value: "+",
						Left:  &Node{Type: NODE_NUMBER, Value: "2"},
						Right: &Node{Type: NODE_NUMBER, Value: "3"},
					},
				},
			},
		},
		{
			name:  "double factorial",
			input: "5!!",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "!",
				Children: []*Node{
					{
						Type:  NODE_FUNCTION,
						Value: "!",
						Children: []*Node{
							{Type: NODE_NUMBER, Value: "5"},
						},
					},
				},
			},
		},
		{
			name:  "factorial in expression",
			input: "3!+2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left: &Node{
					Type:  NODE_FUNCTION,
					Value: "!",
					Children: []*Node{
						{Type: NODE_NUMBER, Value: "3"},
					},
				},
				Right: &Node{Type: NODE_NUMBER, Value: "2"},
			},
		},

		// Assignment
		{
			name:  "simple assignment",
			input: "r=1",
			expected: &Node{
				Type:  NODE_ASSIGN,
				Value: "r",
				Right: &Node{Type: NODE_NUMBER, Value: "1"},
			},
		},
		{
			name:  "assignment with expression",
			input: "r=1+2*3",
			expected: &Node{
				Type:  NODE_ASSIGN,
				Value: "r",
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "1"},
					Right: &Node{
						Type:  NODE_OPERATOR,
						Value: "*",
						Left:  &Node{Type: NODE_NUMBER, Value: "2"},
						Right: &Node{Type: NODE_NUMBER, Value: "3"},
					},
				},
			},
		},
		{
			name:  "assignment with parentheses",
			input: "x=(2+3)*4",
			expected: &Node{
				Type:  NODE_ASSIGN,
				Value: "x",
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "*",
					Left: &Node{
						Type:  NODE_OPERATOR,
						Value: "+",
						Left:  &Node{Type: NODE_NUMBER, Value: "2"},
						Right: &Node{Type: NODE_NUMBER, Value: "3"},
					},
					Right: &Node{Type: NODE_NUMBER, Value: "4"},
				},
			},
		},

		// Variables/Identifiers
		{
			name:  "single variable",
			input: "x",
			expected: &Node{Type: NODE_IDENTIFIER, Value: "x"},
		},
		{
			name:  "variable in expression",
			input: "x+5",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_IDENTIFIER, Value: "x"},
				Right: &Node{Type: NODE_NUMBER, Value: "5"},
			},
		},
		{
			name:  "multiple variables",
			input: "x+y*z",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_IDENTIFIER, Value: "x"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "*",
					Left:  &Node{Type: NODE_IDENTIFIER, Value: "y"},
					Right: &Node{Type: NODE_IDENTIFIER, Value: "z"},
				},
			},
		},

		// Functions
		{
			name:  "function with one argument",
			input: "sin(30)",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "sin",
				Children: []*Node{
					{Type: NODE_NUMBER, Value: "30"},
				},
			},
		},
		{
			name:  "function with two arguments",
			input: "max(5,10)",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "max",
				Children: []*Node{
					{Type: NODE_NUMBER, Value: "5"},
					{Type: NODE_NUMBER, Value: "10"},
				},
			},
		},
		{
			name:  "function with expression argument",
			input: "sqrt(2+3)",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "sqrt",
				Children: []*Node{
					{
						Type:  NODE_OPERATOR,
						Value: "+",
						Left:  &Node{Type: NODE_NUMBER, Value: "2"},
						Right: &Node{Type: NODE_NUMBER, Value: "3"},
					},
				},
			},
		},
		{
			name:  "nested functions",
			input: "sin(cos(45))",
			expected: &Node{
				Type:  NODE_FUNCTION,
				Value: "sin",
				Children: []*Node{
					{
						Type:  NODE_FUNCTION,
						Value: "cos",
						Children: []*Node{
							{Type: NODE_NUMBER, Value: "45"},
						},
					},
				},
			},
		},
		{
			name:  "function in expression",
			input: "2+sin(30)",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_FUNCTION,
					Value: "sin",
					Children: []*Node{
						{Type: NODE_NUMBER, Value: "30"},
					},
				},
			},
		},

		// Complex expressions
		{
			name:  "complex expression",
			input: "2+3*4-5/2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "-",
				Left: &Node{
					Type:  NODE_OPERATOR,
					Value: "+",
					Left:  &Node{Type: NODE_NUMBER, Value: "2"},
					Right: &Node{
						Type:  NODE_OPERATOR,
						Value: "*",
						Left:  &Node{Type: NODE_NUMBER, Value: "3"},
						Right: &Node{Type: NODE_NUMBER, Value: "4"},
					},
				},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "/",
					Left:  &Node{Type: NODE_NUMBER, Value: "5"},
					Right: &Node{Type: NODE_NUMBER, Value: "2"},
				},
			},
		},
		{
			name:  "mixed precedence and associativity",
			input: "2+3*4^2",
			expected: &Node{
				Type:  NODE_OPERATOR,
				Value: "+",
				Left:  &Node{Type: NODE_NUMBER, Value: "2"},
				Right: &Node{
					Type:  NODE_OPERATOR,
					Value: "*",
					Left:  &Node{Type: NODE_NUMBER, Value: "3"},
					Right: &Node{
						Type:  NODE_OPERATOR,
						Value: "^",
						Left:  &Node{Type: NODE_NUMBER, Value: "4"},
						Right: &Node{Type: NODE_NUMBER, Value: "2"},
					},
				},
			},
		},

		// Edge cases
		{
			name:  "single number",
			input: "42",
			expected: &Node{Type: NODE_NUMBER, Value: "42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tt.input)
			assert.NoError(t, err)

			p := Parser{Tokens: tokens}
			ast, err := p.ParseExpression()

			assert.NoError(t, err, "Parser should not return error for valid input")
			assert.Equal(t, tt.expected, ast, "AST should match expected structure")
		})
	}
}

func TestParser_Expression_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "empty expression",
			input: "",
		},
		{
			name:  "unmatched opening parenthesis",
			input: "(2+3",
		},
		{
			name:  "missing operand after plus",
			input: "2+",
		},
		{
			name:  "missing operand after minus",
			input: "5-",
		},
		{
			name:  "missing operand after multiply",
			input: "3*",
		},
		{
			name:  "missing operand after divide",
			input: "8/",
		},
		{
			name:  "missing operand after exponent",
			input: "2^",
		},
		{
			name:  "just an operator",
			input: "+",
		},
		{
			name:  "operator at start",
			input: "*5",
		},
		{
			name:  "unclosed function parenthesis",
			input: "sin(30",
		},

		{
			name:  "assignment without value",
			input: "x=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				// Tokenizer caught the error, that's fine
				return
			}

			p := Parser{Tokens: tokens}
			ast, err := p.ParseExpression()

			assert.Error(t, err, "Expected error for input: "+tt.input)
			assert.Nil(t, ast, "AST should be nil on error")
		})
	}
}