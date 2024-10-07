package lex

import "fmt"

type TokenType int
type Pos int

type Token struct {
	Type          TokenType
	StartPosition Pos
	EndPosition   Pos
	LineNumber    int
	Value         string
}

const (
	// SIGNALS
	TokenEOF   TokenType = iota // Raised when the end of the input is reached
	TokenError TokenType = iota // Raised when an error occurs

	// TOKENS
	TokenNumber     TokenType = iota // 1
	TokenString     TokenType = iota // "Helllo World"
	TokenPlus       TokenType = iota // +
	TokenMinus      TokenType = iota // -
	TokenStar       TokenType = iota // *
	TokenSlash      TokenType = iota // /
	TokenColon      TokenType = iota // :
	TokenKey        TokenType = iota // thisbit:<expression>
	TokenLeftParen  TokenType = iota // (
	TokenRightParen TokenType = iota // )
	TokenTag        TokenType = iota // +HOME, -HOME
	TokenWord       TokenType = iota // hello
	TokenCommand    TokenType = iota // add, modify, etc.
)

// Create a new token of a given type
func NewToken(tokenType TokenType, startPosition Pos, endPosition Pos, lineNumber int, value string) *Token {
	return &Token{
		Type:          tokenType,
		StartPosition: startPosition,
		EndPosition:   endPosition,
		LineNumber:    lineNumber,
		Value:         value,
	}
}

func (t *Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return fmt.Sprint("Error: ", t.Value)
	case TokenString:
		return fmt.Sprint("String: ", t.Value)
	case TokenKey:
		return fmt.Sprint("Key: ", t.Value)
	case TokenTag:
		return fmt.Sprint("Tag: ", t.Value)
	case TokenNumber:
		return fmt.Sprint("Number: ", t.Value)
	case TokenPlus:
		return fmt.Sprint("Plus: ", t.Value)
	case TokenMinus:
		return fmt.Sprint("Minus: ", t.Value)
	case TokenStar:
		return fmt.Sprint("Star: ", t.Value)
	case TokenSlash:
		return fmt.Sprint("Slash: ", t.Value)
	case TokenCommand:
		return fmt.Sprint("Command: ", t.Value)
	case TokenWord:
		return fmt.Sprint("Word: ", t.Value)
  case TokenColon:
    return fmt.Sprint("Colon: ", t.Value)
  case TokenLeftParen:
    return fmt.Sprint("LeftParen: ", t.Value)
  case TokenRightParen:
    return fmt.Sprint("RightParen: ", t.Value)
	default:
		fmt.Println(t.Type)
		panic("Unknown token type")
	}
}