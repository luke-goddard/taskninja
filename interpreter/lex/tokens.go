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
	TokenNumber   TokenType = iota
	TokenString   TokenType = iota
	TokenPair     TokenType = iota
	TokenTag      TokenType = iota
	TokenOperator TokenType = iota
	TokenCommand  TokenType = iota
	TokenWord     TokenType = iota
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
  case TokenPair:
    return fmt.Sprint("Pair: ", t.Value)
  case TokenTag:
    return fmt.Sprint("Tag: ", t.Value)
	case TokenNumber:
		return fmt.Sprint("Number: ", t.Value)
	case TokenOperator:
		return fmt.Sprint("Operator: ", t.Value)
	case TokenCommand:
		return fmt.Sprint("Command: ", t.Value)
	case TokenWord:
		return fmt.Sprint("Word: ", t.Value)
	default:
		fmt.Println(t.Type)
		panic("Unknown token type")
	}
}
