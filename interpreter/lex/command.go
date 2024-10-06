package lex

type Command string

const (
	CommandAdd    Command = "add"
	CommandAll    Command = "all"
	CommandDelete Command = "delete"
	CommandDone   Command = "done"
	CommandList   Command = "list"
	CommandModify Command = "modify"
	CommandReady  Command = "ready"
	CommandStart  Command = "start"
	CommandStop   Command = "stop"
	CommandTags   Command = "tags"
)

func lexCommand(l *Lexer) StateFn {
	var last = l.readUntil(func(r rune) bool {
		return IsWhitespace(r) || r == EOF || r == ':'
	})
	if last == ':' {
		return lexPair
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

		l.emit(TokenCommand)
		return lexStart
	}

	return lexWord
}
