package lex

// Represent different types of commands e.g add
type Command string

const (
	CommandAdd    Command = "add"    // Add a new task
	CommandAll    Command = "all"    // List all tasks
	CommandDelete Command = "delete" // Delete a task
	CommandDone   Command = "done"   // Mark a task as done
	CommandList   Command = "list"   // List tasks
	CommandModify Command = "modify" // Modify a task
	CommandReady  Command = "ready"  // Mark a task as ready
	CommandStart  Command = "start"  // Start a task
	CommandStop   Command = "stop"   // Stop a task
	CommandTags   Command = "tags"   // List all tags
)

func lexCommand(l *Lexer) StateFn {
	var last = l.readUntil(func(r rune) bool {
		return IsWhitespace(r) || r == EOF || r == ':'
	})
	if last == ':' {
		return lexPair
	}

	if l.seenCommand {
		return lexWord
	}

	var lexeme = l.current()

	if lexeme == string(CommandAdd) ||
		lexeme == string(CommandAll) ||
		lexeme == string(CommandDelete) ||
		lexeme == string(CommandDone) ||
		lexeme == string(CommandList) ||
		lexeme == string(CommandModify) ||
		lexeme == string(CommandReady) ||
		lexeme == string(CommandStart) ||
		lexeme == string(CommandStop) ||
		lexeme == string(CommandTags) {

		l.seenCommand = true
		l.emit(TokenCommand)
		return lexStart
	}

	return lexWord
}
