package utils

import "github.com/charmbracelet/lipgloss"

type Colour string

const (
	DEFAULT_PRIMARY_COLOUR    lipgloss.Color = lipgloss.Color("57")
	DEFAULT_SECONDARY_COLOUR  lipgloss.Color = lipgloss.Color("54")
	DEFAULT_FOREGROUND_COLOUR lipgloss.Color = lipgloss.Color("229")
	DEFAULT_WARNING_COLOUR    lipgloss.Color = lipgloss.Color("208")
	DEFAULT_ERROR_COLOUR      lipgloss.Color = lipgloss.Color("196")
)

type Theme struct {
	PrimaryColor   lipgloss.Color
	ForgroundColor lipgloss.Color
	SecondaryColor lipgloss.Color
	WarningColor   lipgloss.Color
	DangerColor    lipgloss.Color
}

func NewTheme() *Theme {
	return &Theme{
		PrimaryColor:   DEFAULT_PRIMARY_COLOUR,
		SecondaryColor: DEFAULT_SECONDARY_COLOUR,
		ForgroundColor: DEFAULT_FOREGROUND_COLOUR,
		WarningColor:   DEFAULT_WARNING_COLOUR,
		DangerColor:    DEFAULT_ERROR_COLOUR,
	}
}

func (t *Theme) SetPrimaryColor(c lipgloss.Color) *Theme {
	t.PrimaryColor = c
	return t
}

func (t *Theme) SetSecondaryColor(c lipgloss.Color) *Theme {
	t.SecondaryColor = c
	return t
}

func (t *Theme) SetWarningColor(c lipgloss.Color) *Theme {
	t.WarningColor = c
	return t
}

func (t *Theme) SetDangerColor(c lipgloss.Color) *Theme {
	t.DangerColor = c
	return t
}

func (t *Theme) SetForgroundColor(c lipgloss.Color) *Theme {
	t.ForgroundColor = c
	return t
}
