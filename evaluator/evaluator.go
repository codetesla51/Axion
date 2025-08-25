package evaluator

import (
	"Axion/parser"
	"fmt"
	"math"
	"strconv"
)

func factorial(n float64) (float64, error) {
	if n < 0 || n != math.Floor(n) {
		return 0, fmt.Errorf("factorial only defined for non-negative integers")
	}
	// IEEE 754 double precision can represent up to approximately 170!
	// 171! ≈ 1.24e+309, which exceeds the maximum finite float64 value
	if n > 170 {
		return 0, fmt.Errorf("factorial too large: %g! exceeds maximum representable value (limit: 170!)", n)
	}
	result := 1.0
	for i := 2; i <= int(n); i++ {
		result *= float64(i)
	}
	return result, nil
}

func Eval(node *parser.Node) (float64, error) {
	if node == nil {
		return 0, fmt.Errorf("Invalid")
	}
	switch node.Type {
	case parser.NODE_NUMBER:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q", node.Value)
		}
		return val, nil

	case parser.NODE_OPERATOR:
		left, err := Eval(node.Left)
		if err != nil {
			return 0, err
		}
		right, err := Eval(node.Right)
		if err != nil {
			return 0, err
		}

		switch node.Value {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return left / right, nil
		case "^":
			return math.Pow(left, right), nil
		default:
			return 0, fmt.Errorf("unknown operator %q", node.Value)
		}

	case parser.NODE_FUNCTION:
		if len(node.Children) < 1 {
			return 0, fmt.Errorf("function %q requires at least 1 argument", node.Value)
		}
		arg1, err := Eval(node.Children[0])
		if err != nil {
			return 0, err
		}

		switch node.Value {
		case "sin":
			return math.Sin(arg1 * math.Pi / 180), nil
		case "cos":
			return math.Cos(arg1 * math.Pi / 180), nil
		case "tan":
			return math.Tan(arg1 * math.Pi / 180), nil
		case "asin":
			return math.Asin(arg1) * 180 / math.Pi, nil
		case "acos":
			return math.Acos(arg1) * 180 / math.Pi, nil
		case "atan":
			return math.Atan(arg1) * 180 / math.Pi, nil
		case "log":
			return math.Log(arg1), nil
		case "log10":
			return math.Log10(arg1), nil
		case "sqrt":
			return math.Sqrt(arg1), nil
		case "exp":
			return math.Exp(arg1), nil
		case "abs":
			return math.Abs(arg1), nil
		case "ceil":
			return math.Ceil(arg1), nil
		case "floor":
			return math.Floor(arg1), nil
		case "!":
			val, err := factorial(arg1)
			if err != nil {
				return 0, nil
			}
			return val, nil
		case "pow", "max", "min":
			if len(node.Children) < 2 {
				return 0, fmt.Errorf("function %q requires 2 arguments", node.Value)
			}
			arg2, err := Eval(node.Children[1])
			if err != nil {
				return 0, err
			}
			switch node.Value {
			case "pow":
				return math.Pow(arg1, arg2), nil
			case "max":
				return math.Max(arg1, arg2), nil
			case "min":
				return math.Min(arg1, arg2), nil
			}
		default:
			return 0, fmt.Errorf("unknown function %q", node.Value)
		}
	}

	return 0, fmt.Errorf("invalid node type")
}
