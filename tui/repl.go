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
	table      table.Model
	dimensions *TerminalDimensions
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

func NewTui() error {

	var dimensions, err = NewTerminalDimensions()
	if err != nil {
		return err
	}

	// https://github.com/charmbracelet/bubbletea/issues/43
	var columns = []table.Column{
		{Title: "ID", Width: dimensions.Width.PercentOrMin(0.02, 4)},
		{Title: "Age", Width: dimensions.Width.PercentOrMin(0.02, 4)},
		{Title: "Name", Width: dimensions.Width.PercentOrMin(0.54, 10)},
		{Title: "Priority", Width: dimensions.Width.PercentOrMin(0.06, 10)},
		{Title: "Project", Width: dimensions.Width.PercentOrMin(0.137, 10)},
		{Title: "Tags", Width: dimensions.Width.PercentOrMin(0.16, 5)},
	}

	var rows = []table.Row{
		{"1", "23", "Buy groceries", "High", "Shopping", "food"},
	}
	var tbl = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(dimensions.Height.PercentOrMin(0.7, 10)),
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
	return nil
}
