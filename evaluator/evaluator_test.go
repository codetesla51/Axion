package evaluator

import (
	"Axion/parser"
	"Axion/tokenizer"
	"math"
	"testing"
)

func TestParser_SimpleExpression(t *testing.T) {
	input := "0.1 + 0.2"
	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		t.Fatalf("tokenizer error: %v", err)
	}
	p := parser.Parser{Tokens: tokens}
	ast := p.ParseExpression()
	got, err := Eval(ast)
	if err != nil {
		t.Fatalf("tokenizer error: %v", err)
		return
	}
	expected := 0.3
	if math.Abs(got-expected) > 1e-9 {
		t.Errorf("expected %v, got %v", expected, got)
	}

}
