package ansi

import (
	"fmt"
)

type Attribute byte

const ESC = 27

// Base attributes
const (
	Bold Attribute = iota + 1
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

func Up(y int) string     { return fmt.Sprintf("\x1b[%dA", y) }
func Down(y int) string   { return fmt.Sprintf("\x1b[%dB", y) }
func Right(x int) string  { return fmt.Sprintf("\x1b[%dC", x) }
func Left(x int) string   { return fmt.Sprintf("\x1b[%dD", x) }
func Pos(x, y int) string { return fmt.Sprintf("\x1b[%d;%dH", y, x) }

const GOTO_TL = "\x1b[1;1H"
const CLEAR_RIGHT = "\x1b[0K"
const CLEAR_LEFT = "\x1b[1K"
const CLEAR_LINE = "\x1b[2K"
const CLEAR_SCREEN = "\x1b[2J"
const CLEAR_UP = "\x1b[0J"
const CLEAR_DOWN = "\x1b[1J"

func Set(vals ...Attribute) string {

	switch len(vals) {
	case 0:
		return "\x1b[0m"
	case 1:
		return fmt.Sprintf("\x1b[%dm", vals[0])
	case 2:
		return fmt.Sprintf("\x1b[%d;%dm", vals[0], vals[1])
	default:
		return fmt.Sprintf("\x1b[%d;%dm", vals[0], vals[1])
		// TODO :: Add support beyond 2
	}

}
