package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/tui/components"
	"github.com/luke-goddard/taskninja/tui/utils"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	tabs       []string
	table      *components.TaskTable
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
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}
	var newTable, _ = m.table.Update(msg)
	m.table = newTable
	return m, cmd
}

func (m model) View() string {
	var document strings.Builder
	var renderedTabs []string
	for i, tab := range m.tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.tabs)-1, i == m.activeTab
		if isActive {
			style = activeTabStyle
		} else {
			style = inactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(tab))
	}
	document.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...))
	document.WriteString("\n")
	// document.WriteString(baseStyle.Render(m.table.View()) + "\n")
	document.WriteString(m.table.View() + "\n")
	document.WriteString(m.table.HelpView() + "\n")
	return document.String()
}

func (m model) Init() tea.Cmd { return nil }

func NewTui() error {

	var dimensions, err = utils.NewTerminalDimensions()
	if err != nil {
		return err
	}

	var theme = utils.NewTheme()

	var tabs = []string{"Tasks", "Projects", "Tags", "Settings"}
	var model = model{
		table:      components.NewTaskTable(baseStyle, dimensions, theme),
		tabs:       tabs,
		dimensions: dimensions,
	}
	tea.NewProgram(model, tea.WithAltScreen()).Run()
	return nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
