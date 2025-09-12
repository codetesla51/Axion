/*
Parser Module - Syntax Analysis and AST Construction
====================================================
Part of Axion CLI Calculator
Author: Uthman
Year: 2025

This module implements recursive descent parsing for mathematical expressions,
converting token sequences from the tokenizer into Abstract Syntax Trees (AST)
ready for evaluation. The parser enforces proper mathematical operator precedence
and handles complex nested expressions.

Architecture:
The parser uses recursive descent methodology with precedence climbing to
ensure correct mathematical evaluation order. Each precedence level is
handled by dedicated parsing functions that build appropriate AST structures.

Precedence Hierarchy (lowest to highest):
1. Assignment operators (=)
2. Addition and subtraction (+, -)
3. Multiplication and division (*, /)
4. Unary operators (-, +)
5. Exponentiation (^) - right associative
6. Postfix operators (factorial !)
7. Primary expressions (numbers, functions, parentheses)

AST Node Types:
- NODE_NUMBER: Terminal nodes containing numeric literals
- NODE_OPERATOR: Binary and unary operation nodes
- NODE_FUNCTION: Function calls with argument lists
- NODE_ASSIGN: Variable assignment operations
- NODE_IDENTIFIER: Variable and constant references

Key Features:
- Operator Precedence: Ensures mathematical correctness (2 + 3 * 4 = 14, not 20)
- Right Associativity: Proper handling of exponentiation (2^3^2 = 2^(3^2) = 512)
- Function Parsing: Multi-argument function support with comma separation
- Assignment Support: Variable assignment with proper precedence
- Error Recovery: Graceful handling of malformed expressions
- Memory Efficiency: Minimal AST node allocation

Expression Examples:
- "2 + 3 * 4" → AST with multiplication evaluated before addition
- "sin(30) + cos(60)" → Function calls with numeric arguments
- "x = 5 + 3" → Assignment node with expression evaluation
- "2^3^2" → Right-associative exponentiation chain

The parser bridges the gap between lexical tokens and evaluable expressions,
ensuring mathematical correctness throughout the parsing process.
*/
package parser

import (
	"Axion/tokenizer"
)

// NodeType categorizes AST node types for evaluation dispatch
type NodeType int

const (
	NODE_NUMBER   NodeType = iota // Terminal nodes containing numeric literals
	NODE_OPERATOR                 // Internal nodes representing operations
	NODE_FUNCTION                 // Function call nodes with argument lists
	NODE_ASSIGN
	NODE_IDENTIFIER
)

// Node represents a single node in the Abstract Syntax Tree
type Node struct {
	Type     NodeType // Node classification for evaluation
	Value    string   // Node content (number, operator symbol, function name)
	Left     *Node    // Left operand for binary operators
	Right    *Node    // Right operand for binary operators
	Children []*Node  // Argument list for function calls
}

// Parser maintains parsing state during recursive descent
type Parser struct {
	Tokens []tokenizer.Token // Input token sequence
	pos    int               // Current parsing position
}

// ParseExpression initiates parsing at the lowest precedence level
func (p *Parser) ParseExpression() *Node {
	return p.parseAssignment()
}

func (p *Parser) parseAssignment() *Node {
	if p.pos+1 < len(p.Tokens) &&
		p.Tokens[p.pos].Type == tokenizer.IDENT &&
		p.Tokens[p.pos+1].Type == tokenizer.ASSIGN {

		varName := p.Tokens[p.pos].Value
		p.pos += 2 // skip IDENT and ASSIGN

		rightNode := p.parseAddSub()

		return &Node{
			Type:  NODE_ASSIGN,
			Value: varName,
			Right: rightNode,
		}
	}

	return p.parseAddSub()
}

func (p *Parser) parseAddSub() *Node {
	node := p.parseMulDiv()

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && (tok.Value == "+" || tok.Value == "-") {
			p.pos++
			right := p.parseMulDiv()
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parseMulDiv() *Node {
	node := p.parseUnary()

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && (tok.Value == "*" || tok.Value == "/") {
			p.pos++
			right := p.parseUnary()
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parseUnary() *Node {
	if p.pos >= len(p.Tokens) {
		return nil
	}

	tok := p.Tokens[p.pos]
	if tok.Type == tokenizer.OPERATOR && (tok.Value == "-" || tok.Value == "+") {
		p.pos++
		child := p.parseExponent()
		if tok.Value == "-" {
			return &Node{
				Type:  NODE_OPERATOR,
				Value: "neg",
				Left:  child,
			}
		}
		return child
	}

	return p.parseExponent()
}

func (p *Parser) parseExponent() *Node {
	node := p.parsePostfix()

	if p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && tok.Value == "^" {
			p.pos++
			right := p.parseUnary()
			return &Node{Type: NODE_OPERATOR, Value: "^", Left: node, Right: right}
		}
	}
	return node
}

func (p *Parser) parsePostfix() *Node {
	node := p.parseFactor()

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.FUNCTION && tok.Value == "!" {
			p.pos++
			node = &Node{
				Type:     NODE_FUNCTION,
				Value:    "!",
				Children: []*Node{node},
			}
		} else {
			break
		}
	}

	return node
}

// parseFactor handles primary expressions
func (p *Parser) parseFactor() *Node {
	if p.pos >= len(p.Tokens) {
		return nil
	}

	tok := p.Tokens[p.pos]
	p.pos++ // consume current token
	var node *Node

	switch tok.Type {
	case tokenizer.NUMBER:
		node = &Node{Type: NODE_NUMBER, Value: tok.Value}

	case tokenizer.IDENT:
		// Check if next token is '(' → function call
		if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
			p.pos++ // consume '('
			var args []*Node

			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
				arg := p.ParseExpression()
				args = append(args, arg)
				for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
					p.pos++ // consume ','
					arg = p.ParseExpression()
					args = append(args, arg)
				}
			}

			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == ")" {
				p.pos++ // consume ')'
			}

			node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
		} else {
			// normal identifier/variable
			node = &Node{Type: NODE_IDENTIFIER, Value: tok.Value}
		}

	case tokenizer.FUNCTION:
		if tok.Value == "!" {
			// Factorial is handled in postfix
			p.pos--
			return nil
		} else {
			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
				p.pos++ // consume '('
				var args []*Node

				if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
					arg := p.ParseExpression()
					args = append(args, arg)
					for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
						p.pos++ // consume ','
						arg = p.ParseExpression()
						args = append(args, arg)
					}
				}

				if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == ")" {
					p.pos++ // consume ')'
				}

				node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
			} else {
				// Function without parentheses - treat as identifier
				node = &Node{Type: NODE_IDENTIFIER, Value: tok.Value}
			}
		}

	case tokenizer.PAREN:
		if tok.Value == "(" {
			node = p.ParseExpression()
			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == ")" {
				p.pos++
			}
		}
	}

	return node
}
