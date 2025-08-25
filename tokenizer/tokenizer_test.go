package tokenizer

import (
	"reflect"
	"testing"
)

func TestTokenizer(t *testing.T) {
	input := "1 + 2"
	got, err := Tokenize(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []Token{
		{Type: 0, Value: "1"},
		{Type: 1, Value: "+"},
		{Type: 0, Value: "2"},
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
