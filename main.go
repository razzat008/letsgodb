package main

import (
	repl "github.com/razzat008/letsgodb/internal/REPl"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
	par "github.com/razzat008/letsgodb/internal/Parser"
)

func main() {
	// stored the initialized buffer value in line buffer
	lineBuffer := repl.InitLineBuffer()
	for {
		repl.PrintDB()
		lineBuffer.UserInput() // repl.InitLineBuffer.UserInput
		// since tok.tokenizer returns slice of token struct 
		par.ParseProgram(tok.Tokenizer(lineBuffer))	
		lineBuffer.Reset()
	}
}
