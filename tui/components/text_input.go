package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/tui/utils"
)

type TextInput struct {
	dimensions *utils.TerminalDimensions
	txtInput   textinput.Model
	enabled    bool
	bus        *bus.Bus
	err        *error
}

func (t *TextInput) Enable() {
	t.enabled = true
}

func (t *TextInput) Disable() {
	t.enabled = false
}

func (t *TextInput) Enabled() bool {
	return t.enabled
}

func (t *TextInput) Disabled() bool {
	return !t.enabled
}

func NewTextInput(dimensions *utils.TerminalDimensions, bus *bus.Bus) *TextInput {
	assert.NotNil(dimensions, "TerminalDimensions is nil")
	assert.NotNil(bus, "Bus is nil")
	var txtIn = textinput.New()
	txtIn.Placeholder = "Type here..."
	txtIn.Focus()
	txtIn.Width = dimensions.Width.Percent(0.95)
	return &TextInput{
		dimensions: dimensions,
		txtInput:   txtIn,
		enabled:    false,
		bus:        bus,
	}
}

func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	var cmd tea.Cmd
	var enabled = t.Enabled()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return t, tea.Quit
		case tea.KeyEnter:
			if t.Enabled() {
				t.Disable()
				t.txtInput.Blur()
				t.submitProgram()
			}
		}

		switch msg.String() {
		case "a":
			t.Enable()
			t.txtInput.Focus()
			if !enabled {
				t.txtInput.SetValue("add \"")
			}
		}
	case *events.Event:
		switch msg.Type {
		case events.EventError:
			var err = msg.Data.(error)
			t.err = &err
			t.enabled = false
		case events.EventRunProgram:
			t.err = nil
		}
	}

	if enabled {
		t.txtInput, cmd = t.txtInput.Update(msg)
	}
	return t, cmd
}

func (t *TextInput) View() string {
	if t.err != nil && !t.enabled {
		return fmt.Sprintf("Error: %s\n", (*t.err).Error())
	}
	if !t.enabled {
		return ""
	}
	return fmt.Sprintf(
		"Enter a test to add\n\n%s\n",
		t.txtInput.View(),
	) + "\n"
}

func (t *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TextInput) submitProgram() {
	var program = t.txtInput.Value()
	t.bus.Publish(events.NewRunProgramEvent(program))
	t.txtInput.SetValue("")
}
