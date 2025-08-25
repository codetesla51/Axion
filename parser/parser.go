package parser

import (
	"Axion/tokenizer"
)

type NodeType int

const (
	NODE_NUMBER NodeType = iota
	NODE_OPERATOR
	NODE_FUNCTION
)

type Node struct {
	Type     NodeType
	Value    string
	Left     *Node
	Right    *Node
	Children []*Node
}

type Parser struct {
	Tokens []tokenizer.Token
	pos    int
}

func (p *Parser) ParseExpression() *Node {
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
func (p *Parser) parseExponent() *Node {
	node := p.parseFactor() // left-hand side
	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && tok.Value == "^" {
			p.pos++
			right := p.parseFactor() // right-hand side
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parseMulDiv() *Node {
	node := p.parseExponent()
	for p.pos < len(p.Tokens) {
		tok := p.Tokens[p.pos]
		if tok.Type == tokenizer.OPERATOR && (tok.Value == "*" || tok.Value == "/") {
			p.pos++
			right := p.parseExponent()
			node = &Node{Type: NODE_OPERATOR, Value: tok.Value, Left: node, Right: right}
		} else {
			break
		}
	}
	return node
}

func (p *Parser) parseFactor() *Node {
	if p.pos >= len(p.Tokens) {
		return nil
	}
	tok := p.Tokens[p.pos]
	p.pos++

	var node *Node

	switch tok.Type {
	case tokenizer.NUMBER:
		node = &Node{Type: NODE_NUMBER, Value: tok.Value}
	case tokenizer.FUNCTION:
    if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "(" {
        p.pos++ // skip '('
        
        var args []*Node
        
        // Parse first argument
        if p.pos < len(p.Tokens) && p.Tokens[p.pos].Value != ")" {
            arg := p.ParseExpression()
            args = append(args, arg)
            
            // Parse additional arguments separated by commas
            for p.pos < len(p.Tokens) && p.Tokens[p.pos].Value == "," {
                p.pos++ // skip ','
                arg = p.ParseExpression()
                args = append(args, arg)
            }
        }
        
        p.pos++ // skip ')'
        node = &Node{Type: NODE_FUNCTION, Value: tok.Value, Children: args}
    } else {
        node = &Node{Type: NODE_FUNCTION, Value: tok.Value}
    }
	case tokenizer.PAREN:
		if tok.Value == "(" {
			node = p.ParseExpression()
			p.pos++ // skip ')'
		}
	}
	if p.pos < len(p.Tokens) && p.Tokens[p.pos].Type == tokenizer.FUNCTION && p.Tokens[p.pos].Value == "!" {
		p.pos++
		node = &Node{
			Type:     NODE_FUNCTION,
			Value:    "!",
			Children: []*Node{node},
		}
	}

	return node
}
