package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/tui/utils"
)

type TaskTable struct {
	table      table.Model
	baseStyle  lipgloss.Style
	dimensions *utils.TerminalDimensions
	theme      *utils.Theme
}

func NewTaskTable(baseStyle lipgloss.Style, dimensions *utils.TerminalDimensions, theme *utils.Theme) *TaskTable {
	assert.NotNil(baseStyle, "baseStyle is nil")
	assert.NotNil(dimensions, "dimensions is nil")
	assert.NotNil(theme, "theme is nil")
	var columns = []table.Column{
		{Title: "ID", Width: dimensions.Width.PercentOrMin(0.02, 4)},
		{Title: "Age", Width: dimensions.Width.PercentOrMin(0.02, 4)},
		{Title: "Name", Width: dimensions.Width.PercentOrMin(0.54, 10)},
		{Title: "Priority", Width: dimensions.Width.PercentOrMin(0.06, 10)},
		{Title: "Project", Width: dimensions.Width.PercentOrMin(0.137, 10)},
		{Title: "Tags", Width: dimensions.Width.PercentOrMin(0.16, 5)},
	}

	var rows = []table.Row{
		// {"1", "23", "Buy groceries", "High", "Shopping", "food"},
	}
	var tbl = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(dimensions.Height.PercentOrMin(0.6, 10)),
	)

	var style = table.DefaultStyles()

	style.Header = style.Header.
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(theme.PrimaryColor).
		BorderBottom(true).
		Bold(true)

	style.Selected = style.Selected.
		Foreground(utils.DEFAULT_FOREGROUND_COLOUR).
		Background(utils.DEFAULT_PRIMARY_COLOUR).
		Bold(true)

	tbl.SetStyles(style)
	return &TaskTable{
		table:     tbl,
		baseStyle: baseStyle,
	}
}

func (m *TaskTable) Update(msg tea.Msg) (*TaskTable, tea.Cmd) {
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
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
		m.table, cmd = m.table.Update(msg)
	case *events.Event:
		switch msg.Type {
		case events.EventListTaskResponse:
			m.handleListTasksResponse(events.DecodeListTasksResponseEvent(msg))
			return m, cmd
		}
	}
	return m, cmd
}

func (m *TaskTable) handleListTasksResponse(e *events.ListTasksResponse) {
	var rows = []table.Row{}
	for _, task := range e.Tasks {
		var columns = []string{}
		columns = append(columns, fmt.Sprintf("%d", task.ID))
		columns = append(columns, "")
		columns = append(columns, task.Title)
		columns = append(columns, "")
		columns = append(columns, "")

		rows = append(rows, columns)
	}
	m.table.SetRows(rows)
}

func (m TaskTable) View() string {
	return m.baseStyle.Render(m.table.View()) + "\n"
}

func (m TaskTable) Init() tea.Cmd {
	return nil
}

func (m TaskTable) HelpView() string {
	return m.table.HelpView()
}
