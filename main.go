package main

import (
	"fmt"
	"Axion/parser"
	"Axion/tokenizer"
	"Axion/evaluator"
)

func main() {
	input := "3.5 + 4 * 5 + sin(30) / (40 * 60)"
	tokens := tokenizer.Tokenize(input)
	p := parser.Parser{Tokens: tokens}
	ast := p.ParseExpression()
	result := evaluator.Eval(ast)
	fmt.Println("Result:", result)
}