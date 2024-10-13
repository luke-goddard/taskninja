package core

import (
	"os"

	"github.com/luke-goddard/taskninja/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// This is going to be the main entry point for the application
type Runner struct {
	args   []string
	config *config.Config
}

func NewRunner(args []string) *Runner {
	return &Runner{args: args}
}

func (r *Runner) Run() {
	r.configDefaultLogger()
	r.loadConfigOrFail()
	r.config.InitLogger()
}

func (r *Runner) configDefaultLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func (r *Runner) loadConfigOrFail() {
	var conf, err = config.GetConfig()
	if err != nil && err.CanBootstrap() {
		conf = config.Bootstrap()
	}
	r.config = conf
}
