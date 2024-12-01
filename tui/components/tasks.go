package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rs/zerolog/log"

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
	TableColumnStarted
	TableColumnName
	TableColumnAge
	TableColumnPriority
	TableColumnProject
	TableColumnTags
	TableColumnDependencies
	TableColumnUrgency
)

type TaskTable struct {
	Table                 table.Model
	baseStyle             lipgloss.Style
	dimensions            *utils.TerminalDimensions
	theme                 *utils.Theme
	bus                   *bus.Bus
	fuzzyFilter           string
	TaskIdsMatchingFilter []int64
	tableStyle            table.Styles
}

type TaskRow table.Row

const NOID int64 = -1

// ===========================================================================
// Task Table
// ===========================================================================

func NewTaskTable(baseStyle lipgloss.Style, dimensions *utils.TerminalDimensions, theme *utils.Theme, bus *bus.Bus) *TaskTable {
	assert.NotNil(bus, "bus is nil")
	assert.NotNil(baseStyle, "baseStyle is nil")
	assert.NotNil(dimensions, "dimensions is nil")
	assert.NotNil(theme, "theme is nil")
	var columns = []table.Column{
		{Title: "ID", Width: dimensions.Width.PercentOrMin(0.05, 0)},
		{Title: "Started", Width: dimensions.Width.PercentOrMin(0.1, 0)},
		{Title: "Name", Width: dimensions.Width.PercentOrMin(0.33, 0)},
		{Title: "Age", Width: dimensions.Width.PercentOrMin(0.05, 0)},
		{Title: "Priority", Width: dimensions.Width.PercentOrMin(0.06, 0)},
		{Title: "Project", Width: dimensions.Width.PercentOrMin(0.06, 0)},
		{Title: "Tags", Width: dimensions.Width.PercentOrMin(0.08, 0)},
		{Title: "Deps", Width: dimensions.Width.PercentOrMin(0.06, 0)},
		{Title: "Urgency", Width: dimensions.Width.PercentOrMin(0.21, 0)},
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
		Table:                 tbl,
		baseStyle:             baseStyle,
		bus:                   bus,
		theme:                 theme,
		fuzzyFilter:           "",
		TaskIdsMatchingFilter: []int64{},
		tableStyle:            style,
	}
}

func (m *TaskTable) Notify(e *events.Event) {
	// Little adapter to allow tea's interface to be compatible with the bus
	m.Update(e)
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
			if m.Table.Cursor() == len(m.Table.Rows())-1 {
				m.Table.SetCursor(m.Table.Cursor() - 1)
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
		case "H":
			m.bus.Publish(events.NewSetPriorityEvent(id, db.TaskPriorityHigh))
		case "M":
			m.bus.Publish(events.NewSetPriorityEvent(id, db.TaskPriorityMedium))
		case "L":
			m.bus.Publish(events.NewSetPriorityEvent(id, db.TaskPriorityLow))
		case "N":
			m.bus.Publish(events.NewSetPriorityEvent(id, db.TaskPriorityNone))
		case "n":
			m.markNextTaskAsNext()
		}
		m.Table, cmd = m.Table.Update(msg)
	case *events.Event:
		switch msg.Type {
		case events.EventListTaskResponse:
			m.handleListTasksResponse(events.DecodeListTasksResponseEvent(msg))
			return m, cmd
		case events.EventTableFuzzySearch:
			m.handleFuzzySearchResponse(events.DecodeTableFuzzySearch(msg))
			return m, cmd
		default:
		}
	}
	return m, cmd
}

func (m *TaskTable) GetIdForCurrentRow() int64 {
	var selectedRow = m.Table.SelectedRow()
	return TaskRow(selectedRow).ID()
}

func (m *TaskTable) CurrentTaskStarted() bool {
	var selectedRow = m.Table.SelectedRow()
	return TaskRow(selectedRow).Started()
}

func (m *TaskTable) GetCurrentRow() TaskRow {
	var selectedRow = m.Table.SelectedRow()
	return TaskRow(selectedRow)
}

func (m *TaskTable) GetRowAtPos(pos int) TaskRow {
	var selectedRow = m.Table.Rows()[pos]
	return TaskRow(selectedRow)
}

func (m *TaskTable) markNextTaskAsNext() {
	var selectedRow = m.Table.SelectedRow()
	var id = TaskRow(selectedRow).ID()
	var cmd = fmt.Sprintf("next %d", id)
	m.bus.Publish(events.NewRunProgramEvent(cmd))
}

func (m *TaskTable) handleListTasksResponse(e *events.ListTasksResponse) {
	var rows = []table.Row{}
	var ids = []int64{}
	var index = 0

	for _, task := range e.Tasks {
		if m.fuzzyFilter != "" {
			if !fuzzy.MatchFold(m.fuzzyFilter, task.Title) {
				continue
			}
		}
		ids = append(ids, task.ID)
		var columns = []string{}
		var started = ""
		var id = fmt.Sprintf("%d", task.ID)
		if task.State == db.TaskStateStarted {
			started = task.PrettyCumTime()
			id = fmt.Sprintf("%s-⏰", id)
		}

		var urgency = task.UrgencyStr()
		var urgencyStyle = lipgloss.
			NewStyle().
			Background(lipgloss.Color(task.UrgencyColourAnsiBackground())).
			Foreground(lipgloss.Color(task.UrgencyColourAnsiForeground()))

		var priority = task.PriorityStr()
		if task.Priority == db.TaskPriorityNone {
			priority = "❌"
		}
		if task.Next {
			priority += " ⭐"
		}

		urgency = urgencyStyle.Render(urgency)

		columns = append(columns, id)                       // ID
		columns = append(columns, started)                  // STARTED
		columns = append(columns, task.Title)               // NAME
		columns = append(columns, task.AgeStr())            // AGE
		columns = append(columns, priority)                 // PRIORITY
		columns = append(columns, task.ProjectNames.String) // PROJECT
		columns = append(columns, "")                       // TAGS
		columns = append(columns, task.Dependencies.String) // DEPENDENCIES
		columns = append(columns, urgency)                  // URGENCY

		index++
		rows = append(rows, columns)
	}
	m.TaskIdsMatchingFilter = ids
	m.Table.SetRows(rows)
}

func (m *TaskTable) handleFuzzySearchResponse(e *events.TableFuzzySearch) {
	m.fuzzyFilter = e.Match
	m.Table.SetCursor(0)
	log.Info().Str("filter", m.fuzzyFilter).Msg("Fuzzy filter")
}

func (m TaskTable) View() string {
	return m.baseStyle.Render(m.Table.View()) + "\n"
}

func (m TaskTable) Init() tea.Cmd {
	return nil
}

func (m TaskTable) HelpView() string {
	return m.Table.HelpView()
}

// ===========================================================================
// TaskRow
// ===========================================================================

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

func (r TaskRow) Title() string {
	assert.NotNil(r, "r is nil")
	if len(r) == 0 {
		return ""
	}
	if len(r) <= TableColumnName {
		return ""
	}
	return r[TableColumnName]
}

func (r TaskRow) UrgencyStr() string {
	assert.NotNil(r, "r is nil")
	if len(r) == 0 {
		return ""
	}
	if len(r) <= TableColumnUrgency {
		return ""
	}
	return r[TableColumnUrgency]
}

func (r TaskRow) PriorityStr() string {
	assert.NotNil(r, "r is nil")
	if len(r) == 0 {
		return ""
	}
	if len(r) <= TableColumnPriority {
		return ""
	}
	return r[TableColumnPriority]
}
