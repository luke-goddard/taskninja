package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/tui/components"
	"github.com/luke-goddard/taskninja/tui/utils"
)

var baseStyle = lipgloss.
	NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	tabs       *components.Tabs
	bus        *bus.Bus
	table      *components.TaskTable
	input      *components.TextInput
	doughnut   *components.Doughnut
	dimensions *utils.TerminalDimensions
	activeTab  int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if m.input.CanQuit() {
				return m, tea.Quit
			}
		}
	case *events.Event:
		var newTable, _ = m.table.Update(msg)
		m.table = newTable

		var newTabs, _ = m.tabs.Update(msg)
		m.tabs = newTabs
	}

	if m.input.Disabled() {
		var newTable, _ = m.table.Update(msg)
		m.table = newTable

		var newTabs, _ = m.tabs.Update(msg)
		m.tabs = newTabs
	}

	var newInput *components.TextInput
	newInput, _ = m.input.Update(msg)
	m.input = newInput

	var newDoughnut *components.Doughnut
	newDoughnut, cmd = m.doughnut.Update(msg)
	m.doughnut = newDoughnut

	return m, cmd
}

func (m model) View() string {
	var document strings.Builder
	document.WriteString(m.tabs.View() + "\n")
	if m.tabs.ActiveTab == 2 {
		document.WriteString("\n")
		document.WriteString(m.doughnut.View() + "\n")
	} else {
		document.WriteString(m.table.View() + "\n")
		document.WriteString(m.table.HelpView() + "\n")
		document.WriteString(m.input.View() + "\n")
	}
	return document.String()
}

func (m model) Notify(e *events.Event) {
	// Little adapter to allow tea's interface to be compatible with the bus
	m.Update(e)
}

func (m model) Init() tea.Cmd {

	m.bus.Subscribe(m)
	m.bus.Publish(events.NewListTasksEvent())
	go m.RefreshTaskListProgramatically()

	return tea.Batch(
		m.table.Init(),
		m.tabs.Init(),
		m.input.Init(),
		m.doughnut.Init(),
	)
}

func (m model) RefreshTaskListProgramatically() {
	var sleep = time.Duration(10) * time.Second
	for {
		time.Sleep(sleep)
		if m.tabs.ActiveTab != 0 {
			continue
		}
		m.bus.Publish(events.NewListTasksEvent())
	}
}

func NewTui(bus *bus.Bus) (*tea.Program, error) {

	var dimensions, err = utils.NewTerminalDimensions()
	if err != nil {
		return nil, err
	}

	var theme = utils.NewTheme()

	var tabs = components.NewTabs()
	var model = model{
		bus:        bus,
		input:      components.NewTextInput(dimensions, bus),
		table:      components.NewTaskTable(baseStyle, dimensions, theme, bus),
		doughnut:   components.NewDonut(dimensions),
		tabs:       tabs,
		dimensions: dimensions,
	}

	var program = tea.NewProgram(model, tea.WithAltScreen())
	return program, nil
}
