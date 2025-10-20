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
	"fmt"
)

// NodeType categorizes AST node types for evaluation dispatch
type NodeType int

const (
	NODE_NUMBER     NodeType = iota // Terminal nodes containing numeric literals
	NODE_OPERATOR                   // Internal nodes representing operations
	NODE_FUNCTION                   // Function call nodes with argument lists
	NODE_ASSIGN                     // Variable assignment nodes
	NODE_IDENTIFIER                 // Variable and constant references
	NODE_OR
	NODE_AND
	NODE_COMPARISON
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
func (p *Parser) ParseExpression() (*Node, error) {
	if len(p.Tokens) == 0 {
		return nil, fmt.Errorf("empty expression")
	}

	node, err := p.parseAssignment()
	if err != nil {
		return nil, err
	}

	if p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		return nil, fmt.Errorf("unexpected token '%s' at position %d", tok.Value, p.pos)
	}

	return node, nil
}

func (p *Parser) parseAssignment() (*Node, error) {
	if p.pos+1 < len(p.Tokens) &&
		p.Tokens[p.pos].Type == tokenizer.IDENT &&
		p.Tokens[p.pos+1].Type == tokenizer.ASSIGN {

		varName := p.Tokens[p.pos].Value
		p.pos += 2

		rightNode, err := p.parseLogicalOr()
		if err != nil {
			return nil, err
		}
		if rightNode == nil {
			return nil, fmt.Errorf("expected expression after '='")
		}

		return &Node{
			Type:  NODE_ASSIGN,
			Value: varName,
			Right: rightNode,
		}, nil
	}

	return p.parseLogicalOr() // ← Changed from parseAddSub()
}

func (p *Parser) parseLogicalOr() (*Node, error) {
	node, err := p.parseLogicalAnd()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.LOGICAL && tok.Value == "||" {
			p.pos++
			rightNode, err := p.parseLogicalAnd()
			if err != nil {
				return nil, err
			}
			if rightNode == nil {
				return nil, fmt.Errorf("expected expression after '%s'", tok.Value)
			}
			node = &Node{
				Type:  NODE_OR,
				Value: "||",
				Left:  node,
				Right: rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseLogicalAnd() (*Node, error) {
	node, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.LOGICAL && tok.Value == "&&" {
			p.pos++
			rightNode, err := p.parseComparison()
			if err != nil {
				return nil, err
			}
			if rightNode == nil {
				return nil, fmt.Errorf("expected expression after '%s'", tok.Value)
			}
			node = &Node{
				Type:  NODE_AND,
				Value: "&&",
				Left:  node,
				Right: rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseComparison() (*Node, error) {
	node, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.COMPARISON {
			p.pos++
			rightNode, err := p.parseAddSub()
			if err != nil {
				return nil, err
			}
			if rightNode == nil {
				return nil, fmt.Errorf("expected expression after '%s'", tok.Value)
			}
			node = &Node{
				Type:  NODE_COMPARISON,
				Value: tok.Value,
				Left:  node,
				Right: rightNode,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseAddSub() (*Node, error) {
	node, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("expected expression")
	}

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && (tok.Value == "+" || tok.Value == "-") {
			p.pos++
			right, err := p.parseMulDiv()
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, fmt.Errorf("expected expression after '%s'", tok.Value)
			}
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseMulDiv() (*Node, error) {
	node, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("expected expression")
	}

	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && (tok.Value == "*" || tok.Value == "/") {
			p.pos++
			right, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, fmt.Errorf("expected expression after '%s'", tok.Value)
			}
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseUnary() (*Node, error) {
	if p.pos >= len(p.Tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	tok := p.Tokens[p.pos]
	if tok.Type == tokenizer.OPERATOR && (tok.Value == "-" || tok.Value == "+") {
		p.pos++
		child, err := p.parseExponent()
		if err != nil {
			return nil, err
		}
		if child == nil {
			return nil, fmt.Errorf("expected expression after unary '%s'", tok.Value)
		}
		if tok.Value == "-" {
			return &Node{
				Type:  NODE_OPERATOR,
				Value: "neg",
				Left:  child,
			}, nil
		}
		return child, nil
	}

	return p.parseExponent()
}

func (p *Parser) parseExponent() (*Node, error) {
	node, err := p.parsePostfix()
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("expected expression")
	}

	if p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && tok.Value == "^" {
			p.pos++
			// RIGHT ASSOCIATIVE: recursively call parseExponent
			right, err := p.parseExponent()
			if err != nil {
				return nil, err
			}
			if right == nil {
				return nil, fmt.Errorf("expected expression after '^'")
			}
			return &Node{Type: NODE_OPERATOR, Value: "^", Left: node, Right: right}, nil
		}
	}
	return node, nil
}

func (p *Parser) parsePostfix() (*Node, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("expected expression")
	}

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

	return node, nil
}

// parseFactor handles primary expressions
func (p *Parser) parseFactor() (*Node, error) {
	if p.pos >= len(p.Tokens) {
		return nil, fmt.Errorf("unexpected end of expression")
	}

	tok := p.Tokens[p.pos]
	p.pos++ // consume current token
	var node *Node

	switch tok.Type {
	case tokenizer.NUMBER:
		node = &Node{Type: NODE_NUMBER, Value: tok.Value}

	case tokenizer.IDENT:
		if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
			p.pos++ // consume '('
			var args []*Node

			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
				arg, err := p.parseAssignment()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)
				for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
					p.pos++ // consume ','
					arg, err = p.parseAssignment()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
				}
			}

			if p.pos >= len(p.Tokens) || p.Tokens[p.pos].Value != ")" {
				return nil, fmt.Errorf("unmatched opening parenthesis in function call '%s'", tok.Value)
			}
			p.pos++ // consume ')'

			node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
		} else {
			// normal identifier/variable
			node = &Node{Type: NODE_IDENTIFIER, Value: tok.Value}
		}

	case tokenizer.FUNCTION:
		if tok.Value == "!" {
			// Factorial is handled in postfix
			p.pos--
			return nil, fmt.Errorf("unexpected factorial operator")
		} else {
			if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
				p.pos++ // consume '('
				var args []*Node

				if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
					arg, err := p.parseAssignment()
					if err != nil {
						return nil, err
					}
					args = append(args, arg)
					for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
						p.pos++ // consume ','
						arg, err = p.parseAssignment()
						if err != nil {
							return nil, err
						}
						args = append(args, arg)
					}
				}

				if p.pos >= len(p.Tokens) || p.Tokens[p.pos].Value != ")" {
					return nil, fmt.Errorf("unmatched opening parenthesis in function '%s'", tok.Value)
				}
				p.pos++ // consume ')'

				node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
			} else {
				// Function without parentheses - treat as identifier
				node = &Node{Type: NODE_IDENTIFIER, Value: tok.Value}
			}
		}

	case tokenizer.PAREN:
		if tok.Value == "(" {
			subExpr, err := p.parseAssignment()
			if err != nil {
				return nil, err
			}
			if subExpr == nil {
				return nil, fmt.Errorf("empty parentheses")
			}
			if p.pos >= len(p.Tokens) || p.Tokens[p.pos].Value != ")" {
				return nil, fmt.Errorf("unmatched opening parenthesis")
			}
			p.pos++ // consume ')'
			node = subExpr
		} else if tok.Value == ")" {
			return nil, fmt.Errorf("unexpected closing parenthesis")
		}

	default:
		return nil, fmt.Errorf("unexpected token: %s", tok.Value)
	}

	return node, nil
}
