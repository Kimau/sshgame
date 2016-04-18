package ansi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
func Hor(x int) string    { return fmt.Sprintf("\x1b[%dG", x) }

const CUP = "\x1b[A"
const CDOWN = "\x1b[B"
const CRIGHT = "\x1b[C"
const CLEFT = "\x1b[D"
const CHOME = "\x1b[999D"
const CEND = "\x1b[999C"

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

var IBMTable = [...]string{
	"\u00C7", "\u00FC", "\u00E9", "\u00E2", "\u00E4", "\u00E0", "\u00E5", "\u00E7", "\u00EA", "\u00EB", "\u00E8", "\u00EF", "\u00EE", "\u00EC", "\u00C4", "\u00C5", "\u00C9", "\u00E6", "\u00C6", "\u00F4", "\u00F6", "\u00F2", "\u00FB", "\u00F9", "\u00FF", "\u00D6", "\u00DC", "\u00A2", "\u00A3", "\u00A5", "\u20A7", "\u0192", "\u00E1", "\u00ED", "\u00F3", "\u00FA", "\u00F1", "\u00D1", "\u00AA", "\u00BA", "\u00BF", "\u2310", "\u00AC", "\u00BD", "\u00BC", "\u00A1", "\u00AB", "\u00BB", "\u2591", "\u2592", "\u2593", "\u2502", "\u2524", "\u2561", "\u2562", "\u2556", "\u2555", "\u2563", "\u2551", "\u2557", "\u255D", "\u255C", "\u255B", "\u2510", "\u2514", "\u2534", "\u252C", "\u251C", "\u2500", "\u253C", "\u255E", "\u255F", "\u255A", "\u2554", "\u2569", "\u2566", "\u2560", "\u2550", "\u256C", "\u2567", "\u2568", "\u2564", "\u2565", "\u2559", "\u2558", "\u2552", "\u2553", "\u256B", "\u256A", "\u2518", "\u250C", "\u2588", "\u2584", "\u258C", "\u2590", "\u2580", "\u03B1", "\u00DF", "\u0393", "\u03C0", "\u03A3", "\u03C3", "\u00B5", "\u03C4", "\u03A6", "\u0398", "\u03A9", "\u03B4", "\u221E", "\u03C6", "\u03B5", "\u2229", "\u2261", "\u00B1", "\u2265", "\u2264", "\u2320", "\u2321", "\u00F7", "\u2248", "\u00B0", "\u2219", "\u00B7", "\u221A", "\u207F", "\u00B2", "\u25A0", "\u00A0",
}

func IBMExtend(src byte) string {
	if src < 128 {
		return string(src)
	}

	return IBMTable[int(src)-128]

}

func AnsFileToStr(data []byte) string {

	bigStr := ""
	for _, v := range data {
		if v == 0x1A {
			break
		}
		bigStr += IBMExtend(v)
	}

	bigStr += Set()

	return bigStr
}

type attribStruct struct {
	fgCol      Attribute
	bgCol      Attribute
	other      []Attribute
	textOffset int
}

func AnsFileTrim(src string, xLimit int, yLimit int) string {

	//skipForward := regexp.MustCompile("\xb1\\[([0-9]*)C")
	//attrSet := regexp.MustCompile("\xb1\\[([0-9]*;*[0-9]*)m")
	//
	allescape := regexp.MustCompile("\\[([0-9\\;]*)[A-Za-z]")
	space := regexp.MustCompile("\\[([0-9]*)C")
	attrib := regexp.MustCompile("\\[([0-9\\;]*)m")
	notAttrib := regexp.MustCompile("\\[([0-9\\;]*)[^m]")
	numBits := regexp.MustCompile("[0-9]+")
	escBits := regexp.MustCompile("")

	// Remove Cursor Movement
	noCurString := space.ReplaceAllStringFunc(src, func(x string) string {
		f := space.FindStringSubmatch(x)
		res := ""
		for i, _ := strconv.Atoi(f[1]); i > 0; i -= 1 {
			res += " "
		}
		return res
	})

	noCurString = notAttrib.ReplaceAllString(noCurString, "")

	// Split into Lines and Trim
	lines := strings.Split(noCurString, "\n")
	if yLimit <= 0 {
		yLimit = len(lines)
	}
	lines = lines[0:yLimit]

	// Attrin
	attribListIdx := make(map[int][]attribStruct)
	var prevAttrib *attribStruct
	prevAttrib = nil

	for y, ln := range lines {

		offset := 0
		attribListIdx[y] = []attribStruct{}
		if y > 0 && prevAttrib != nil {
			offset = 1
			attribListIdx[y] = []attribStruct{*prevAttrib}
		}

		// Atrrib
		textString := attrib.ReplaceAllStringFunc(ln, func(x string) string {
			f := numBits.FindAllStringSubmatch(x, -1)
			attrb := attribStruct{}
			if prevAttrib != nil {
				attrb = *prevAttrib
			}

			for _, a := range f {
				intVal, e := strconv.Atoi(a[0])

				if e != nil {
					fmt.Println(e)
				} else {
					attVal := Attribute(intVal)
					switch attVal {
					case 0:
						attrb = attribStruct{}

					case Bold, Faint, Italic, Underline, BlinkSlow, BlinkRapid, ReverseVideo, Concealed, CrossedOut:
						attrb.other = append(attrb.other, attVal)
					case FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite, FgHiBlack, FgHiRed, FgHiGreen, FgHiYellow, FgHiBlue, FgHiMagenta, FgHiCyan, FgHiWhite:
						attrb.fgCol = attVal
					case BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite, BgHiBlack, BgHiRed, BgHiGreen, BgHiYellow, BgHiBlue, BgHiMagenta, BgHiCyan, BgHiWhite:
						attrb.bgCol = attVal
					}
				}
			}

			prevAttrib = &attrb
			attribListIdx[y] = append(attribListIdx[y], attrb)
			return ""
		})

		//
		indexBits := escBits.FindAllStringIndex(textString, -1)
		for i, v := range indexBits {
			attribListIdx[y][i+offset].textOffset = v - i
		}
		lines[y] = escBits.ReplaceAllString(textString, "")

	}

	textString := strings.Join(lines, Left(-xLimit)+Down(1))

	// Remove All other bits
	return textString
}
