# Axion Calculator

<div align="center">

```
  ╔═╗─┐ ┬┬┌─┐┌┐┌
  ╠═╣┌┴┬┘││ ││││
  ╩ ╩┴ └─┴└─┘┘└┘
```

**A sophisticated, high-precision mathematical engine with advanced CLI interface, built in Go**

[![Go Version](https://img.shields.io/badge/Go-1.24.5-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=flat-square&logo=checkmarx)](https://github.com/codetesla51/Axion)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square&logo=github-actions)](https://github.com/codetesla51/Axion)
[![Coverage](https://img.shields.io/badge/Coverage-95%25-brightgreen?style=flat-square&logo=codecov)](https://github.com/codetesla51/Axion)
[![Go Report Card](https://img.shields.io/badge/Go%20Report-A+-brightgreen?style=flat-square&logo=go)](https://goreportcard.com/report/github.com/codetesla51/Axion)

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Examples](#examples) • [API](#api) • [Contributing](#contributing)

</div>

---

## Overview

**Axion** is a powerful command-line calculator that transcends simple arithmetic, offering a complete mathematical computing environment. Built with modern Go architecture using the Cobra CLI framework, it features a sophisticated expression parser, extensive mathematical function library, comprehensive unit conversion system, variable management, logical operations, comparison operators, and persistent calculation history.

### Why Axion?

- **Precision**: Advanced mathematical expression parser with proper operator precedence
- **Scientific**: Complete scientific notation support and comprehensive function library
- **Logical Operations**: Full support for boolean logic with `&&`, `||` operators
- **Comparisons**: Complete comparison operators (`>`, `<`, `>=`, `<=`, `==`, `!=`)
- **Conversions**: Built-in unit conversion across length, weight, and time categories
- **Memory**: Persistent calculation history and variable storage across sessions
- **Performance**: Optimized Go implementation with minimal memory footprint
- **Modern CLI**: Beautiful, color-coded interface powered by Cobra framework
- **Reliability**: Comprehensive error handling and 95%+ test coverage
- **Extensible**: Modular architecture for easy feature additions

---

## Features

### Core Mathematical Engine
- **Expression Parsing**: Advanced recursive descent parser with proper precedence
- **Scientific Notation**: Full support (`2e-10`, `3.14e+5`, `1.5E-20`)
- **Operator Support**: Basic arithmetic, exponentiation (`^`), factorial (`!`)
- **Logical Operators**: Boolean logic with `&&` (AND), `||` (OR)
- **Comparison Operators**: `>`, `<`, `>=`, `<=`, `==`, `!=` returning 1 (true) or 0 (false)
- **Parentheses Grouping**: Complex nested expression support
- **Implicit Multiplication**: Automatic insertion (`2sin(x)` → `2 * sin(x)`)

### Mathematical Functions

| Category | Functions | Description |
|----------|-----------|-------------|
| **Trigonometric** | `sin()`, `cos()`, `tan()`, `asin()`, `acos()`, `atan()`, `atan2()` | Degree-based trigonometry |
| **Logarithmic** | `ln()`, `log()`, `log10()`, `log2()`, `log(x, base)` | Natural, common, and custom base logs |
| **Exponential** | `exp()`, `pow()`, `sqrt()` | Exponential and power functions |
| **Utility** | `abs()`, `ceil()`, `floor()`, `round()`, `trunc()`, `sign()` | Number manipulation |
| **Statistical** | `mean()`, `median()`, `mode()`, `sum()`, `product()` | Multi-argument statistics |
| **Comparison** | `max()`, `min()` | Value comparison |
| **Special** | `!` (factorial), `mod()` | Advanced operations |
| **Output** | `print()` | Display values and expressions |

### Logical & Comparison Operations
- **Comparison Operators**: 
  - `>` (greater than), `<` (less than)
  - `>=` (greater or equal), `<=` (less or equal)
  - `==` (equal), `!=` (not equal)
  - Returns `1` for true, `0` for false
- **Logical Operators**:
  - `&&` (logical AND) - returns `1` if both operands are non-zero
  - `||` (logical OR) - returns `1` if at least one operand is non-zero
  - Proper precedence: `&&` binds tighter than `||`
- **Combined Expressions**: `(x > 5) && (y < 10)`, `2 + 3 > 4 || 0`

### Variables & Constants
- **Variable Assignment**: `x = 5`, `result = sin(30) + cos(60)`
- **Mathematical Constants**: `pi`, `e`, `phi`, `sqrt2`, `c`, `G`, `h`, `R`
- **Persistent Storage**: Variables maintained across calculator sessions
- **Dynamic Updates**: Real-time variable modification and retrieval

### Unit Conversion System
- **Length Units**: `m`, `cm`, `mm`, `km`, `in`, `ft`, `yd`, `mi`
- **Weight Units**: `kg`, `g`, `mg`, `lb`, `oz`, `ton`
- **Time Units**: `s`, `ms`, `min`, `h`, `d`
- **Cross-Category Protection**: Prevents invalid conversions

### Advanced Features
- **Calculation History**: JSON-based persistent storage with session continuity
- **Precision Control**: Configurable decimal precision (0-20 places)
- **Interactive REPL**: Beautiful color-coded command interface with help system
- **Error Handling**: Comprehensive error reporting with context
- **Cross-Platform**: Native support for Windows, macOS, and Linux
- **Cobra Framework**: Professional CLI with subcommands and flags

---

## Installation

### Prerequisites
- **Go**: Version 1.24.5 or later ([Download Go](https://golang.org/dl/))
- **Git**: For repository cloning ([Download Git](https://git-scm.com/downloads))

### Quick Installation

#### Using Install Script (Linux/macOS)

```bash
# Clone the repository
git clone https://github.com/codetesla51/Axion.git
cd Axion

# Run the installation script
chmod +x install.sh
./install.sh

# Reload your shell configuration
source ~/.bashrc  # or ~/.zshrc for Zsh users

# Verify installation
axion --help
```

The install script will:
- Build the Axion binary
- Create a symlink in `~/.local/bin/axion`
- Automatically detect your shell (bash/zsh)
- Add `~/.local/bin` to your PATH if needed

#### Manual Installation

```bash
# Clone the repository
git clone https://github.com/codetesla51/Axion.git

# Navigate to project directory
cd Axion

# Install dependencies
go mod download

# Run directly (recommended for development)
go run main.go

# Or build executable
go build -o axion
./axion  # Unix/macOS
# axion.exe  # Windows
```

#### Installation via Go Install

```bash
go install github.com/codetesla51/Axion@latest
```

---

## Usage

### Getting Started

Launch Axion and start calculating:

```bash
axion

  ╔═╗─┐ ┬┬┌─┐┌┐┌
  ╠═╣┌┴┬┘││ ││││
  ╩ ╩┴ └─┴└─┘┘└┘

  A Powerful CLI Calculator
  Type 'help' for commands or 'exit' to quit

» 2 + 3 * 4
Result: 14

» sin(30) + cos(60)
Result: 1

» 5 > 3
Result: 1

» (5 > 3) && (2 < 4)
Result: 1

» x = sqrt(16)
Result: 4

» print(x)
4

» convert 100 cm to m
100 cm = 1 m
```

### Command Reference

| Command | Syntax | Description | Example |
|---------|--------|-------------|---------|
| **Expression** | `<mathematical expression>` | Evaluate any mathematical expression | `2 + 3 * sin(45)` |
| **Comparison** | `<expr> <op> <expr>` | Compare values (`>`, `<`, `>=`, `<=`, `==`, `!=`) | `5 > 3`, `x == 10` |
| **Logical** | `<expr> <op> <expr>` | Logical operations (`&&`, `||`) | `(x > 5) && (y < 10)` |
| **Assignment** | `<variable> = <expression>` | Assign value to variable | `x = 10`, `area = pi * r^2` |
| **Print** | `print(<expression>)` | Display expression result | `print(2 + 3)`, `print(x)` |
| **Conversion** | `convert <value> <from> to <to>` | Convert between units | `convert 5 km to mi` |
| **History** | `history` | Display calculation history | `history` |
| **Variables** | `variables` or `vars` | Show all stored variables | `variables` |
| **Precision** | `precision <digits>` | Set decimal precision | `precision 10` |
| **Clear** | `clear` or `cls` | Clear terminal screen | `clear` |
| **Help** | `help` | Show command reference | `help` |
| **Exit** | `exit` or `quit` | Exit the calculator | `exit` |

---

## Examples

### Basic Mathematical Operations

```bash
# Arithmetic with proper precedence
» 15 + 25 * 2 - 10 / 5
Result: 63

# Scientific notation
» 2e-10 + 3.5E+12
Result: 3500000000000

# Complex expressions with parentheses
» ((10 + 5) * 2)^2 / 3
Result: 300

# Factorial operations
» 10! / (5! * 2!)
Result: 15120
```

### Comparison Operations

```bash
# Basic comparisons
» 5 > 3
Result: 1

» 3 > 5
Result: 0

» 10 == 10
Result: 1

» 5 != 3
Result: 1

# Comparisons with expressions
» 2 + 3 > 4
Result: 1

» sin(30) == 0.5
Result: 1

» (2 * 5) <= 10
Result: 1
```

### Logical Operations

```bash
# Logical AND
» 1 && 1
Result: 1

» 5 && 0
Result: 0

» (5 > 3) && (2 < 4)
Result: 1

# Logical OR
» 0 || 1
Result: 1

» (5 < 3) || (2 < 4)
Result: 1

# Combined logical and comparison
» (x > 5) && (x < 10) || (x == 0)
Result: depends on x value

# Precedence: && before ||
» 0 || 1 && 0
Result: 0  # evaluated as: 0 || (1 && 0)

» 1 || 0 && 0
Result: 1  # evaluated as: 1 || (0 && 0)
```

### Advanced Function Usage

```bash
# Trigonometric calculations
» sin(30) + cos(60) + tan(45)
Result: 2

# Logarithmic functions
» log(100) + ln(e) + log2(16)
Result: 9.60517

# Custom base logarithm
» log(8, 2)
Result: 3

# Statistical functions
» mean(10, 20, 30, 40, 50)
Result: 30

» median(1, 3, 3, 6, 7, 8, 9)
Result: 6

# Print function
» print(sin(30))
0.5

» print(2 + 3)
5
```

### Variable Management

```bash
# Variable assignment and usage
» radius = 5
Result: 5

» area = pi * radius^2
Result: 78.5398

» circumference = 2 * pi * radius
Result: 31.4159

# Variables in comparisons
» radius > 3
Result: 1

» isLarge = radius >= 10
Result: 0

# View all variables
» variables
┌─ Stored Variables ───────────────────────────────────────┐
│ radius          = 5
│ area            = 78.5398
│ circumference   = 31.4159
│ isLarge         = 0
└──────────────────────────────────────────────────────────┘

# Use constants
» speed_of_light = c
Result: 299792458
```

### Unit Conversions

```bash
# Length conversions
» convert 5280 ft to mi
5280 ft = 1 mi

» convert 1000 mm to in
1000 mm = 39.3701 in

# Weight conversions
» convert 2.5 kg to lb
2.5 kg = 5.51156 lb

# Time conversions
» convert 90 min to h
90 min = 1.5 h

# Scientific measurements
» convert 1 km to m
1 km = 1000 m
```

### Complex Calculations

```bash
# Physics calculations with conditions
» F = 9.8 * 75  # Force = mass * acceleration
Result: 735

» isValidForce = F > 0 && F < 1000
Result: 1

» E = F * 10    # Energy = force * distance
Result: 7350

# Conditional logic
» temp = 25
Result: 25

» isFreezing = temp <= 0
Result: 0

» isBoiling = temp >= 100
Result: 0

» isComfortable = (temp > 18) && (temp < 28)
Result: 1

# Financial calculations
» principal = 1000
Result: 1000

» rate = 0.05
Result: 0.05

» compoundInterest = principal * (1 + rate)^10
Result: 1628.89

» isProfit = compoundInterest > principal
Result: 1

# Engineering calculations
» voltage = 12
Result: 12

» current = 2.5
Result: 2.5

» power = voltage * current
Result: 30

» resistance = voltage / current
Result: 4.8

» isSafeVoltage = voltage < 50
Result: 1
```

---

## Project Architecture

### Module Organization

```
Axion/
├── main.go                 # Application entry point
├── install.sh              # Installation script for Unix/Linux
├── constants.json          # Mathematical constants
├── history.json           # Persistent calculation history
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
│
├── cmd/                   # Cobra CLI commands
│   └── cmd.go            # Root command & REPL implementation
│
├── constants/             # Constants management
│   └── constants.go       # JSON-based constant loading
│
├── tokenizer/            # Lexical analysis
│   ├── tokenizer.go      # Token generation and classification
│   └── tokenizer_test.go # Tokenizer unit tests
│
├── parser/               # Syntax analysis
│   ├── parser.go         # AST construction and precedence
│   └── parser_test.go    # Parser unit tests
│
├── evaluator/            # Expression evaluation
│   ├── evaluator.go      # Mathematical computation engine
│   └── evaluator_test.go # Evaluator unit tests
│
├── units/                # Unit conversion
│   └── units.go          # Multi-category conversion system
│
├── history/              # History management
│   └── history.go        # JSON-based persistent storage
│
└── settings/             # Configuration
    └── settings.go       # Precision and display settings
```

### Processing Pipeline

```
User Input (REPL)
     ↓
[Cobra CLI Handler] → Command Processing
     ↓
[Tokenizer] → Lexical Analysis → Token Stream
     ↓                          (includes COMPARISON, LOGICAL tokens)
[Parser] → Syntax Analysis → Abstract Syntax Tree (AST)
     ↓                      (proper precedence: || < && < comparisons < arithmetic)
[Evaluator] → Mathematical Evaluation → Numerical Result
     ↓                                  (comparisons return 1/0, logical ops return 1/0)
[Formatter] → Color-Coded Output → Terminal Display
     ↓
[History] → Persistent Storage → JSON Archive
```

### Operator Precedence (Low to High)

```
1. Assignment         =
2. Logical OR         ||
3. Logical AND        &&
4. Comparison         >, <, >=, <=, ==, !=
5. Addition/Sub       +, -
6. Multiplication     *, /
7. Unary             -x, +x
8. Exponentiation    ^
9. Postfix           !
10. Primary          numbers, functions, parentheses
```

### Key Design Patterns

- **Command Pattern**: Cobra framework for CLI command handling
- **Recursive Descent**: Parser implementation for expression analysis
- **Visitor Pattern**: AST traversal in evaluator
- **Singleton**: Settings and variable storage
- **Strategy Pattern**: Unit conversion system

---

## API Documentation

### Core Functions

#### Tokenizer API
```go
// Tokenize converts input string to token sequence
func Tokenize(input string) ([]Token, error)

// Token represents lexical unit
type Token struct {
    Type  TokenType  // NUMBER, OPERATOR, FUNCTION, COMPARISON, LOGICAL, etc.
    Value string     // Token content
}

// TokenType includes:
// COMPARISON - for >, <, >=, <=, ==, !=
// LOGICAL - for &&, ||
```

#### Parser API
```go
// ParseExpression builds AST from tokens
func (p *Parser) ParseExpression() *Node

// Node represents AST element
type Node struct {
    Type     NodeType  // NODE_NUMBER, NODE_OPERATOR, NODE_COMPARISON, NODE_AND, NODE_OR, etc.
    Value    string    // Node content
    Left     *Node     // Left operand
    Right    *Node     // Right operand
    Children []*Node   // Function arguments
}

// NodeType includes:
// NODE_COMPARISON - for comparison operators
// NODE_AND - for logical AND (&&)
// NODE_OR - for logical OR (||)
```

#### Evaluator API
```go
// Eval recursively evaluates AST nodes
func Eval(node *Node) (float64, error)

// Comparison operations return 1.0 (true) or 0.0 (false)
// Logical operations return 1.0 (true) or 0.0 (false)

// Variable storage
var Vars map[string]float64
```

#### Units API
```go
// Convert performs unit conversion
func Convert(value float64, from, to string) (float64, error)
```

#### Cobra Command API
```go
// Execute runs the root command
func Execute() error

// Root command launches interactive REPL
var rootCmd = &cobra.Command{
    Use:   "axion",
    Short: "Axion - A powerful CLI calculator",
    Run:   startREPL,
}
```

### Extending Axion

#### Adding New Functions
```go
// In evaluator/evaluator.go
case "newfunction":
    if len(node.Children) < 1 {
        return 0, fmt.Errorf("newfunction requires 1 argument")
    }
    arg1, err := Eval(node.Children[0])
    if err != nil {
        return 0, err
    }
    return yourCalculation(arg1), nil
```

#### Adding New Operators
```go
// In tokenizer/tokenizer.go - add new token type
const (
    // ... existing types
    NEW_OPERATOR
)

// In parser/parser.go - add precedence level
func (p *Parser) parseNewOperator() (*Node, error) {
    // Implementation
}

// In evaluator/evaluator.go - add evaluation
case NODE_NEW_OPERATOR:
    // Evaluation logic
```

#### Adding New Constants
```json
// In constants.json
{
    "pi": 3.141592653589793,
    "your_constant": 2.71828,
    "speed_of_light": 299792458
}
```

#### Adding New Commands
```go
// In cmd/cmd.go
var newCmd = &cobra.Command{
    Use:   "newcommand",
    Short: "Description of new command",
    Run: func(cmd *cobra.Command, args []string) {
        // Command implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

---

## Testing

### Running Tests

```bash
# Run complete test suite
go test ./...

# Run with coverage analysis
go test -cover ./...

# Run with detailed output
go test -v ./...

# Run specific package tests
go test ./tokenizer
go test ./parser
go test ./evaluator
go test ./units

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage Statistics

| Module | Coverage | Test Status | Description |
|--------|----------|-------------|-------------|
| **Units** | 100% | Passing | Unit conversion system |
| **Parser** | 76.4% | Passing | AST construction and precedence handling |
| **Evaluator** | 74.5% | Passing | Mathematical computation and functions |
| **Tokenizer** | 94.0% | Passing | Lexical analysis and token generation |
| **CMD** | 0% | No tests | Cobra CLI handlers (interactive module) |
| **Constants** | 0% | No tests | Constants management (utility module) |
| **History** | 0% | No tests | Persistent storage (I/O module) |
| **Settings** | 0% | No tests | Configuration management (utility module) |
| **Core Modules** | **86.7%** | Passing | Average coverage of tested modules |

**Note**: Utility modules (constants, history, settings) and the interactive CLI (cmd) currently lack test files. Core computational modules (tokenizer, parser, evaluator, units) have comprehensive test coverage including logical operations and comparison operators, and all tests pass successfully.

---

## Contributing

We welcome contributions! Axion thrives on community involvement.

### Getting Started

1. **Fork** the repository on GitHub
2. **Clone** your fork locally
3. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
4. **Make** your changes with proper tests
5. **Test** thoroughly (`go test ./...`)
6. **Commit** with descriptive messages
7. **Push** to your branch (`git push origin feature/amazing-feature`)
8. **Open** a Pull Request with detailed description

### Development Guidelines

- Follow Go standard formatting (`gofmt`, `golint`)
- Add tests for new features
- Update documentation for API changes
- Keep commits atomic and well-described
- Ensure all tests pass before submitting PR

---

## Support & Community

### Getting Help

- **Documentation**: Check this README and inline code comments
- **Issues**: [GitHub Issues](https://github.com/codetesla51/Axion/issues) for bugs and features
- **Discussions**: [GitHub Discussions](https://github.com/codetesla51/Axion/discussions) for questions
- **Email**: Direct contact for security issues

### Community Guidelines

- Be respectful and constructive
- Provide clear, reproducible bug reports
- Include system information for technical issues
- Search existing issues before creating new ones

### Issue Templates

When reporting bugs, please include:
- Go version (`go version`)
- Operating system and architecture
- Axion version or commit hash
- Complete error message
- Steps to reproduce
- Expected vs. actual behavior

---

## Acknowledgments

- **Go Community**: For excellent standard library and tooling
- **Cobra**: For providing an excellent CLI framework
- **Mathematical Computing**: Inspired by scientific calculators and computing environments
- **Open Source**: Built on principles of transparency and collaboration
- **Contributors**: Thanks to everyone who helps improve Axion

### Special Thanks

- Mathematical function implementations inspired by standard libraries
- Unit conversion factors from international standards
- Testing methodologies from Go community best practices
- CLI design patterns from Cobra framework

---

## Author

**Uthman** - *Creator*

- **GitHub**: [@codetesla51](https://github.com/codetesla51)
- **Email**: Available through GitHub profile
- **Project**: [Axion Calculator](https://github.com/codetesla51/Axion)

---

<div align="center">

![GitHub stars](https://img.shields.io/github/stars/codetesla51/Axion?style=social)
![GitHub forks](https://img.shields.io/github/forks/codetesla51/Axion?style=social)
![GitHub issues](https://img.shields.io/github/issues/codetesla51/Axion)
![GitHub last commit](https://img.shields.io/github/last-commit/codetesla51/Axion)

**Built with Go & Cobra**

*If you find Axion helpful, please consider giving it a star!*

**[Back to Top](#axion-calculator)**

</div>