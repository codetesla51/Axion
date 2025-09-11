/*
Parser Module - Recursive Descent Parser
=========================================

This module implements a recursive descent parser for mathematical expressions.
It constructs an Abstract Syntax Tree (AST) from the token sequence produced
by the tokenizer, ensuring proper operator precedence and associativity.

Operator Precedence (highest to lowest):
1. Primary expressions (numbers, functions, parentheses)
2. Postfix operators (factorial !)
3. Exponentiation (^) - right associative
4. Unary operators (+, -) 
5. Multiplication and division (*, /) - left associative
6. Addition and subtraction (+, -) - left associative

The parser uses recursive descent with each precedence level implemented
as a separate function. Higher precedence operations are parsed by calling
functions that handle lower-precedence levels.

Right associativity for exponentiation means: 2^3^2 = 2^(3^2) = 512
Left associativity for other operators means: 8/4/2 = (8/4)/2 = 1
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
        Type     NodeType  // Node classification for evaluation
        Value    string    // Node content (number, operator symbol, function name)
        Left     *Node     // Left operand for binary operators
        Right    *Node     // Right operand for binary operators
        Children []*Node   // Argument list for function calls
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

// parseAddSub handles addition and subtraction (precedence level 1 - lowest)
// Implements left associativity: a - b + c = ((a - b) + c)
func (p *Parser) parseAssignment() *Node {
    if p.pos+1 < len(p.Tokens) &&
       p.Tokens[p.pos].Type == tokenizer.IDENT &&
       p.Tokens[p.pos+1].Type == tokenizer.ASSIGN {

        varName := p.Tokens[p.pos].Value
        p.pos += 2 // skip IDENT and ASSIGN

        rightNode := p.parseAddSub() // still a Node

        return &Node{
            Type:  NODE_ASSIGN, 
            Value: varName,       // store variable name in Value
            Right: rightNode,     // RHS of assignment
        }
    }

    return p.parseAddSub()
}
func (p *Parser) parseAddSub() *Node {
        // Parse left operand at higher precedence level
        node := p.parseMulDiv()

        // Continue parsing same-level operators with left associativity
        for p.pos < len(p.Tokens) {
                tok := p.Tokens[p.pos]
                if tok.Type == tokenizer.OPERATOR && (tok.Value == "+" || tok.Value == "-") {
                        p.pos++ // Consume operator token
                        right := p.parseMulDiv() // Parse right operand
                        // Create binary operator node with left associativity
                        node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
                } else {
                        break // No more operators at this precedence level
                }
        }
        return node
}

// parseMulDiv handles multiplication and division (precedence level 2)
// Implements left associativity: a / b * c = ((a / b) * c)
func (p *Parser) parseMulDiv() *Node {
        // Parse left operand at higher precedence level
        node := p.parseUnary()

        // Continue parsing same-level operators with left associativity
        for p.pos < len(p.Tokens) {
                tok := p.Tokens[p.pos]
                if tok.Type == tokenizer.OPERATOR && (tok.Value == "*" || tok.Value == "/") {
                        p.pos++ // Consume operator token
                        right := p.parseUnary() // Parse right operand
                        // Create binary operator node with left associativity
                        node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
                } else {
                        break // No more operators at this precedence level
                }
        }
        return node
}

// parseUnary handles unary plus and minus operators (precedence level 3)
// Unary operators have lower precedence than exponentiation:
// -3^2 should parse as -(3^2) = -9, not (-3)^2 = 9
func (p *Parser) parseUnary() *Node {
        // Check for end of input
        if p.pos >= len(p.Tokens) {
                return nil
        }

        tok := p.Tokens[p.pos]
        if tok.Type == tokenizer.OPERATOR && (tok.Value == "-" || tok.Value == "+") {
                p.pos++ // Consume unary operator
                // Parse operand at higher precedence (exponentiation comes first)
                child := p.parseExponent()
                if tok.Value == "-" {
                        // Create unary negation node
                        return &Node{
                                Type:  NODE_OPERATOR,
                                Value: "neg", // Internal representation for unary minus
                                Left:  child,
                        }
                }
                // Unary plus is effectively a no-op
                return child
        }

        // No unary operator present, continue to higher precedence
        return p.parseExponent()
}

// parseExponent handles exponentiation (precedence level 4)
// Implements right associativity: 2^3^4 = 2^(3^4) = 2^81 = large number
func (p *Parser) parseExponent() *Node {
        // Parse base operand at higher precedence level
        node := p.parsePostfix()

        // Check for exponentiation operator
        if p.pos < len(p.Tokens) {
                tok := p.Tokens[p.pos]
                if tok.Type == tokenizer.OPERATOR && tok.Value == "^" {
                        p.pos++ // Consume exponentiation operator
                        // Right associativity: parse exponent at same precedence level
                        right := p.parseUnary()
                        return &Node{Type: NODE_OPERATOR, Value: "^", Left: node, Right: right}
                }
        }
        return node
}

// parsePostfix handles postfix operators like factorial (precedence level 5)
// Multiple factorials can be applied: 5!! = (5!)! = 120! (mathematically valid)
func (p *Parser) parsePostfix() *Node {
        // Parse primary expression first
        node := p.parseFactor()

        // Apply all consecutive postfix operators
        for p.pos < len(p.Tokens) {
                tok := p.Tokens[p.pos]
                if tok.Type == tokenizer.FUNCTION && tok.Value == "!" {
                        p.pos++ // Consume factorial operator
                        // Wrap current node as argument to factorial function
                        node = &Node{
                                Type:     NODE_FUNCTION,
                                Value:    "!",
                                Children: []*Node{node},
                        }
                } else {
                        break // No more postfix operators
                }
        }

        return node
}

// parseFactor handles primary expressions (precedence level 6 - highest)
// Primary expressions: numbers, function calls, parenthesized expressions
func (p *Parser) parseFactor() *Node {
        // Check for end of input
        if p.pos >= len(p.Tokens) {
                return nil
        }

        tok := p.Tokens[p.pos]
        p.pos++ // Consume current token

        var node *Node

        switch tok.Type {
        // Handle numeric literals
        case tokenizer.NUMBER:
                node = &Node{Type: NODE_NUMBER, Value: tok.Value}
case tokenizer.IDENT:
    node = &Node{
        Type:  NODE_IDENTIFIER,
        Value: tok.Value,
    }
        // Handle function calls
        case tokenizer.FUNCTION:
                // Special case: factorial should not be parsed as prefix function
                if tok.Value == "!" {
                        p.pos-- // Backtrack - factorial is handled in postfix parsing
                        return nil
                }

                // Check for function argument list in parentheses
                if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
                        p.pos++ // Consume opening parenthesis
                        var args []*Node

                        // Parse function arguments
                        if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
                                // Parse first argument as full expression
                                arg := p.ParseExpression()
                                args = append(args, arg)

                                // Parse additional arguments separated by commas
                                for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
                                        p.pos++ // Consume comma separator
                                        arg = p.ParseExpression()
                                        args = append(args, arg)
                                }
                        }

                        // Consume closing parenthesis if present
                        if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == ")" {
                                p.pos++ // Consume closing parenthesis
                        }
                        // Create function call node with argument list
                        node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
                } else {
                        // Function without parentheses - create node with empty argument list
                        node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: []*Node{}}
                }

        // Handle parenthesized expressions
        case tokenizer.PAREN:
                if tok.Value == "(" {
                        // Parse subexpression within parentheses
                        node = p.ParseExpression()
                        // Consume matching closing parenthesis
                        if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == ")" {
                                p.pos++ // Consume closing parenthesis
                        }
                }
        }

        return node
}