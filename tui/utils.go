package tui

import (
	"math"

	"github.com/gdamore/tcell"
)

type Dimension int

type TerminalDimensions struct {
	Width  Dimension
	Height Dimension
}

func NewTerminalDimensions() (*TerminalDimensions, error) {
	var screen, err = tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	defer screen.Fini()
	var width, height = screen.Size()

	return &TerminalDimensions{
		Width:  Dimension(width),
		Height: Dimension(height),
	}, nil
}

func (d *Dimension) Percent(p float64) int {
	return int(math.Floor(float64(*d) * p))
}

func (d *Dimension) PercentOrMax(p float64, max int) int {
	var percent = d.Percent(p)
	if percent > max {
		return max
	}
	return percent
}

func (d *Dimension) PercentOrMin(p float64, min int) int {
	var percent = d.Percent(p)
	if percent < min {
		return min
	}
	return percent
}
