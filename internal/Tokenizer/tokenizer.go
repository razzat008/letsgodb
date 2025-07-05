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
	TokenSemiColon  TokenType = "SEMICOLON"
	TokenEOF        TokenType = "EOF"
	TokenComma      TokenType = "COMMA"
	TokenAsterisk   TokenType = "ASTERISK"
)

// break input string into clean token parts
func tokenizeInput(input string) []string {
	var tokens []string
	var current strings.Builder

	for _, ch := range input {
		switch ch {
		case ' ', '\t', '\n':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case ';', ',', '*', '=', '<', '>', '(', ')':
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

// tokenizes the user input []byte into string
// and returns the string
func Tokenizer(lb *repl.LineBuffer) []Token {
	rawTokens := tokenizeInput(string(lb.Buffer))
	var tokens []Token

	for _, currentToken := range rawTokens {
		upperToken := strings.ToUpper(currentToken)
		switch upperToken {
		case "SELECT":
			tokens = append(tokens, Token{Type: TokenSelect, CurrentToken: upperToken})
		case "FROM":
			tokens = append(tokens, Token{Type: TokenFrom, CurrentToken: upperToken})
		case "WHERE":
			tokens = append(tokens, Token{Type: TokenWhere, CurrentToken: upperToken})
		case "=", ">", "<", ">=", "<=", "!=":
			tokens = append(tokens, Token{Type: TokenOperator, CurrentToken: upperToken})
		case ";":
			tokens = append(tokens, Token{Type: TokenSemiColon, CurrentToken: upperToken})
		case ",":
			tokens = append(tokens, Token{Type: TokenComma, CurrentToken: upperToken})
		case "*":
			tokens = append(tokens, Token{Type: TokenAsterisk, CurrentToken: upperToken})
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
