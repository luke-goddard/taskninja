package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

type Tabs struct {
	Tabs      []string
	ActiveTab int
}

func NewTabs() *Tabs {
	return &Tabs{
		ActiveTab: 0,
		Tabs: []string{
			"Tasks (1)",
			"Projects (2)",
			"Tags (3)",
			"People (4)",
			"Context (5)",
			"Study (6)",
			"Notes (7)",
			"Settings (8)",
		},
	}
}

func (model Tabs) View() string {
	var document strings.Builder
	var renderedTabs []string
	for i, tab := range model.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(model.Tabs)-1, i == model.ActiveTab
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
	return document.String()
}

func (m Tabs) Init() tea.Cmd { return nil }

func (m Tabs) Update(msg tea.Msg) (*Tabs, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return &m, tea.Quit
		case "right", "l", "tab":
			m.ActiveTab = min(m.ActiveTab+1, len(m.Tabs)-1)
			return &m, nil
		case "left", "h", "p", "shift+tab":
			m.ActiveTab = max(m.ActiveTab-1, 0)
			return &m, nil
		}
	}
	return &m, cmd
}
