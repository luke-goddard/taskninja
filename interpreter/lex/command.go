package lex

import (
	"strings"

	"github.com/luke-goddard/taskninja/interpreter/token"
)

// Represent different types of commands e.g add
type Command string

const (
	CommandAdd     Command = "add"     // Add a new task
	CommandDepends Command = "depends" // Add a dependency
	// CommandAll    Command = "all"    // List all tasks
	// CommandDelete Command = "delete" // Delete a task
	// CommandDone   Command = "done"   // Mark a task as done
	// CommandList   Command = "list"   // List tasks
	// CommandModify Command = "modify" // Modify a task
	// CommandReady  Command = "ready"  // Mark a task as ready
	// CommandStart  Command = "start"  // Start a task
	// CommandStop   Command = "stop"   // Stop a task
	// CommandTags   Command = "tags"   // List all tags
)

func lexCommand(l *Lexer) StateFn {
	var last = l.readUntil(func(r rune) bool {
		return !IsAlphabet(r)
	})
	if last == ':' {
		return lexPair
	}

	var lexeme = l.current()

	if lexeme == string(CommandAdd) ||
		lexeme == string(CommandDepends) {
		if !l.seenCommand {
			l.seenCommand = true
			l.emit(token.Command)
			return lexStart
		}
	}

	if strings.ToLower(lexeme) == "or" {
		l.emit(token.Or)
		return lexStart
	}

	if strings.ToLower(lexeme) == "and" {
		l.emit(token.And)
		return lexStart
	}

	return lexWord
}
