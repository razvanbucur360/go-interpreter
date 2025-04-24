package main

import (
	"fmt"
	"os"
)

var hasError = false
var hasRuntimeError = false

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	if !(command == "tokenize" || command == "parse" || command == "evaluate" || command == "run") {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	source := string(fileContents)
	scanner := &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}

	tokens := scanner.ScanTokens()

	if command == "tokenize" {
		printTokens(tokens)

		if hasError {
			os.Exit(65)
		}
		return
	}

	parser := NewParser(tokens)
	statements := parser.parse()

	if command == "parse" {
		astPrinter := NewAstPrinter()
		fmt.Println(astPrinter.Print(statements))

		if hasError {
			os.Exit(65) // Exit with code 65 for compile errors
		}
		return
	}

	if command == "run" || command == "evaluate" {
		if hasError {
			os.Exit(65) // Exit immediately if parse errors exist
		}
		interpreter := NewInterpreter(command == "evaluate")
		resolver := NewResolver(interpreter)
		resolver.resolve(statements)
		if hasError {
			os.Exit(65)
		}
		interpreter.interpret(statements)

		if hasRuntimeError {
			os.Exit(70)
		}
	}
}

// printTokens prints all tokens
func printTokens(tokens []Token) {
	for _, token := range tokens {
		var LiteralStr string
		if token.Literal == nil {
			LiteralStr = "null"
		} else if num, ok := token.Literal.(float64); ok {
			if token.TokenType == NUMBER {
				if num == float64(int(num)) {
					// For whole numbers, add .0
					LiteralStr = fmt.Sprintf("%.1f", num)
				} else {
					// For decimal numbers, use the original Lexeme to preserve exact digits
					LiteralStr = token.Lexeme
				}
			} else {
				LiteralStr = fmt.Sprintf("%v", num)
			}
		} else {
			LiteralStr = fmt.Sprintf("%v", token.Literal)
		}
		fmt.Printf("%s %s %s\n", token.TokenType, token.Lexeme, LiteralStr)
	}
}

func error(Line int, message string) {
	report(Line, "", message)
}

func report(Line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", Line, where, message)
	hasError = true
}
