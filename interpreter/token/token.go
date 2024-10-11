package token

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
	Eof   TokenType = iota // Raised when the end of the input is reached
	Error TokenType = iota // Raised when an error occurs

	// TOKENS
	Number     TokenType = iota // 1
	String     TokenType = iota // "Helllo World"
	Plus       TokenType = iota // +
	Minus      TokenType = iota // -
	Star       TokenType = iota // *
	Slash      TokenType = iota // /
	Colon      TokenType = iota // :
	Key        TokenType = iota // thisbit:<expression>
	LeftParen  TokenType = iota // (
	RightParen TokenType = iota // )
	Tag        TokenType = iota // +HOME, -HOME
	Command    TokenType = iota // add, modify, etc.
	LessThan         TokenType = iota // <
	Equal         TokenType = iota // =
	Or         TokenType = iota // or
	And        TokenType = iota // and
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
	case Eof:
		return "EOF"
	case Error:
		return fmt.Sprint("Error: ", t.Value)
	case String:
		return fmt.Sprint("String: ", t.Value)
	case Key:
		return fmt.Sprint("Key: ", t.Value)
	case Tag:
		return fmt.Sprint("Tag: ", t.Value)
	case Number:
		return fmt.Sprint("Number: ", t.Value)
	case Plus:
		return fmt.Sprint("Plus: ", t.Value)
	case Minus:
		return fmt.Sprint("Minus: ", t.Value)
	case Star:
		return fmt.Sprint("Star: ", t.Value)
	case Slash:
		return fmt.Sprint("Slash: ", t.Value)
	case Command:
		return fmt.Sprint("Command: ", t.Value)
	case Colon:
		return fmt.Sprint("Colon: ", t.Value)
	case LeftParen:
		return fmt.Sprint("LeftParen: ", t.Value)
	case RightParen:
		return fmt.Sprint("RightParen: ", t.Value)
	case LessThan:
		return fmt.Sprint("LT: ", t.Value)
	case Equal:
		return fmt.Sprint("EQ: ", t.Value)
	case Or:
		return fmt.Sprint("Or: ", t.Value)
	case And:
		return fmt.Sprint("And: ", t.Value)
	default:
		var err = fmt.Errorf("Unknown token type: %d", t.Type)
		panic(err)
	}
}

func (t *TokenType) String() string {
	switch *t {
	case Eof:
		return "EOF"
	case Error:
		return "Error"
	case String:
		return "String"
	case Key:
		return "Key"
	case Tag:
		return "Tag"
	case Number:
		return "Number"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"
	case Command:
		return "Command"
	case Colon:
		return "Colon"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LessThan:
		return "LT"
	case Equal:
		return "EQ"
	case Or:
		return "Or"
	case And:
		return "And"
	default:
		var err = fmt.Errorf("Unknown token type: %d", *t)
		panic(err)
	}
}
