package core

import (
	"os"
	"strings"

	"github.com/luke-goddard/taskninja/config"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/interpreter/ast"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// This is going to be the main entry point for the application
type Runner struct {
	args        string
	config      *config.Config
	interpreter *interpreter.Interpreter
	store       *db.Store
}

func normalizeArgs(args []string) string {
	var sb = strings.Builder{}
	for i, arg := range args {
		if i == 0 {
			continue
		}
		sb.WriteString(" ")
		sb.WriteString(arg)
	}
	return sb.String()
}

func NewRunner(args []string) *Runner {
	return &Runner{
		args:        normalizeArgs(args),
		interpreter: interpreter.NewInterpreter(),
	}
}

func (r *Runner) Run() {
	r.configDefaultLogger()
	r.loadConfigOrFail()
	r.config.InitLogger()

	var db, err = db.NewStore(&r.config.Connection)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create database store")
	}
	r.store = db

	var cmd, errs = r.interpreter.ParserString(r.args)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error().Msg(err.Error())
		}
		return
	}

	switch cmd.Kind {
	case ast.CommandKindList:
	default:
		log.Fatal().
			Str("kind", cmd.Kind.String()).
			Msg("unsupported command kind")
	}
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

func (r *Runner) loadDatabaseOrFail() {
}
