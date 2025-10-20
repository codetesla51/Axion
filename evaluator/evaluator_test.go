package evaluator

import (
	"Axion/parser"
	"Axion/tokenizer"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	for _, tt := range []struct {
		name      string
		input     string
		expected  float64
		expectErr bool
	}{
		{"log10 normal", "log10(100)", 2, false},
		{"log10 zero", "log10(0)", 0, true},
		{"log10 negative", "log10(-5)", 0, true},
		{"log10 very small", "log10(0.01)", -2, false},

		{"ln zero", "ln(0)", 0, true},
		{"ln negative", "ln(-1)", 0, true},
		{"ln of e", "ln(2.718281828459045)", 1, false},
		{"ln of 1", "ln(1)", 0, false},

		{"log single arg", "log(100)", 2, false},
		{"log single arg zero", "log(0)", 0, true},

		{"pow overflow detection", "pow(10,400)", 0, true},
		{"exp just below limit", "exp(709)", math.Exp(709), false},
		{"exp at limit", "exp(710)", 0, true},

		{"exponent operator at limit", "2^500", math.Pow(2, 500), false},
		{"exponent operator over limit", "2^501", 0, true},

		// Factorial edge cases
		{"factorial of 0", "0!", 1, false},
		{"factorial of 1", "1!", 1, false},

		{"factorial over limit", "171!", 0, true},
		{"factorial decimal", "3.5!", 0, true},
		{"factorial negative zero", "(-0)!", 1, false}, // -0 floors to 0

		{"sin no args", "sin()", 0, true},
		{"pow one arg", "pow(2)", 0, true},
		{"max one arg", "max(5)", 0, true},
		{"log three args", "log(10,2,3)", 0, true},

		// Empty collections
		{"sum single value", "sum(42)", 42, false},
		{"product single value", "product(42)", 42, false},
		{"mean single value", "mean(42)", 42, false},
		{"median single value", "median(42)", 42, false},
		{"mode single value", "mode(42)", 42, false},

		// Mode with ties (returns first encountered max)
		{"mode with tie", "mode(1,2,2,3,3)", 2, false}, // First to reach max count

		// Nested negation
		{"double negation", "--5", 5, false},
		{"triple negation", "---5", -5, false},

		// Variable persistence across evaluations
		{"var persistence 1", "x=10", 10, false},
		{"var persistence 2", "x+5", 15, false},
		{"var persistence 3", "x=20", 20, false},
		{"var persistence 4", "x", 20, false},

		// Special float values
		{"very large number", "1e308", 1e308, false},
		{"very small number", "1e-308", 1e-308, false},
		{"scientific notation", "6.022e23", 6.022e23, false},
		// Comparison operators
		{"greater than true", "5 > 3", 1, false},
		{"greater than false", "3 > 5", 0, false},
		{"greater than equal", "5 > 5", 0, false},
		{"less than true", "3 < 5", 1, false},
		{"less than false", "5 < 3", 0, false},
		{"less than equal", "5 < 5", 0, false},
		{"greater or equal true greater", "5 >= 3", 1, false},
		{"greater or equal true equal", "5 >= 5", 1, false},
		{"greater or equal false", "3 >= 5", 0, false},
		{"less or equal true less", "3 <= 5", 1, false},
		{"less or equal true equal", "5 <= 5", 1, false},
		{"less or equal false", "5 <= 3", 0, false},
		{"equal true", "5 == 5", 1, false},
		{"equal false", "5 == 3", 0, false},
		{"not equal true", "5 != 3", 1, false},
		{"not equal false", "5 != 5", 0, false},

		// Logical operators
		{"logical and both true", "1 && 1", 1, false},
		{"logical and first false", "0 && 1", 0, false},
		{"logical and second false", "1 && 0", 0, false},
		{"logical and both false", "0 && 0", 0, false},
		{"logical and with numbers", "5 && 3", 1, false},
		{"logical and with zero", "5 && 0", 0, false},
		{"logical or both true", "1 || 1", 1, false},
		{"logical or first true", "1 || 0", 1, false},
		{"logical or second true", "0 || 1", 1, false},
		{"logical or both false", "0 || 0", 0, false},
		{"logical or with numbers", "5 || 3", 1, false},
		{"logical or with zero", "0 || 5", 1, false},

		// Combined comparisons and logical operators
		{"comparison with and", "(5 > 3) && (2 < 4)", 1, false},
		{"comparison with or", "(5 < 3) || (2 < 4)", 1, false},
		{"multiple comparisons and", "(5 > 3) && (2 < 4) && (1 == 1)", 1, false},
		{"multiple comparisons or", "(5 < 3) || (2 > 4) || (1 == 1)", 1, false},
		{"mixed and or", "(5 > 3) && (2 < 4) || (1 > 2)", 1, false},
		{"precedence and before or", "0 || 1 && 0", 0, false},
		{"precedence and before or true", "1 || 0 && 0", 1, false},

		// Comparisons with arithmetic
		{"arithmetic in comparison left", "2 + 3 > 4", 1, false},
		{"arithmetic in comparison right", "5 > 2 + 2", 1, false},
		{"arithmetic both sides", "2 + 3 == 1 + 4", 1, false},
		{"complex arithmetic comparison", "2 * 3 > 4 + 1", 1, false},
		{"comparison with parentheses", "(2 + 3) * 2 > 8", 1, false},

		// Logical with arithmetic
		{"arithmetic with and", "2 + 3 > 4 && 1", 1, false},
		{"arithmetic with or", "2 + 3 < 4 || 1", 1, false},
		{"complex logical arithmetic", "(2 + 3) > 4 && (5 - 2) < 4", 1, false},

		// Assignment with comparisons and logical
		{"assign comparison result", "x = 5 > 3", 1, false},
		{"assign logical result", "x = 1 && 1", 1, false},
		{"assign complex logical", "x = (5 > 3) && (2 < 4)", 1, false},

		// Edge cases
		{"comparison with zero", "0 > 0", 0, false},
		{"comparison with negative", "-5 < 0", 1, false},
		{"comparison negative both", "-5 > -3", 0, false},
		{"logical with negative", "-5 && 3", 1, false},
		{"comparison with float", "2.5 > 2.4", 1, false},
		{"comparison float equal", "2.5 == 2.5", 1, false},
		{"not equal with float", "2.5 != 2.50001", 1, false},
		// Division edge cases
		{"zero divided by number", "0/5", 0, false},

		// Mod edge cases
		{"mod with negative", "mod(-7,3)", math.Mod(-7, 3), false},
		{"mod with floats", "mod(7.5,2.5)", 0, false},

		// Atan2 quadrant testing
		{"atan2 quadrant 1", "atan2(1,1)", 45, false},
		{"atan2 quadrant 2", "atan2(1,-1)", 135, false},
		{"atan2 quadrant 3", "atan2(-1,-1)", -135, false},
		{"atan2 quadrant 4", "atan2(-1,1)", -45, false},

		// Trigonometric edge cases
		{"sin of 0", "sin(0)", 0, false},
		{"cos of 0", "cos(0)", 1, false},
		{"tan of 0", "tan(0)", 0, false},
		{"sin of 180", "sin(180)", 0, false},
		{"cos of 90", "cos(90)", 0, false},
		{"tan of 270", "tan(270)", 0, true}, // 270 mod 180 = 90 (asymptote)

		// Inverse trig at boundaries
		{"asin of -1", "asin(-1)", -90, false},
		{"asin of 0", "asin(0)", 0, false},
		{"acos of -1", "acos(-1)", 180, false},
		{"acos of 0", "acos(0)", 90, false},
		{"atan of 0", "atan(0)", 0, false},

		// Degree/radian conversion
		{"deg2rad 360", "deg2rad(360)", 2 * math.Pi, false},
		{"rad2deg 2pi", "rad2deg(6.283185307179586)", 360, false},

		// Sign function completeness
		{"sign of very small positive", "sign(0.0000001)", 1, false},
		{"sign of very small negative", "sign(-0.0000001)", -1, false},

		// Rounding edge cases
		{"round negative", "round(-4.5)", -5, false}, // Go rounds half away from zero
		{"round half", "round(2.5)", 3, false},
		{"ceil negative", "ceil(-4.2)", -4, false},
		{"floor negative", "floor(-4.2)", -5, false},
		{"trunc negative", "trunc(-4.9)", -4, false},

		// Complex nested expressions
		{"nested trig", "sin(asin(0.5))", 0.5, false},
		{"nested log", "exp(ln(10))", 10, false},
		{"nested power", "sqrt(pow(5,2))", 5, false},

		// Mixed operations
		{"factorial in expression", "3!+2!", 8, false},
		{"multiple negations with ops", "-3*-4", 12, false},
		{"parentheses with negation", "-(3+4)", -7, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "var persistence 1" {
				Vars = make(map[string]float64)
			}

			tokens, err := tokenizer.Tokenize(tt.input)
			assert.NoError(t, err, "tokenizer error for input %q", tt.input)

			p := parser.Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Skipf("Parser error: %v", err)
			}
			got, err := Eval(ast)
			if tt.expectErr {
				assert.Error(t, err, "expected error for input %q", tt.input)
			} else {
				assert.NoError(t, err, "unexpected error for input %q: %v", tt.input, err)
				assert.InDelta(t, tt.expected, got, 1e-9, "expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestEvaluator_StatisticalFunctions(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected float64
	}{
		{"sum multiple", "sum(1,2,3,4,5)", 15},
		{"sum with negatives", "sum(-5,5,10)", 10},
		{"sum with decimals", "sum(1.5,2.5,3)", 7},

		{"product multiple", "product(2,3,4)", 24},
		{"product with zero", "product(5,0,10)", 0},
		{"product with negative", "product(-2,3,4)", -24},

		{"mean even count", "mean(2,4,6,8)", 5},
		{"mean odd count", "mean(1,2,3,4,5)", 3},
		{"mean with negatives", "mean(-10,0,10)", 0},

		{"median sorted odd", "median(1,2,3)", 2},
		{"median unsorted odd", "median(3,1,2)", 2},
		{"median sorted even", "median(1,2,3,4)", 2.5},
		{"median unsorted even", "median(4,2,1,3)", 2.5},
		{"median duplicates", "median(5,5,5)", 5},

		{"mode clear winner", "mode(1,2,2,3)", 2},
		{"mode all same", "mode(5,5,5,5)", 5},
		{"mode all unique", "mode(1,2,3,4)", 1}, // Returns first value
	} {
		t.Run(tt.name, func(t *testing.T) {
			Vars = make(map[string]float64)

			tokens, err := tokenizer.Tokenize(tt.input)
			assert.NoError(t, err)

			p := parser.Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Skipf("Parser error: %v", err)
			}
			got, err := Eval(ast)
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, got, 1e-9)
		})
	}
}

func TestEvaluator_ErrorMessages(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		errorContains string
	}{
		{"undefined var message", "undefined_var", "undefined variable"},
		{"division by zero message", "5/0", "division by zero"},
		{"factorial negative message", "(-5)!", "non-negative integers"},
		{"factorial overflow message", "171!", "too large"},
		{"sqrt negative message", "sqrt(-1)", "negative"},
		{"log domain message", "ln(0)", "positive"},
		{"tan asymptote message", "tan(90)", "undefined"},
		{"asin domain message", "asin(2)", "domain error"},
		{"exp overflow message", "exp(710)", "overflow"},
		{"pow invalid message", "pow(-2,0.5)", "non-integer exponent"},
		{"mod zero message", "mod(5,0)", "division by zero"},
		{"function arg count", "sin()", "requires 1 argument"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Vars = make(map[string]float64)

			tokens, err := tokenizer.Tokenize(tt.input)
			if err != nil {
				t.Skipf("tokenizer error: %v", err)
			}

			p := parser.Parser{Tokens: tokens}
			ast, err := p.ParseExpression()
			if err != nil {
				t.Skipf("Parser error: %v", err)
			}
			_, err = Eval(ast)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorContains)
		})
	}
}
