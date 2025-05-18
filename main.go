package main

import (
	repl "github.com/razzat008/letsgodb/internal/REPl"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
)

func main(){
// stored the initialized buffer value in line buffer
	lineBuffer := repl.InitLineBuffer() 	
	for {
		repl.PrintDB()
		lineBuffer.UserInput() // repl.InitLineBuffer.UserInput
		tok.Tokenizer(lineBuffer)
		lineBuffer.Reset()
	}
}
