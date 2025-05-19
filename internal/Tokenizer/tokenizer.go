/*
	- In the tokenization process we are going to tokenize the currentToken on each white space
	 and store tokens in a slice for further processing

	- Then we are going to define a behavour of what each token does which is also
		managed by the tokenizer
*/

package Tokenizer

import (
	repl "github.com/razzat008/letsgodb/internal/REPl"
	"strings"
)

type TokenType string // a custom type for storing our types

type Token struct {
	CurrentToken string    // the token we're currenlty looking at
	Type         TokenType // the type of the CurrentToken
}

// specifying what possible token	can be
const (
	TokenSelect     TokenType = "SELECT"
	TokenFrom       TokenType = "FROM"
	TokenWhere      TokenType = "WHERE"
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenOperator   TokenType = "OPERATOR"
	TokenValue      TokenType = "VALUE"
	TokenUnknown    TokenType = "UNKNOWN"
)

// tokenizes the user input []byte into string
// and returns the string
func Tokenizer(lb *repl.LineBuffer) []Token {
	rawTokens := strings.Fields(string(lb.Buffer))
	var tokens []Token

	for _, currentToken := range rawTokens {
		switch strings.ToUpper(currentToken) {
		case "SELECT": 
			tokens = append(tokens, Token{Type: TokenSelect, CurrentToken: currentToken})
		case "FROM":
			tokens = append(tokens, Token{Type: TokenFrom, CurrentToken: currentToken})
		case "WHERE":
			tokens = append(tokens, Token{Type: TokenWhere, CurrentToken: currentToken})
		case "=", ">", "<", ">=", "<=", "!=":
			tokens = append(tokens, Token{Type: TokenOperator, CurrentToken: currentToken})
		default:
			if strings.HasPrefix(currentToken, "'") && strings.HasSuffix(currentToken, "'") {
				// checking for values like 'School' ; i.e. quoted values
				tokens = append(tokens, Token{Type: TokenValue, CurrentToken: currentToken})
			} else {
				// if nothing match treat it as an identifier
				tokens = append(tokens, Token{Type: TokenIdentifier, CurrentToken: currentToken})
			}
		}
	}
	return tokens
}
/*
A Token may have two characteristics 
- it's name(value) and it's type
a struct Token is created
and some possible token types is then declared
=====
the contents of the buffer is now looped through looking for matching tokens
to assign their types individually 
*/
