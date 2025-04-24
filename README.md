# Go Lox Interpreter

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

This is a Go implementation of a Lox interpreter, originally developed as part of CodeCrafters' ["Build your own Interpreter" Challenge](https://app.codecrafters.io/courses/interpreter/overview).

## Origins and Credits

This implementation was created by [Your Name] with:
- Guidance from Robert Nystrom's excellent book [Crafting Interpreters](https://craftinginterpreters.com/)
- Initial project structure from [CodeCrafters](https://codecrafters.io)
- Valuable assistance from the developer community

The interpreter follows the tree-walk architecture described in the book, implementing the full Lox language specification including:
- Lexical analysis (scanning)
- Parsing (recursive descent)
- AST generation
- Runtime interpretation
- Classes and inheritance
- Closures and function calls

## Differences from Reference Implementations

While based on the Java and C implementations in Crafting Interpreters, this version:
- Is implemented in idiomatic Go
- Includes additional error handling
- Features some optimizations for Go's runtime
- Contains custom extensions to the testing framework

## Getting Started

### Prerequisites
- Go 1.24 or higher

### Running the Interpreter
```sh
./your_program.sh run [script.lox]