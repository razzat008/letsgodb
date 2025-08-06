/*
	- In the tokenization process we are going to tokenize the currentToken on each white space
	 and store tokens in a slice for further processing

	- Then we are going to define a behavour of what each token does which is also
		managed by the tokenizer
*/

package Tokenizer

import (
	"strings"

	repl "github.com/razzat008/letsgodb/internal/REPl"
)

type TokenType string // a custom type for storing our types

type Token struct {
	CurrentToken string    // the token we're currenlty looking at
	Type         TokenType // the type of the CurrentToken
}

// specifying what possible token	can be
const (
	TokenSelect        TokenType = "SELECT"
	TokenInsert        TokenType = "INSERT"
	TokenDelete        TokenType = "DELETE"
	TokenFrom          TokenType = "FROM"
	TokenUse           TokenType = "USE"
	TokenWhere         TokenType = "WHERE"
	TokenCreate        TokenType = "CREATE"
	TokenTable         TokenType = "TABLE"
	TokenIdentifier    TokenType = "IDENTIFIER"
	TokenOperator      TokenType = "OPERATOR"
	TokenValue         TokenType = "VALUE"
	TokenInto          TokenType = "INTO"
	TokenShow          TokenType = "SHOW"
	TokenValues        TokenType = "VALUES"
	TokenUnknown       TokenType = "UNKNOWN"
	TokenSemiColon     TokenType = "SEMICOLON"
	TokenEOF           TokenType = "EOF"
	TokenComma         TokenType = "COMMA"
	TokenAsterisk      TokenType = "ASTERISK"
	TokenLeftParen     TokenType = "LEFT_PAREN"
	TokenRightParen    TokenType = "RIGHT_PAREN"
	TokenStringLiteral TokenType = "STRING_LITERAL"
	TokenAnd           TokenType = "AND"
	TokenOr            TokenType = "OR"
	TokenDatabase      TokenType = "DATABASE"
	TokenDrop          TokenType = "DROP"
	TokenList          TokenType = "LIST"
	TokenPrimaryKey    TokenType = "PRIMARY_KEY"
)

// break input string into clean token parts
func tokenizeInput(input string) []string {
	var tokens []string
	var current strings.Builder
	i := 0
	for i < len(input) {
		ch := input[i]

		// Handle whitespace
		if ch == ' ' || ch == '\t' || ch == '\n' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			i++
			continue
		}

		// Handle multi-character operators
		if i+1 < len(input) {
			twoChar := input[i : i+2]
			if twoChar == ">=" || twoChar == "<=" || twoChar == "!=" {
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
				tokens = append(tokens, twoChar)
				i += 2
				continue
			}
		}

		// Handle single-character symbols
		if strings.ContainsRune(";,*=<>()", rune(ch)) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(ch))
			i++
			continue
		}

		// Default: add to current token
		current.WriteByte(ch)
		i++
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
		case "INSERT":
			tokens = append(tokens, Token{Type: TokenInsert, CurrentToken: upperToken})
		case "AND":
			tokens = append(tokens, Token{Type: TokenAnd, CurrentToken: upperToken})
		case "LIST":
			tokens = append(tokens, Token{Type: TokenList, CurrentToken: upperToken})
		case "OR":
			tokens = append(tokens, Token{Type: TokenOr, CurrentToken: upperToken})
		case "INTO":
			tokens = append(tokens, Token{Type: TokenInto, CurrentToken: upperToken})
		case "USE":
			tokens = append(tokens, Token{Type: TokenUse, CurrentToken: upperToken})
		case "DELETE":
			tokens = append(tokens, Token{Type: TokenDelete, CurrentToken: upperToken})
		case "VALUES":
			tokens = append(tokens, Token{Type: TokenValues, CurrentToken: upperToken})
		case "CREATE":
			tokens = append(tokens, Token{Type: TokenCreate, CurrentToken: upperToken})
		case "TABLE":
			tokens = append(tokens, Token{Type: TokenTable, CurrentToken: upperToken})
		case "DATABASE":
			tokens = append(tokens, Token{Type: TokenDatabase, CurrentToken: upperToken})
		case "SHOW":
			tokens = append(tokens, Token{Type: TokenShow, CurrentToken: upperToken})
		case "DROP":
			tokens = append(tokens, Token{Type: TokenDrop, CurrentToken: upperToken})
		case "PRIMARY_KEY":
			tokens = append(tokens, Token{Type: TokenPrimaryKey, CurrentToken: upperToken})
		case "=", ">", "<", ">=", "<=", "!=":
			tokens = append(tokens, Token{Type: TokenOperator, CurrentToken: upperToken})
		case ";":
			tokens = append(tokens, Token{Type: TokenSemiColon, CurrentToken: upperToken})
		case ",":
			tokens = append(tokens, Token{Type: TokenComma, CurrentToken: upperToken})
		case "*":
			tokens = append(tokens, Token{Type: TokenAsterisk, CurrentToken: upperToken})
		case "(":
			tokens = append(tokens, Token{Type: TokenLeftParen, CurrentToken: upperToken})
		case ")":
			tokens = append(tokens, Token{Type: TokenRightParen, CurrentToken: upperToken})
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
