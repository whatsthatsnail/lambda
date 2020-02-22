package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"lambda/lexer"
	"lambda/parser"
	"lambda/interpreter"
	"os"
)

// Gets arguments when using 'go run *.go -- ...'
func main() {

	var args []string
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	if len(args) >= 1 && args[0] == "repl" {
		repl()
	} else if len(args) >= 3 && args[0] == "file" && args[2] == "-v" {
		file(args[1], false)
	} else if len(args) >= 2 && args[0] == "file" {
		file(args[1], true)
	} else {
		repl()
	}
}

// Helper function to check for errors when reading files
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func repl() {
	fmt.Printf("Entering REPL:\n>>> ")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			os.Exit(0)
		}

		lex := lexer.NewLexer(line, true)
		tokens, _ := lex.ScanTokens()

		lexer.PrintTokens(tokens)
	}
}

// Reads file into lexer, tokenizes, and prints tokens
func file(path string, quiet bool) {
	dat, err := ioutil.ReadFile(path)
	check(err)

	lex := lexer.NewLexer(string(dat), false)
	tokens, _ := lex.ScanTokens()

	if !quiet {
		fmt.Println(path + ":" + "\n")
		fmt.Println(string(dat) + "\n")
		lexer.PrintTokens(tokens)
	}

	parser := parser.NewParser(tokens)
	tree, _ := parser.Parse()

	inter := interpreter.NewInterpreter(tree)
	fmt.Println(inter.Evaluate())
}
