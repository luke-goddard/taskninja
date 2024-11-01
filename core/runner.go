package core

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/bus/handler"
	"github.com/luke-goddard/taskninja/config"
	"github.com/luke-goddard/taskninja/db"
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
	assert.NotNil(r.bus, "Bus is nil")
	assert.NotNil(r.interpreter, "Interpreter is nil")
	assert.NotNil(r.args, "Args is nil")

	var err error
	var store *db.Store
	var program *tea.Program

	r.loadConfigOrFail()
	r.configDefaultLogger()
	r.config.InitLogger()

	store, err = db.NewStore(&r.config.Connection)
	defer store.Close()

	assert.NoError(err, "Failed to create database store")
	assert.True(store.IsConnected(), "Store is not connected")

	r.store = store
	r.service = services.NewServiceHandler(r.interpreter, r.store)
	r.handler = handler.NewEventHandler(r.service, r.bus)
	r.bus.Subscribe(r.handler)

	program, err = tui.NewTui(r.bus)

	assert.NoError(err, "Failed to create TUI")
	assert.NotNil(program, "Program is nil")
	assert.NotNil(r.store, "Store is nil")
	assert.NotNil(r.service, "Service is nil")
	assert.NotNil(r.handler, "Handler is nil")
	assert.True(r.bus.HasSubscribers(), "Bus has no subscribers")

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
