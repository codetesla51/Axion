package parser

import (
	"Axion/tokenizer"
	"reflect"
	"testing"
)

func TestParser_SimpleExpression(t *testing.T) {
	input := "(2 + (3 * 4))"
	tokens, err := tokenizer.Tokenize(input)
	if err != nil {
		t.Fatalf("tokenizer error: %v", err)
	}

	p := Parser{Tokens: tokens}
	ast := p.ParseExpression()

	expected := &Node{
		Type:  NODE_OPERATOR,
		Value: "+",
		Left:  &Node{Type: NODE_NUMBER, Value: "2"},
		Right: &Node{
			Type:  NODE_OPERATOR,
			Value: "*",
			Left:  &Node{Type: NODE_NUMBER, Value: "3"},
			Right: &Node{Type: NODE_NUMBER, Value: "4"},
		},
	}

	if !reflect.DeepEqual(ast, expected) {
		t.Errorf("expected %+v, got %+v", expected, ast)
	}
}
