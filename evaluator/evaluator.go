package evaluator

import(
  	"math"
	"strconv"
		"Axion/parser"

  )

func Eval(node *parser.Node) float64 {
	if node == nil {
		return 0
	}
	switch node.Type {
	case parser.NODE_NUMBER:
		val, _ := strconv.ParseFloat(node.Value, 64)
		return val
	case parser.NODE_OPERATOR:
		left := Eval(node.Left)
		right := Eval(node.Right)
		switch node.Value {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			return left / right
		case "^":
			return math.Pow(left, right)
		}
	case parser.NODE_FUNCTION:
		arg := Eval(node.Children[0])
		switch node.Value {
		case "sin":
			return math.Sin(arg * math.Pi / 180) 
		case "cos":
			return math.Cos(arg * math.Pi / 180)
		case "tan":
			return math.Tan(arg * math.Pi / 180)
		}
	}
	return 0
}
