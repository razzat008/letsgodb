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

type Expr interface {
	exprNode()
}

// AST structs for parsed SQL statements
type SelectStatement struct {
	Columns []string
	Table   string
	Where   Expr
}

func (s *SelectStatement) StatementNode() {}

type InsertStatement struct {
	Table   string
	Values  [][]string
	Columns []string
}

func (i *InsertStatement) StatementNode() {}

// AST struct for CREATE TABLE
type CreateTableStatement struct {
	TableName string
	Columns   []string
}

type DropStatement struct {
	Database string
	Table string
	Columns []string
}
func (d *DropStatement) StatementNode() {}

type DeleteStatement struct {
	Table string
	Where Expr
}
func (d *DeleteStatement) StatementNode() {}

func (c *CreateTableStatement) StatementNode() {}


type Condition struct {
	Column   string
	Operator string
	Value    string
}
func (c *Condition) exprNode() {}

type BinaryExpr struct {
	Left     Expr
	Operator string // "AND" or "OR" for now
	Right    Expr
}
func (b *BinaryExpr) exprNode() {}

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
/* Parsing the where clause for Select statement */ 
func (p *Parser) parseExpr() Expr {
	left := p.parsePrimaryExpr()
	if left == nil {
		return nil
	}

	for p.currentToken.Type == tok.TokenAnd || p.currentToken.Type == tok.TokenOr {
		op := p.currentToken.CurrentToken
		p.nextToken()
		right := p.parseExpr()
		if right == nil {
			return nil
		}
		left = &BinaryExpr{
			Left:     left,
			Operator: op,
			Right:    right,
		}
	}
	return left
}

func (p *Parser) parsePrimaryExpr() Expr {
	if p.currentToken.Type == tok.TokenLeftParen {
		p.nextToken()
		expr := p.parseExpr()
		if p.currentToken.Type != tok.TokenRightParen {
			fmt.Println("Syntax error: expected ')' after expression")
			return nil
		}
		p.nextToken()
		return expr
	}
	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Println("Syntax error: expected column name")
		return nil
	}
	column := p.currentToken.CurrentToken
	p.nextToken()

	if p.currentToken.Type != tok.TokenOperator {
		fmt.Println("Syntax error: expected operator")
		return nil
	}
	operator := p.currentToken.CurrentToken
	p.nextToken()

	if p.currentToken.Type != tok.TokenIdentifier && p.currentToken.Type != tok.TokenValue {
		fmt.Println("Syntax error: expected value")
		return nil
	}
	value := p.currentToken.CurrentToken
	p.nextToken()

	return &Condition{
		Column:   column,
		Operator: operator,
		Value:    value,
	}
}

