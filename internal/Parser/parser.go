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

// Initially the value stored in the parser
func (p* Parser) initParser(Tokens []tok.Token){
	p.Tokens 			 = Tokens 								//Stored tokens
	p.position     = 0          						//Tracks the current position in parsing 
	p.currentToken = Tokens[p.position]  		//Initially current token's index = 0 
	p.peekToken 	 = Tokens[p.position + 1] //Next token after current token = 1
}


func ParseProgram (Tokens []tok.Token){
	p := &Parser{}
	p.initParser(Tokens)

	// This is the switch case for your next token type
	switch p.currentToken.Type {
	case tok.TokenSelect:
		p.parseSelect()
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

// Parsing starts here
func (p* Parser) parseSelect(){
	// get next token after select
	p.nextToken()
	// collecting the columns name {Identifier}
	columns := []string{}
	for {
		if (p.currentToken.Type != tok.TokenIdentifier) && 
			 (p.currentToken.Type == tok.TokenAsterisk){
			break;
		}
		columns = append(columns, p.currentToken.CurrentToken)
		p.nextToken()
		// later we will implement for comma
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
	table := p.currentToken.CurrentToken
}
