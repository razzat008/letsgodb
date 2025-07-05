/*
idea
-	Fist the tokens recieved is sent to the ParseProgram() {entry point of parser}
-	Then (Token[0].type: TokenType) is passed into the switch statment
- Which would then execute the parse function for each operation
- Ast will be constructed accordingly to the parse function
*/
package parser

import (
	"fmt"
	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
)

// Parser Struct
type Parser struct {
	Tokens       []tok.Token // array of tokens from the tokenizer
	position     int         // current position in the token stream
	currentToken tok.Token   // currently processed token
	peekToken    tok.Token   // lookahead token (next token)
}

/* Initializing Parser  */
func (p *Parser) initParser(Tokens []tok.Token) {
	p.Tokens = Tokens // Store the full slice of tokens
	p.position = 0    // Start parsing from the beginning

	// Assign the current and peek token safely
	if len(Tokens) > 0 {
		p.currentToken = Tokens[0] // First token
	}
	if len(Tokens) > 1 {
		p.peekToken = Tokens[1] // Second token (lookahead)
	} else {
		p.peekToken = tok.Token{Type: tok.TokenEOF} // If no second token, mark as EOF
	}
}

/* Entry point of the parser */
func ParseProgram(Tokens []tok.Token) {
	if len(Tokens) == 0 {
		fmt.Println("Empty input: no tokens to parse")
		return
	}

	p := &Parser{}
	p.initParser(Tokens)

	/* gets token type of Tokens[0] */
	switch p.currentToken.Type {
	case tok.TokenSelect:
		p.parseSelect()
	/* Todo: Insert, Update, Delete.....*/
	default:
		fmt.Printf("Unknown or unsupported operation: %v\n", p.currentToken.Type)
	}
}

/* Advance to the next token */
func (p *Parser) nextToken() {
	p.position++

	// Safely assign currentToken
	if p.position < len(p.Tokens) {
		p.currentToken = p.Tokens[p.position]
	} else {
		p.currentToken = tok.Token{Type: tok.TokenEOF} // No more tokens left
	}

	// Safely assign peekToken
	if p.position+1 < len(p.Tokens) {
		p.peekToken = p.Tokens[p.position+1]
	} else {
		p.peekToken = tok.Token{Type: tok.TokenEOF}
	}
}

/*
		-----Currently Parses like this--------------------
	 	SELECT COLUMN FROM TABLE {optional:WHERE Condition} ; EOF
*/
func (p *Parser) parseSelect() {
	/* If token after SELECT is not an identifier or an asterisk */
	if p.peekToken.Type != tok.TokenIdentifier && p.peekToken.Type != tok.TokenAsterisk {
		fmt.Printf("Syntax error: expected column name or '*' after SELECT, got %v\n", p.peekToken.Type)
		return
	}

	/* Move to the token after SELECT */
	p.nextToken()

	/* Collecting column names or '*' */
	columns := []string{}

	/* Loop through all valid column tokens (identifiers or asterisk) */
	for {
		/* If current token is not valid as a column name, break the loop */
		if p.currentToken.Type != tok.TokenIdentifier && p.currentToken.Type != tok.TokenAsterisk {
			break
		}

		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()

		// TODO: Handle comma-separated columns like SELECT name, age FROM ...
	}

	/* Expecting FROM keyword after columns */
	if p.currentToken.Type != tok.TokenFrom {
		fmt.Printf("Syntax error: expected FROM clause, got %v\n", p.currentToken.Type)
		return
	}

	p.nextToken()

	/* Expecting a valid table name (identifier) after FROM */
	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Printf("Syntax error: expected table name after FROM, got %v\n", p.currentToken.Type)
		return
	}

	/* Store the table name for later use */
	table := p.currentToken.CurrentToken

	p.nextToken()

	/* Optional: Handle WHERE clause */
	if p.currentToken.Type == tok.TokenWhere {
		/* TODO: Parse WHERE conditions here */
		fmt.Println("WHERE clause found — not yet implemented")
		p.nextToken()
	}

	/* Expecting semicolon to end the query */
	if p.currentToken.Type != tok.TokenSemiColon {
		fmt.Printf("Syntax error: expected ';', got %v\n", p.currentToken.Type)
		return
	}

	p.nextToken()

	/* Ensure no extra tokens after semicolon */
	if p.currentToken.Type != tok.TokenEOF {
		fmt.Printf("Syntax warning: unexpected token after semicolon: %v\n", p.currentToken.Type)
	}

	/*
		Note:
		- For now, we are ignoring multiline SQL input
		- Semicolon is required to indicate end of statement
	*/

	/* Debug output for now */
	fmt.Println("Parsed SELECT query:")
	fmt.Println("  Columns:", columns)
	fmt.Println("  Table:", table)
}