// Parse CREATE TABLE statement
func (p *Parser) parseCreateTable() *CreateTableStatement {
	// Expect: CREATE TABLE table_name (col1, col2, ...)
	p.nextToken() // move to TABLE
	if p.currentToken.Type != tok.TokenTable {
		fmt.Printf("Syntax error: expected TABLE after CREATE, got %v (did you forget the TABLE keyword?)\n", p.currentToken.Type)
		return nil
	}
	p.nextToken() // move to table name
	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Printf("Syntax error: expected table name, got %v\n", p.currentToken.Type)
		return nil
	}
	tableName := p.currentToken.CurrentToken
	p.nextToken()
	if p.currentToken.Type != tok.TokenLeftParen {
		fmt.Printf("Syntax error: expected '(' after table name, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()
	columns := []string{}
	for p.currentToken.Type == tok.TokenIdentifier {
		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()
		if p.currentToken.Type == tok.TokenComma {
			p.nextToken()
		}
	}
	if p.currentToken.Type != tok.TokenRightParen {
		fmt.Printf("Syntax error: expected ')' after column list, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()
	if p.currentToken.Type != tok.TokenSemiColon {
		fmt.Printf("Syntax error: expected ';' at end of statement, got %v\n", p.currentToken.Type)
		return nil
	}
	return &CreateTableStatement{TableName: tableName, Columns: columns}
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
	case tok.TokenInsert:
		stmt := p.parseInsert()
		b, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println("Parsed INSERT statement:", string(b))
		return stmt
	case tok.TokenCreate:
		stmt := p.parseCreateTable()
		b, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println("Parsed CREATE TABLE statement:", string(b))
		return stmt
	case tok.TokenDrop:
		stmt := p.parseDrop()
		b, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println("Parsed Drop TABLE statement:", string(b))
		return stmt
	case tok.TokenDelete: 
		stmt := p.parseDelete()
		b, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println("Parsed Delete TABLE statement:", string(b))
		return stmt
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
	if p.position+1 < len(p.Tokens) { // to advance to the next token
		p.peekToken = p.Tokens[p.position+1] // to assign the next token
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

	var where Expr
	if p.currentToken.Type == tok.TokenWhere {
		p.nextToken()
		where = p.parseExpr()
		if where == nil {
			return nil
		}
	}

	return &SelectStatement{
		Columns: columns,
		Table:   table,
		Where:   where,
	}
}

func (p *Parser) parseInsert() *InsertStatement {
	if p.peekToken.Type != tok.TokenInto {
		fmt.Printf("Syntax error: expected INTO , got %v\n", p.peekToken.Type)
		return nil
	}
	p.nextToken()
	p.nextToken()

	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Printf("Syntax error: expected table name after INTO, got %v\n", p.currentToken.Type)
		return nil
	}
	table := p.currentToken.CurrentToken
	p.nextToken()

	if p.currentToken.Type != tok.TokenLeftParen {
		fmt.Printf("Syntax error: expected '(' after table name, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	columns := p.parseColumns()
	if columns == nil {
		return nil
	}

	if p.currentToken.Type != tok.TokenRightParen {
		fmt.Printf("Syntax error: expected ')' after column list, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	if p.currentToken.Type != tok.TokenValues {
		fmt.Printf("Syntax error: expected 'VALUES' after column list, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	if p.currentToken.Type != tok.TokenLeftParen {
		fmt.Printf("Syntax error: expected '(' after VALUES, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	values := p.parseValues()
	if values == nil {
		return nil
	}

	if p.currentToken.Type != tok.TokenRightParen {
		fmt.Printf("Syntax error: expected ')' after value list, got %v\n", p.currentToken.Type)
		return nil
	}

	return &InsertStatement{
		Table:   table,
		Values:  values,
		Columns: columns,
	}
}

func (p *Parser) parseColumns() []string {
	columns := []string{}
	for p.currentToken.Type == tok.TokenIdentifier {
		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()
		if p.currentToken.Type == tok.TokenComma {
			p.nextToken()
		}
	}

	if p.currentToken.Type != tok.TokenRightParen {
		fmt.Printf("Syntax error: expected ')' after column list, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	return columns
}

func (p *Parser) parseValues() [][]string {
	if p.currentToken.Type != tok.TokenLeftParen {
		fmt.Printf("Syntax error: expected '(' after VALUES, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	values := [][]string{}
	for p.currentToken.Type == tok.TokenStringLiteral {
		values = append(values, []string{p.currentToken.CurrentToken})
		p.nextToken()
		if p.currentToken.Type == tok.TokenComma {
			p.nextToken()
		}
	}

	if p.currentToken.Type != tok.TokenRightParen {
		fmt.Printf("Syntax error: expected ')' after value list, got %v\n", p.currentToken.Type)
		return nil
	}
	p.nextToken()

	return values
}


func (p *Parser) parseDrop() *DropStatement{
	/* 
		Examples: 
		- Drop table table_name 
		- Drop table table_name ( columns ) --> drop columns of table_name
		- Drop table table_name1 table_name2 --> drop multiple tables
	*/
	var database string 
	var columns []string
	var table string
	p.nextToken() // database or table token
	switch p.currentToken.Type {
	case tok.TokenTable:
	p.nextToken()
	if p.currentToken.Type != tok.TokenIdentifier {
			fmt.Printf("Syntax error: expected IDENTIFIER, got %v\n", p.currentToken.Type)
			return nil
		}
		table = p.currentToken.CurrentToken

	if p.peekToken.Type == tok.TokenSemiColon{
			break
		} 
			if p.peekToken.Type != tok.TokenLeftParen{
			fmt.Printf("Syntax error: expected ( , got %v\n", p.peekToken.Type)
			return nil
		} 
		p.nextToken() // at parenthesis 
		p.nextToken() // at identifier 

		//loops through identifier
		for p.currentToken.Type == tok.TokenIdentifier {
			columns = append(columns, p.currentToken.CurrentToken)
			p.nextToken()
			if p.currentToken.Type == tok.TokenComma {
				p.nextToken()
			}
		}
		if p.currentToken.Type != tok.TokenRightParen {
			fmt.Printf("Syntax error: expected ) , got %v\n", p.peekToken.Type)
			return nil
		}
		 

	case tok.TokenDatabase:
		p.nextToken()
		if p.currentToken.Type!= tok.TokenIdentifier{
			fmt.Printf("Syntax error: expected IDENTIFIER, got %v\n", p.currentToken.Type)
			return nil
		}
		database = p.currentToken.CurrentToken
		if p.peekToken.Type != tok.TokenSemiColon{
			fmt.Printf("Syntax error: expected Semicolon, got %v\n", p.peekToken.Type)
			return nil
		}
	default:
		fmt.Printf("Syntax error: expected Table or Database, got %v\n",p.currentToken.Type)
	}
	return &DropStatement{
		Database: database,
		Table: table,
		Columns : columns,
	}
}

func (p *Parser) parseDelete() *DeleteStatement {
	//  DELETE FROM table_name [WHERE condition]
	if p.peekToken.Type != tok.TokenFrom {
		fmt.Printf("Syntax error: expected FROM after DELETE, got %v\n", p.peekToken.Type)
		return nil
	}
	p.nextToken()
	p.nextToken()

	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Printf("Syntax error: expected IDENTIFIER , got %v\n", p.currentToken.Type)
		return nil
	}

	table := p.currentToken.CurrentToken
	p.nextToken()

	var where Expr
	if p.currentToken.Type == tok.TokenWhere {
		p.nextToken()
		where = p.parseExpr()
		if where == nil {
			return nil
		}
	}

	return &DeleteStatement{
		Table: table,
		Where: where,
	}
}
