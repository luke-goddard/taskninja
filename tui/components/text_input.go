package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/tui/utils"
)

type TextInput struct {
	dimensions *utils.TerminalDimensions
	txtInput   textinput.Model
}

func NewTextInput(dimensions *utils.TerminalDimensions) *TextInput {
	var txtIn = textinput.New()
	txtIn.Placeholder = "Type here..."
	txtIn.Focus()
	txtIn.Width = dimensions.Width.Percent(0.95)
	return &TextInput{
		dimensions: dimensions,
		txtInput:   txtIn,
	}
}

func (t *TextInput) Update(msg tea.Msg) (*TextInput, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return t, tea.Quit
		}
	}

	t.txtInput, cmd = t.txtInput.Update(msg)
	return t, cmd
}

func (t *TextInput) View() string {
	return fmt.Sprintf(
		"What’s your favorite Pokémon?\n\n%s\n\n%s",
		t.txtInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (t *TextInput) Init() tea.Cmd {
	return textinput.Blink
}
