package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/db"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/tui/utils"
)

type TableRowIndex int

const (
	TableColumnID int = iota
	TableColumnUrgency
	TableColumnStarted
	TableColumnName
	TableColumnAge
	TableColumnPriority
	TableColumnProject
	TableColumnTags
)

type TaskTable struct {
	table      table.Model
	baseStyle  lipgloss.Style
	dimensions *utils.TerminalDimensions
	theme      *utils.Theme
	bus        *bus.Bus
}

type TaskRow table.Row

const NOID int64 = -1

func (r TaskRow) ID() int64 {
	assert.NotNil(r, "r is nil")
	if len(r) == 0 {
		return NOID
	}
	if len(r) <= TableColumnID {
		return NOID
	}
	assert.True(len(r) > TableColumnID, "r does not have a column for ID")
	var str = r[TableColumnID]
	str = strings.TrimSuffix(str, "-⏰")
	var id, err = strconv.ParseInt(str, 10, 64)
	assert.Nil(err, "failed to convert ID to int")
	return id
}

func (r TaskRow) Started() bool {
	assert.NotNil(r, "r is nil")
	if len(r) == 0 {
		return false
	}
	var started = r[TableColumnStarted]
	return started != ""
}

func NewTaskTable(baseStyle lipgloss.Style, dimensions *utils.TerminalDimensions, theme *utils.Theme, bus *bus.Bus) *TaskTable {
	assert.NotNil(bus, "bus is nil")
	assert.NotNil(baseStyle, "baseStyle is nil")
	assert.NotNil(dimensions, "dimensions is nil")
	assert.NotNil(theme, "theme is nil")
	var columns = []table.Column{
		{Title: "ID", Width: dimensions.Width.PercentOrMin(0.05, 4)},
		{Title: "Urgency", Width: dimensions.Width.PercentOrMin(0.08, 4)},
		{Title: "Started", Width: dimensions.Width.PercentOrMin(0.08, 4)},
		{Title: "Name", Width: dimensions.Width.PercentOrMin(0.43, 10)},
		{Title: "Age", Width: dimensions.Width.PercentOrMin(0.05, 4)},
		{Title: "Priority", Width: dimensions.Width.PercentOrMin(0.06, 10)},
		{Title: "Project", Width: dimensions.Width.PercentOrMin(0.134, 10)},
		{Title: "Tags", Width: dimensions.Width.PercentOrMin(0.08, 5)},
	}

	var rows = []table.Row{}
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
		bus:       bus,
		theme:     theme,
	}
}

func (m *TaskTable) Update(msg tea.Msg) (*TaskTable, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		var id = m.GetIdForCurrentRow()
		if id == NOID {
			break
		}
		switch msg.String() {
		case "d":
			m.bus.Publish(events.NewCompleteEvent(id))
		case "D":
			if m.table.Cursor() == len(m.table.Rows())-1 {
				m.table.SetCursor(m.table.Cursor() - 1)
			}
			m.bus.Publish(events.NewDeleteTaskByIdEvent(id))
		case "s":
			if !m.CurrentTaskStarted() {
				m.bus.Publish(events.NewStartTaskEvent(id))
			} else {
				m.bus.Publish(events.NewStopTaskByIdEvent(id))
			}
		case "+":
			m.bus.Publish(events.NewIncreasePriorityEvent(id))
		case "-":
			m.bus.Publish(events.NewDecreasePriorityEvent(id))
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

func (m *TaskTable) GetIdForCurrentRow() int64 {
	var selectedRow = m.table.SelectedRow()
	return TaskRow(selectedRow).ID()
}

func (m *TaskTable) CurrentTaskStarted() bool {
	var selectedRow = m.table.SelectedRow()
	return TaskRow(selectedRow).Started()
}

func (m *TaskTable) handleListTasksResponse(e *events.ListTasksResponse) {
	var rows = []table.Row{}
	for _, task := range e.Tasks {
		var columns = []string{}
		var started = ""
		var id = fmt.Sprintf("%d", task.ID)
		if task.State == db.TaskStateStarted {
			started = task.PrettyCumTime()
			id = fmt.Sprintf("%s-⏰", id)
		}
		columns = append(columns, id)                       // ID
		columns = append(columns, task.UrgencyStr())        // URGENCY
		columns = append(columns, started)                  // STARTED
		columns = append(columns, task.Title)               // NAME
		columns = append(columns, task.AgeStr())            // AGE
		columns = append(columns, task.PriorityStr())       // PRIORITY
		columns = append(columns, task.ProjectNames.String) // PROJECT
		columns = append(columns, "")                       // TAGS

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
