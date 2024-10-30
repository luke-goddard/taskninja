package core

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/config"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/bus/handler"
	"github.com/luke-goddard/taskninja/interpreter"
	"github.com/luke-goddard/taskninja/services"
	"github.com/luke-goddard/taskninja/tui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// This is going to be the main entry point for the application
type Runner struct {
	bus         *bus.Bus
	service     *services.ServiceHandler
	handler     *handler.EventHandler
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
		bus:         bus.NewBus(),
		args:        normalizeArgs(args),
		interpreter: interpreter.NewInterpreter(),
	}
}

func (r *Runner) Run() {
	var err error
	var store *db.Store
	var program *tea.Program

	r.loadConfigOrFail()
	store, err = db.NewStore(&r.config.Connection)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database store")
	}

	r.store = store
	r.service = services.NewServiceHandler(r.interpreter, r.store)
	r.handler = handler.NewEventHandler(r.service, r.bus)
	r.configDefaultLogger()
	r.config.InitLogger()
	r.bus.Subscribe(r.handler)

	program, err = tui.NewTui(r.bus)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create TUI")
	}
	program.Run()
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
