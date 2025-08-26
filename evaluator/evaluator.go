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
	// Unary minus
	if node.Value == "neg" {
		left, err := Eval(node.Left)
		if err != nil {
			return 0, err
		}
		return -left, nil
	}
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
		case "sin", "cos", "tan":

			radians := arg1 * math.Pi / 180
			switch node.Value {
			case "sin":
				return math.Sin(radians), nil
			case "cos":
				return math.Cos(radians), nil
			case "tan":
				if math.Mod(arg1, 180) == 90 {
					return 0, fmt.Errorf("tan(%g°): undefined (asymptote)", arg1)
				}
				return math.Tan(radians), nil
			}

		case "asin":
			if arg1 < -1 || arg1 > 1 {
				return 0, fmt.Errorf("asin: domain error - input must be between -1 and 1, got %g", arg1)
			}
			return math.Asin(arg1) * 180 / math.Pi, nil

		case "acos":
			if arg1 < -1 || arg1 > 1 {
				return 0, fmt.Errorf("acos: domain error - input must be between -1 and 1, got %g", arg1)
			}
			return math.Acos(arg1) * 180 / math.Pi, nil

		case "atan":
			return math.Atan(arg1) * 180 / math.Pi, nil

		case "log":
			if arg1 <= 0 {
				return 0, fmt.Errorf("log: domain error - logarithm undefined for non-positive numbers, got %g", arg1)
			}
			return math.Log(arg1), nil

		case "log10":
			if arg1 <= 0 {
				return 0, fmt.Errorf("log10: domain error - logarithm undefined for non-positive numbers, got %g", arg1)
			}
			return math.Log10(arg1), nil

		case "sqrt":
			if arg1 < 0 {
				return 0, fmt.Errorf("sqrt: domain error - cannot take square root of negative number, got %g", arg1)
			}
			return math.Sqrt(arg1), nil

		case "exp":
			if arg1 > 709 {
				return 0, fmt.Errorf("exp: overflow error - exp(%g) would exceed maximum representable value", arg1)
			}
			return math.Exp(arg1), nil

		case "abs", "ceil", "floor":

			switch node.Value {
			case "abs":
				return math.Abs(arg1), nil
			case "ceil":
				return math.Ceil(arg1), nil
			case "floor":
				return math.Floor(arg1), nil
			}

		case "!":
			val, err := factorial(arg1)
			if err != nil {
				return 0, err
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
				if arg1 == 0 && arg2 < 0 {
					return 0, fmt.Errorf("pow: domain error - 0 raised to negative power is undefined")
				}
				if arg1 < 0 && arg2 != math.Floor(arg2) {
					return 0, fmt.Errorf("pow: domain error - negative base with non-integer exponent")
				}
				result := math.Pow(arg1, arg2)
				if math.IsInf(result, 0) {
					return 0, fmt.Errorf("pow: overflow error - pow(%g, %g) exceeds representable range", arg1, arg2)
				}
				return result, nil
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
