package main

import (
	"os"
	"strings"

	"github.com/luke-goddard/taskninja/interpreter"
)

func main() {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}
	var sb = strings.Builder{}
	for _, arg := range args {
		if strings.Contains(arg, " ") {
			sb.WriteString("\"")
			sb.WriteString(arg)
			sb.WriteString("\"")
		}
		sb.WriteString(arg)
	}
	var command = sb.String()
	var interpreter = interpreter.NewInterpreter(command)
	interpreter.Execute()
}
