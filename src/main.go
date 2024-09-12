package main

import (
	"limLang/evaluator"
	"limLang/lexer"
	"limLang/object"
	"limLang/parser"
	"log"
	"os"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Hello %s! This is  ling lang!\n", user.Username)

	if len(os.Args) < 1 {
		log.Fatalf("Usage: %s <file-path>", os.Args[0])
	}

	// Get the file path from the command line arguments
	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	// Convert byte slice to string
	fileContent := string(data)

	// Print the file content
	// fmt.Println(fileContent)
	l := lexer.New(fileContent)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	evaluator.Eval(program, env)
}
