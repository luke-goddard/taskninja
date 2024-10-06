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
	var command = strings.Join(args, " ")
	var interpreter = interpreter.NewInterpreter()
	interpreter.Execute(command)
}
