/*
idea
-	Fist the tokens recieved is sent to the ParseProgram() {entry point of parser}
-	Then (Token[0].type: TokenType) is passed into the switch statment
- Which would then execute the parse function for each operation
- Ast will be constructed accordingly to the parse function
*/
package parser

import (
	tok	"github.com/razzat008/letsgodb/internal/Tokenizer"
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

	// Until and unless token Type EOF is found keep running the loop
	for p.currentToken.Type != tok.TokenEOF {
		p.ParseStatement() //parses the the token based on it's type
		p.nextToken() 		 //gets the next token if exist
	}
}

func (p* Parser) ParseStatement(){
	switch p.Tokens[p.position].Type {

	case tok.TokenSelect:
		p.parseSwitchStatement()
		break
	case tok.TokenFrom:
		break
	case tok.TokenWhere:
		break
	case tok.TokenOperator:
		break
	case tok.TokenComma:
		break
	case tok.TokenIdentifier:
		break
	case tok.TokenValue:
		break
	case tok.TokenUnknown:
		break
	}
}

// Parsing switch statement starts here
func (p* Parser) parseSwitchStatement(){}


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
