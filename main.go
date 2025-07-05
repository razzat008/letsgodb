package main

import (
	par "github.com/razzat008/letsgodb/internal/Parser"
	repl "github.com/razzat008/letsgodb/internal/REPl"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
)

// to print help message
func printHelp() {
	println("letsgodb Help:")
	println("  Type SQL commands to interact with the database.")
	println("  Type 'help' to see this message.")
	println("  Type '\\e' to exit.")
	println("  To learn more about letsgodb, visit https://github.com/razzat008/letsgodb")
}

// main entry point of the program
func main() {
	// stored the initialized buffer value in line buffer
	lineBuffer := repl.InitLineBuffer()
	for {
		repl.PrintDB()
		lineBuffer.UserInput() // repl.InitLineBuffer.UserInput
		input := string(lineBuffer.Buffer)
		if input == "help" {
			printHelp()
			lineBuffer.Reset()
			continue
		} else if input == "\\e" {
			println("Exiting letsgodb...")
			println("Bye!!")
			break
		}
		// since tok.tokenizer returns slice of token struct
		par.ParseProgram(tok.Tokenizer(lineBuffer))
		lineBuffer.Reset() // after one instance of the buffer has been sent the buffer is cleared
	}
}
