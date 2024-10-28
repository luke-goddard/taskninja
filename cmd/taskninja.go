package main

import (
	"os"

	"github.com/luke-goddard/taskninja/core"
)

func main() {
	var runner = core.NewRunner(os.Args)
	runner.Run()
}
