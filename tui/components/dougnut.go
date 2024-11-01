// donut.go
//
// An implementation in Go of Andy Sloan's donut.c
//
// Original:
// http://www.a1k0n.net/2011/07/20/donut-math.html
//
// To run:
//    $ go mod init github.com/GaryBoone/GoDonut
//    $ go mod tidy
//    $ go run .
//
// Author:      gary.boone@gmail.com
// History:     20121012      • initial version
//              20121013      • upload to GitHub
//              20240322 	  • update run directions for Go 1.122.1

package components

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luke-goddard/taskninja/tui/utils"
)

var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

const frame_delay = 60
const theta_spacing = 0.07
const phi_spacing = 0.02

const R1 = 1.0
const R2 = 2.0
const K2 = 5.0

type Doughnut struct {
	screen     *Screen
	a          float64
	b          float64
	k1         float64
	dimensions *utils.TerminalDimensions
	Interval   time.Duration
}

type TickMsg struct{ id int }

func NewDonut(dimensions *utils.TerminalDimensions) *Doughnut {
	var w = int(dimensions.Width) - 15
	var h = int(dimensions.Height) - 15
	var s = min(w, h)
	screen := newScreen(s)
	A, B, K1 := 1.0, 1.0, float64(screen.dim)*K2*3.0/(8.0*(R1+R2))
	return &Doughnut{screen, A, B, K1, dimensions, time.Millisecond * frame_delay * 2}
}

type Screen struct {
	dim  int
	data [][]byte
}

func newZBuffer(d int) *[][]float64 {
	b := make([][]float64, d)
	for i := range b {
		b[i] = make([]float64, d)
	}
	return &b
}

func newScreen(d int) *Screen {
	b := make([][]byte, d)
	for i := range b {
		b[i] = make([]byte, d)
	}
	return &Screen{d, b}
}

func (screen Screen) render() string {
	var sb strings.Builder
	// fmt.S("\x1b[H") // bring cursor to "home" location
	for j := 0; j < screen.dim; j++ {
		sb.WriteString(fmt.Sprintf("%s\n", screen.data[j]))
	}
	return sb.String()
}

func (d *Doughnut) Init() tea.Cmd {
	return tea.Batch(func() tea.Msg { return TickMsg{id: nextID()} }, tick(d.Interval))
}

func (d *Doughnut) Update(msg tea.Msg) (*Doughnut, tea.Cmd) {
	switch msg.(type) {
	case TickMsg:
		return d, tick(d.Interval)
	}
	return d, nil
}

func (d *Doughnut) View() string {
	d.a += 0.07
	d.b += 0.03
	d.screen.computeFrame(d.a, d.b, d.k1)
	var sb strings.Builder
	// sb.WriteString("\x1b[H") // bring cursor to "home" location
	for j := 0; j < d.screen.dim; j++ {
		var data = fmt.Sprintf("%s\n", d.screen.data[j])
		sb.WriteString(data)
	}
	return sb.String()
}

func (screen *Screen) clear() {
	for i, _ := range screen.data {
		for j, _ := range screen.data[i] {
			screen.data[i][j] = ' '
		}
	}
}

func (screen *Screen) computeFrame(A, B, K1 float64) {

	// precompute sines and cosines of A and B
	cosA := math.Cos(A)
	sinA := math.Sin(A)
	cosB := math.Cos(B)
	sinB := math.Sin(B)

	screen.clear()
	zbuffer := newZBuffer(screen.dim)

	// theta goes around the cross-sectional circle of a torus
	for theta := 0.0; theta < 2.0*math.Pi; theta += theta_spacing {
		// precompute sines and cosines of theta
		costheta := math.Cos(theta)
		sintheta := math.Sin(theta)

		// phi goes around the center of revolution of a torus
		for phi := 0.0; phi < 2.0*math.Pi; phi += phi_spacing {
			// precompute sines and cosines of phi
			cosphi := math.Cos(phi)
			sinphi := math.Sin(phi)

			// the x,y coordinate of the circle, before revolving (factored out of the above equations)
			circlex := R2 + R1*costheta
			circley := R1 * sintheta

			// final 3D (x,y,z) coordinate after rotations, directly from our math above
			x := circlex*(cosB*cosphi+sinA*sinB*sinphi) - circley*cosA*sinB
			y := circlex*(sinB*cosphi-sinA*cosB*sinphi) + circley*cosA*cosB
			z := K2 + cosA*circlex*sinphi + circley*sinA
			ooz := 1 / z // "one over z"

			// x and y projection.  note that y is negated here, because y goes up in
			// 3D space but down on 2D displays.
			xp := int(float64(screen.dim)/2.0 + K1*ooz*x)
			yp := int(float64(screen.dim)/2.0 - K1*ooz*y)

			// calculate luminance.  ugly, but correct.
			L := cosphi*costheta*sinB - cosA*costheta*sinphi - sinA*sintheta +
				cosB*(cosA*sintheta-costheta*sinA*sinphi)
			// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface is
			// pointing away from us, so we won't bother trying to plot it.
			if L > 0 {
				// test against the z-buffer.  larger 1/z means the pixel is closer to
				// the viewer than what's already plotted.
				if ooz > (*zbuffer)[yp][xp] {
					(*zbuffer)[yp][xp] = ooz
					luminance_index := int(L * 8.0) // this brings L into the range 0..11 (8*sqrt(2) = 11.3)
					// now we lookup the character corresponding to the luminance and plot it in our output:
					screen.data[yp][xp] = ".,-~:;=!*#$@"[luminance_index]
				}
			}
		}
	}
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// adapted from:
// https://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/
func GetWinsize() (*winsize, error) {
	ws := new(winsize)

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)), // #nosec G103
	)

	if int(r1) == -1 {
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}

func animate(screen *Screen) {
	// Calculate K1 based on screen size: the maximum x-distance occurs roughly at
	// the edge of the torus, which is at x=R1+R2, z=0.  we want that to be
	// displaced 3/8ths of the width of the screen, which is 3/4th of the way from
	// the center to the side of the screen.
	// screen_width*3/8 = K1*(R1+R2)/(K2+0)
	// screen_width*K2*3/(8*(R1+R2)) = K1
	A, B, K1 := 1.0, 1.0, float64(screen.dim)*K2*3.0/(8.0*(R1+R2))

	c := time.Tick(frame_delay * time.Millisecond) // create timer channel
	for range c {
		A += 0.07
		B += 0.03
		screen.computeFrame(A, B, K1)
		screen.render()
	}
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return TickMsg{id: nextID()}
	})
}
