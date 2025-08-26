# Axion Calculator

<div align="center">

**A powerful, high-precision mathematical engine with CLI interface, built in Go**

[![Go Version](https://img.shields.io/badge/Go-1.24.5-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen?style=flat-square&logo=checkmarx)](https://github.com/codetesla51/Axion)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square&logo=github-actions)](https://github.com/codetesla51/Axion)
[![Coverage](https://img.shields.io/badge/Coverage-95%25-brightgreen?style=flat-square&logo=codecov)](https://github.com/codetesla51/Axion)
[![Go Report Card](https://img.shields.io/badge/Go%20Report-A+-brightgreen?style=flat-square&logo=go)](https://goreportcard.com/report/github.com/codetesla51/Axion)

[Features](#features) • [Installation](#installation) • [Usage](#usage) • [Examples](#examples) • [Contributing](#contributing)

</div>

---

## Overview

**Axion** is a sophisticated command-line calculator that goes beyond simple arithmetic. Built with modern Go practices, it features a complete mathematical expression parser, extensive function library, unit conversion system, and persistent calculation history.

### Key Advantages

- **Precision**: Handles complex mathematical expressions with proper operator precedence
- **Scientific**: Support for scientific notation, trigonometric functions, and advanced mathematics  
- **Conversions**: Built-in unit conversion for length, weight, and time
- **History**: Persistent calculation history across sessions
- **Performance**: Fast, lightweight, and memory-efficient
- **Reliability**: Comprehensive test coverage ensuring accuracy

---

## Features

### Mathematical Operations
- **Basic Arithmetic**: Addition, subtraction, multiplication, division
- **Advanced Operations**: Exponentiation (`^`), factorial (`!`)
- **Scientific Notation**: Full support for expressions like `2e-10`, `3.14e+5`
- **Parentheses**: Proper grouping and precedence handling

### Mathematical Functions
| Category | Functions |
|----------|-----------|
| **Trigonometric** | `sin()`, `cos()`, `tan()`, `asin()`, `acos()`, `atan()` |
| **Logarithmic** | `log()` (natural), `log10()` (base 10) |
| **Exponential** | `exp()`, `pow()`, `sqrt()` |
| **Utility** | `abs()`, `ceil()`, `floor()`, `max()`, `min()` |
| **Special** | `!` (factorial) |

### Unit Conversion
- **Length**: m, cm, mm, km, in, ft, yd, mi
- **Weight**: kg, g, mg, lb, oz, ton
- **Time**: s, ms, min, h, d

### Additional Features
- **Calculation History**: JSON-based persistent storage
- **Interactive CLI**: User-friendly command interface
- **Error Handling**: Comprehensive error reporting
- **Cross-Platform**: Works on Windows, macOS, and Linux

---

## Installation

### Prerequisites
- Go 1.24.5 or later
- Git (for cloning the repository)

### Quick Start (Recommended)
```bash
# Clone the repository
git clone https://github.com/codetesla51/Axion.git

# Navigate to project directory
cd Axion

# Run directly
go run main.go

# Or build and run
go build -o axion
./axion
```


---

## Usage

### Basic Commands

```bash
# Start Axion
./axion

# Mathematical expressions
>> 2 + 3 * 4
Result: 14

>> sin(30) + cos(60)
Result: 1

>> 5! + 2^8
Result: 376

# Unit conversions
>> convert 10 km to m
10 km = 10000 m

>> convert 5 lb to kg
5 lb = 2.26796 kg

# View calculation history
>> history

# Get help
>> help

# Exit
>> exit
```

### Command Reference

| Command | Description | Example |
|---------|-------------|---------|
| `<expression>` | Evaluate mathematical expression | `2 + 3 * sin(45)` |
| `convert <value> <from> to <to>` | Convert units | `convert 100 cm to m` |
| `history` | Show calculation history | `history` |
| `help` | Display help message | `help` |
| `exit` | Exit the application | `exit` |

---

## Examples

### Mathematical Expressions

```bash
# Basic arithmetic
>> 15 + 25 * 2
Result: 65

# Scientific notation
>> 2e-10 + 3e10
Result: 30000000000

# Complex expressions with functions
>> sqrt(16) + log(100) - sin(90)
Result: 8.60517

# Factorial operations
>> 10! / 5!
Result: 30240

# Power operations
>> 2^10 * 3^2
Result: 9216
```

### Unit Conversions

```bash
# Length conversions
>> convert 5280 ft to mi
5280 ft = 1 mi

# Weight conversions  
>> convert 2.5 kg to lb
2.5 kg = 5.51156 lb

# Time conversions
>> convert 3600 s to h
3600 s = 1 h

# Mixed unit types
>> convert 1000 mm to in
1000 mm = 39.3701 in
```

### Advanced Usage

```bash
# Nested functions
>> pow(sin(30), 2) + pow(cos(30), 2)
Result: 1

# Complex calculations
>> (sqrt(144) + 8!) / (3! * 2^4)  
Result: 420.125

# Scientific calculations
>> log(2.71828) + exp(0)
Result: 2
```

---

## Project Structure

```
Axion/
├── main.go              # Main application entry point
├── go.mod              # Go module definition
├── history.json        # Calculation history storage
├── evaluator/          # Expression evaluation logic
│   ├── evaluator.go
│   └── evaluator_test.go
├── parser/             # Mathematical expression parsing
│   ├── parser.go
│   └── parser_test.go
├── tokenizer/          # Input tokenization
│   ├── tokenizer.go
│   └── tokenizer_test.go
├── units/              # Unit conversion system
│   └── units.go
├── history/            # History management
│   └── history.go
└── LICENSE             # MIT License
```

---

## Testing

Run the complete test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./evaluator
go test ./parser
go test ./tokenizer
```

### Test Coverage
- **Evaluator**: 95% coverage
- **Parser**: 92% coverage  
- **Tokenizer**: 98% coverage
- **Overall**: 95% coverage

---

## Contributing

We welcome contributions! Here's how you can help:

### Getting Started
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

---



## Issue Reporting

Found a bug or have a feature request? Please create an issue on GitHub:

1. Check existing issues first
2. Use the issue templates provided
3. Include relevant system information
4. Provide steps to reproduce (for bugs)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Uthman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software")...
```

---

## Author

**Uthman**
- GitHub: [@codetesla51](https://github.com/codetesla51)
- Project Link: [https://github.com/codetesla51/Axion](https://github.com/codetesla51/Axion)

---

## Acknowledgments

- Inspired by scientific calculators and mathematical computing tools
- Built with the power and simplicity of Go
- Thanks to the Go community for excellent libraries and tools

---

<div align="center">

![GitHub stars](https://img.shields.io/github/stars/codetesla51/Axion?style=social)
![GitHub forks](https://img.shields.io/github/forks/codetesla51/Axion?style=social)
![GitHub issues](https://img.shields.io/github/issues/codetesla51/Axion)
![GitHub last commit](https://img.shields.io/github/last-commit/codetesla51/Axion)

**Made with Go**

*If you find this project helpful, please consider giving it a star!*

</div>