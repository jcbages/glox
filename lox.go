package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var hadError = false
var hadRuntimeError = false
var interpreter = NewInterpreter()

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	run(string(content))

	if hadError {
		os.Exit(65)
	}

	if hadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewScanner(os.Stdin)

	for reader.Scan() {
		fmt.Print("> ")
		run(reader.Text())
		hadError = false
		hadRuntimeError = false
	}
}

func run(content string) {
	scanner := NewScanner(content)
	tokens := scanner.ScanTokens()

	fmt.Println("--- BEGIN TOKENS --- ")
	for _, token := range tokens {
		fmt.Println(token)
	}
	fmt.Println("---- END TOKENS ---- ")

	parser := NewParser(tokens)
	program, _ := parser.Parse()
	if hadError {
		return
	}

	fmt.Println("--- BEGIN AST ---")
	fmt.Println((&AstPrinter{}).print(program))
	fmt.Println("---- END AST ----")

	interpreter.Interpret(program)
}

func LoxError(line int, message string) {
	Report(line, "", message)
}

func LoxTokenError(token Token, message string) {
	if token.TokenType == EOF {
		Report(token.Line, "at end", message)
	} else {
		Report(token.Line, fmt.Sprintf("at '%v'", token.Lexeme), message)
	}
}

func LoxRuntimeError(token Token, message string) {
	log.Println(message)
	log.Printf("[line %v]", token.Line)
	hadRuntimeError = true
}

func Report(line int, where string, message string) {
	log.Printf("[line %v] Error %v: %v\n", line, where, message)
	hadError = true
}
