package main

import (
	"compilers/interpreter"
	"compilers/lexer"
	"compilers/parser"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Parse command-line arguments
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatalf("Usage: %s <path to .mygo file>", os.Args[0])
	}
	filePath := flag.Arg(0)

	// Read the .mygo file
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert the code to a string
	input := string(code)

	// Create the lexer
	l := lexer.New(input)

	// Create the parser
	p := parser.New(l)

	// Create and execute the interpreter
	i := interpreter.NewInterpreter(p)
	i.Interpret()
}
