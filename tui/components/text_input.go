package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/assert"
	"github.com/luke-goddard/taskninja/bus"
	"github.com/luke-goddard/taskninja/events"
	"github.com/luke-goddard/taskninja/tui/utils"
	"github.com/rs/zerolog/log"
)

const RINGBUFFER_SIZE = 30

type InputMode int

const (
	InputModeNone InputMode = iota
	InputModeCmd
	InputModeSearch
)

type TextInput struct {
	dimensions *utils.TerminalDimensions
	txtInput   textinput.Model
	enabled    bool
	bus        *bus.Bus
	err        *error
	history    *InputHistoryRingBuffer
	inputMode  InputMode
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

func (t *TextInput) CanQuit() bool {
	return t.Disabled()
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
		history:    NewInputHistoryRingBuffer(),
		inputMode:  InputModeNone,
	}
}

func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	var cmd tea.Cmd
	var enabled = t.Enabled()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			t.ClearErr()
			if t.Enabled() {
				t.Disable()
				t.txtInput.Blur()
				t.handleSubmit()
				t.inputMode = InputModeNone
			}
			return t, cmd
		case tea.KeyEscape:
			if t.Enabled() {
				t.Disable()
				t.txtInput.Blur()
			}
		case tea.KeyUp:
			if t.enabled {
				var previous = t.history.GetPrevious()
				t.txtInput.SetValue(previous)
			}
		case tea.KeyDown:
			if t.enabled {
				var next = t.history.GetNext()
				t.txtInput.SetValue(next)
			}
		}

		switch msg.String() {
		case "a":
			if !enabled {
				t.inputMode = InputModeCmd
				t.ClearErr()
				t.Enable()
				t.txtInput.Focus()
				t.txtInput.SetValue("add \"")
				return t, cmd
			}
		case "/":
			if !enabled {
				t.inputMode = InputModeSearch
				t.ClearErr()
				t.Enable()
				t.txtInput.Focus()
				t.txtInput.SetValue("")
				return t, cmd
			}
		}
	case *events.Event:
		switch msg.Type {
		case events.EventError:
			var err = msg.Data.(error)
			t.err = &err
			t.enabled = false
			return t, cmd
		case events.EventRunProgram:
			return t, cmd
		}
	}

	// check if tea.KeyMsg is a valid message
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if enabled {
			t.txtInput, cmd = t.txtInput.Update(msg)
			if t.inputMode == InputModeSearch {
				log.Info().Interface("msg", msg).Msg("TextInput Update")
				t.submitSearch()
			}
		}
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
	if t.inputMode == InputModeCmd {
		return fmt.Sprintf(
			"Enter a command:\n\n%s\n",
			t.txtInput.View(),
		) + "\n"
	}
	if t.inputMode == InputModeSearch {
		return fmt.Sprintf(
			"Search:\n\n%s\n",
			t.txtInput.View(),
		) + "\n"
	}
	return fmt.Sprintf(
		"Enter a comamand:\n\n%s\n",
		t.txtInput.View(),
	) + "\n"
}

func (t *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (t *TextInput) handleSubmit() {
	log.Info().Str("inputMode", fmt.Sprintf("%d", t.inputMode)).Msg("Submitting input")
	switch t.inputMode {
	case InputModeNone:
		return
	case InputModeCmd:
		t.submitProgram()
		t.txtInput.SetValue("")
	case InputModeSearch:
		t.submitSearch()
		t.txtInput.SetValue("")
	default:
		log.Error().Str("inputMode", fmt.Sprintf("%d", t.inputMode)).Msg("Unknown input mode")
		assert.Fail("Unknown input mode")
	}
}

func (t *TextInput) submitProgram() {
	log.Info().Str("program", t.txtInput.Value()).Msg("Submitting program")
	t.err = nil
	var program = t.txtInput.Value()
	t.history.Add(program)
	t.bus.Publish(events.NewRunProgramEvent(program))
	t.txtInput.SetValue("")
}

func (t *TextInput) submitSearch() {
	log.Info().Str("search", t.txtInput.Value()).Msg("Submitting search")
	t.err = nil
	var search = t.txtInput.Value()
	t.bus.Publish(events.NewTableFuzzySearch(search))
	t.bus.Publish(events.NewListTasksEvent())
}

func (t *TextInput) ClearErr() {
	t.err = nil
}

type InputHistoryRingBuffer struct {
	history     [RINGBUFFER_SIZE]string
	lookupIndex int
	insertIndex int
}

func NewInputHistoryRingBuffer() *InputHistoryRingBuffer {
	return &InputHistoryRingBuffer{
		history:     [RINGBUFFER_SIZE]string{},
		lookupIndex: 0,
		insertIndex: 0,
	}
}

func (r *InputHistoryRingBuffer) Add(input string) {
	r.history[r.insertIndex] = input
	r.insertIndex++
	if r.insertIndex >= RINGBUFFER_SIZE {
		r.insertIndex = 0
	}
	if r.insertIndex == r.lookupIndex {
		r.lookupIndex++
		if r.lookupIndex >= RINGBUFFER_SIZE {
			r.lookupIndex = 0
		}
	}
	r.lookupIndex = r.insertIndex
}

func (r *InputHistoryRingBuffer) Get(index int) string {
	return r.history[index]
}

func (r *InputHistoryRingBuffer) GetPrevious() string {
	r.lookupIndex--
	if r.lookupIndex < 0 {
		r.lookupIndex = RINGBUFFER_SIZE - 1
	}
	var cmd = r.history[r.lookupIndex]
	if cmd == "" {
		r.lookupIndex = RINGBUFFER_SIZE - 1
	}
	return r.history[r.lookupIndex]
}

func (r *InputHistoryRingBuffer) GetNext() string {
	r.lookupIndex++
	if r.lookupIndex >= RINGBUFFER_SIZE {
		r.lookupIndex = 0
	}
	var cmd = r.history[r.lookupIndex]
	if cmd == "" {
		r.lookupIndex = 0
	}
	return r.history[r.lookupIndex]
}
