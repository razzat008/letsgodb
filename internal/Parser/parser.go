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

//Parser Struct
type Parser struct {
	Tokens 				[]tok.Token
	position 		  int	
	currentToken 	tok.Token
	peekToken 		tok.Token
}

/* Initializing Parser  */
func (p* Parser) initParser(Tokens []tok.Token){
	p.Tokens 			 = Tokens 								//Stored tokens
	p.position     = 0          						//Tracks the current position in parsing 
	p.currentToken = Tokens[p.position]  		//Initially current token's index = 0 
	p.peekToken 	 = Tokens[p.position + 1] //Next token after current token = 1
}


/* Entry point of the parser */
func ParseProgram (Tokens []tok.Token){
	p := &Parser{}
	p.initParser(Tokens)

	/* gets token type of Tokens[0] */
	switch p.currentToken.Type {
	case tok.TokenSelect:
		p.parseSelect()
	/* Todo: Insert, Update, Delete.....*/
	default: 
		fmt.Println("Unknown Operation detected")
	}
}

func (p* Parser) nextToken (){
	p.position ++
	// If no more tokens is left it will make current token as TokenEOF
	if p.position  < len(p.Tokens){
		p.currentToken = p.Tokens[p.position]
	} else {
		p.currentToken.Type = tok.TokenEOF
	}

	if p.position+1 < len(p.Tokens){
		p.peekToken = p.Tokens[p.position + 1]
	}
}


/*
	-----Currently Parses like this--------------------
 	SELECT COLUMN FROM TABLE {optional:WHERE Condition} ; EOF
*/
func (p* Parser) parseSelect(){
	/* If token after select is not an identifier or an asterisk */
		if (p.peekToken.Type != tok.TokenIdentifier) && 
			 (p.peekToken.Type != tok.TokenAsterisk){
		fmt.Println("Column name or All column expected")
		return 
	}

	/* getting next token after select*/
	p.nextToken()

	/* collecting the column names or asterisk */
	columns := []string{}

	/* loops through all identifiers */
	for {
		/* If current token is not an identifier or asterisk (could be FROM) exit loop */
		if (p.currentToken.Type != tok.TokenIdentifier) && 
			 (p.currentToken.Type != tok.TokenAsterisk){
			break;
		}
		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()
		/* TODO: later we will implement for comma */
	}

	if p.currentToken.Type != tok.TokenFrom {
  	fmt.Println("Expected From clause")
		return
	}

	p.nextToken()
	if p.currentToken.Type != tok.TokenIdentifier {
		fmt.Println("Expected Table name")
		return 
	}
	/* stores the table name to perform operation on the table*/
	table := p.currentToken.CurrentToken

	p.nextToken()
	if  p.currentToken.Type == tok.TokenWhere {
		/* TODO: Complete the Where part */	
		/* What levels of complexity are We adding on it */
	}
	if p.currentToken.Type != tok.TokenSemiColon {
		fmt.Println("Expected Semicolon")
		return
	}
	p.nextToken()	
	if p.currentToken.Type != tok.TokenEOF {
		fmt.Println("Additional token recieved after Semicolon")	
	}
	
	/* 
	Note:
	- For Now we are ignoring the multiline input parsing like in sql 
		if semicolon is not recieved 
	*/
}
