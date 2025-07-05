/*
idea
-	Fist the tokens recieved is sent to the ParseProgram() {entry point of parser}
-	Then (Token[0].type: TokenType) is passed into the switch statment
- Which would then execute the parse function for each operation
- Ast will be constructed accordingly to the parse function
*/
package parser

import (
	"encoding/json"
	"fmt"

	tok "github.com/razzat008/letsgodb/internal/Tokenizer"
)

// Statement interface for all SQL statements
type Statement interface {
	StatementNode()
}

// AST structs for parsed SQL statements
type SelectStatement struct {
	Columns []string
	Table   string
	Where   *WhereClause // nil if no WHERE clause
}

func (s *SelectStatement) StatementNode() {}

type InsertStatement struct {
	Table   string
	Columns []string
	Values  []string
}

func (i *InsertStatement) StatementNode() {}

type WhereClause struct {
	Column   string
	Operator string
	Value    string
}

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

func ParseProgram(Tokens []tok.Token) Statement {
	if len(Tokens) == 0 {
		fmt.Println("Empty input: no tokens to parse")
		return nil
	}

	p := &Parser{}
	p.initParser(Tokens)

	switch p.currentToken.Type {
	case tok.TokenSelect:
		stmt := p.parseSelect()
		b, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println("Parsed SELECT statement:", string(b))
		return stmt
	// case tok.TokenInsert:
	// 	stmt := p.parseInsert()
	// 	b, _ := json.MarshalIndent(stmt, "", "  ")
	// 	fmt.Println("Parsed INSERT statement:", string(b))
	// 	return stmt
	// Todo: Update, Delete, etc.
	default:
		fmt.Printf("Unknown or unsupported operation: %v\n", p.currentToken.Type)
		return nil
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

func (p *Parser) parseSelect() *SelectStatement {
	// If token after SELECT is not an identifier or an asterisk
	if p.peekToken.Type != tok.TokenIdentifier && p.peekToken.Type != tok.TokenAsterisk {
		fmt.Printf("Syntax error: expected column name or '*' after SELECT, got %v\n", p.peekToken.Type)
		return nil
	}

	// Move to the token after SELECT
	p.nextToken()

	// Collecting column names or '*'
	columns := []string{}

	// Loop through all valid column tokens (identifiers or asterisk)
	for {
		if p.currentToken.Type != tok.TokenIdentifier && p.currentToken.Type != tok.TokenAsterisk {
			break
		}

		// append column name to columns slice
		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()

		// Handle comma-separated columns like SELECT name, age FROM ...
		if p.currentToken.Type == tok.TokenComma {
			p.nextToken()
			continue
		}
	}

	// Expecting FROM keyword after columns
	if p.currentToken.Type != tok.TokenFrom {
		fmt.Printf("Syntax error: expected FROM clause, got %v\n", p.currentToken.Type)
		return nil
	}

	p.nextToken()

	// Expecting a valid table name (identifier) after FROM
	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Printf("Syntax error: expected table name after FROM, got %v\n", p.currentToken.Type)
		return nil
	}

	// Store the table name for later use
	table := p.currentToken.CurrentToken

	p.nextToken()

	// Optional: Handle WHERE clause
	var where *WhereClause
	if p.currentToken.Type == tok.TokenWhere {
		p.nextToken()
		// Expect: IDENTIFIER OPERATOR VALUE
		if p.currentToken.Type != tok.TokenIdentifier {
			fmt.Printf("Syntax error: expected column name in WHERE, got %v\n", p.currentToken.Type)
			return nil
		}
		whereColumn := p.currentToken.CurrentToken
		p.nextToken()
		if p.currentToken.Type != tok.TokenOperator {
			fmt.Printf("Syntax error: expected operator in WHERE, got %v\n", p.currentToken.Type)
			return nil
		}
		whereOperator := p.currentToken.CurrentToken
		p.nextToken()
		if p.currentToken.Type != tok.TokenIdentifier && p.currentToken.Type != tok.TokenValue {
			fmt.Printf("Syntax error: expected value in WHERE, got %v\n", p.currentToken.Type)
			return nil
		}
		whereValue := p.currentToken.CurrentToken
		p.nextToken()
		where = &WhereClause{
			Column:   whereColumn,
			Operator: whereOperator,
			Value:    whereValue,
		}

		// Expecting semicolon to end the query
		if p.currentToken.Type != tok.TokenSemiColon {
			fmt.Printf("Syntax error: expected ';', got %v\n", p.currentToken.Type)
			return nil
		}

		p.nextToken()

		// Ensure no extra tokens after semicolon
		if p.currentToken.Type != tok.TokenEOF {
			fmt.Printf("Syntax warning: unexpected token after semicolon: %v\n", p.currentToken.Type)
		}

		return &SelectStatement{
			Columns: columns,
			Table:   table,
			Where:   where,
		}
	}

	return &SelectStatement{
		Columns: columns,
		Table:   table,
		Where:   where,
	}
}
