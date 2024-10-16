package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/tui/components"
	"github.com/luke-goddard/taskninja/tui/utils"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	tabs       *components.Tabs
	table      *components.TaskTable
	input      *components.TextInput
	dimensions *utils.TerminalDimensions
	activeTab  int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var newTable, _ = m.table.Update(msg)
	m.table = newTable

	var newTabs, _ = m.tabs.Update(msg)
	m.tabs = newTabs

	var newInput, _ = m.input.Update(msg)
	m.input = newInput

	return m, cmd
}

func (m model) View() string {
	var document strings.Builder
	document.WriteString(m.tabs.View() + "\n")
	document.WriteString(m.table.View() + "\n")
	document.WriteString(m.table.HelpView() + "\n")
	document.WriteString(m.input.View() + "\n")
	return document.String()
}

func (m model) Init() tea.Cmd { return nil }

func NewTui() error {

	var dimensions, err = utils.NewTerminalDimensions()
	if err != nil {
		return err
	}

	var theme = utils.NewTheme()

	var tabs = components.NewTabs()
	var model = model{
		input:      components.NewTextInput(dimensions),
		table:      components.NewTaskTable(baseStyle, dimensions, theme),
		tabs:       tabs,
		dimensions: dimensions,
	}

	tea.NewProgram(model, tea.WithAltScreen()).Run()
	return nil
}
