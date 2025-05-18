/*
	- In the tokenization process we are going to tokenize the word on each white space
	 and store tokens in a slice for further processing

	- Then we are going to define a behavour of what each token does which is also 
		managed by the tokenizer
*/

package Tokenizer

import (
	repl "DataGoesBase/internal/REPL"
	"strings"
)

// Now we have to provide meaning to our obtained tokens


// tokenizes the user input []byte into string
//and returns the string 
func Tokenizer(lb* repl.LineBuffer) []string {
	tokens := strings.Fields(string(lb.Buffer))
	return tokens 
}
