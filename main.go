package main

import (
	"fmt"

	repl "github.com/razzat008/letsgodb/internal/REPl"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
)

func main() {
	// stored the initialized buffer value in line buffer
	lineBuffer := repl.InitLineBuffer()
	for {
		repl.PrintDB()
		lineBuffer.UserInput() // repl.InitLineBuffer.UserInput
		toks := tok.Tokenizer(lineBuffer)
		for _, tok := range toks {
			fmt.Printf("[%s] %s\n", tok.CurrentToken, tok.Type)
		}

		lineBuffer.Reset()
	}
}
