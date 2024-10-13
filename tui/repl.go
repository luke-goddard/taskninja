package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func (m model) Init() tea.Cmd { return nil }

func NewTui() {
	var columns = []table.Column{
		{Title: "ID", Width: 2},
		{Title: "Name", Width: 30},
		{Title: "Priority", Width: 10},
		{Title: "Project", Width: 10},
		{Title: "Tags", Width: 10},
	}
	var rows = []table.Row{
		{"1", "Task 1", "1", "Project 1", "tag1, tag2"},
		{"2", "Task 2", "2", "Project 2", "tag3, tag4"},
		{"3", "Task 3", "3", "Project 3", "tag5, tag6"},
		{"4", "Task 4", "4", "Project 4", "tag7, tag8"},
	}
	var tbl = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	var style = table.DefaultStyles()

	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF00FF")).
		BorderBottom(true).
		Bold(true)

	style.Selected = style.Selected.
		Foreground(lipgloss.Color("205")).
		Background(lipgloss.Color("240")).
		Bold(true)

	tbl.SetStyles(style)

	var model = model{table: tbl}
	tea.NewProgram(model, tea.WithAltScreen()).Run()
}
