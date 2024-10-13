package main

import (
	"os"

	"github.com/luke-goddard/taskninja/core"
)

func main() {
	var runner = core.NewRunner(os.Args)
	runner.Run()
	// var args []string
	// if len(os.Args) > 1 {
	// 	args = os.Args[1:]
	// }
	// var sb = strings.Builder{}
	// for _, arg := range args {
	// 	sb.WriteString(arg)
	// }
	// var command = sb.String()
	// var interpreter = interpreter.NewInterpreter()
	// interpreter.Execute(command)
}
