package ansi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

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

func RemoveCursorMovement(src string) string {
	space := regexp.MustCompile("\\[([0-9]*)C")

	// Remove Cursor Movement
	noForward := space.ReplaceAllStringFunc(src, func(x string) string {
		f := space.FindStringSubmatch(x)
		res := ""
		for i, _ := strconv.Atoi(f[1]); i > 0; i -= 1 {
			res += " "
		}
		return res
	})

	// TODO :: Handle other cursor moves ABDH

	return noForward
}

func AnsFileTrimHeight(src string, yLimit int) (txtRes string, ansRes string) {
	noCurString := RemoveCursorMovement(src)

	// Split into Lines and Trim
	lines := strings.Split(noCurString, "\n")
	if yLimit > 0 {
		lines = lines[0:yLimit]
	}

	allescape := regexp.MustCompile("\\[([0-9\\;]*)[^\\;0-9]")
	txtRes = allescape.ReplaceAllString(strings.Join(lines, "\n\r"), "")
	ansRes = txtRes

	// Remove All other bits
	return txtRes, ansRes
}

func StripANSI(src string) string {
	allescape := regexp.MustCompile("\\[([0-9\\;]*)[^\\;0-9]")
	return allescape.ReplaceAllString(src, "")
}

func AnsFileTrim(src string, xLimit int, yLimit int) string {
	return AnsFileBoxTrim(src, 0, xLimit, 0, yLimit)
}

func AnsFileBoxTrim(src string, xMin int, xMax int, yMin int, yMax int) string {
	attrib := regexp.MustCompile("\\[([0-9\\;]*)m")
	numBits := regexp.MustCompile("[0-9]+")

	// Remove Cursor Movement
	noCurString := RemoveCursorMovement(src)

	// Split into Lines and Trim
	lines := strings.Split(noCurString, "\n")

	// Bounds Check
	if xMin < 0 {
		xMin = 0
	}
	if yMin < 0 {
		yMin = 0
	}
	if xMax <= xMin {
		xMax = xMin + 1
	}

	if yMax > len(lines) {
		yMax = len(lines)
	}
	if yMax <= yMin {
		yMax = yMin + 1
	}

	lines = lines[yMin:yMax]

	// Per Line
	for y, ln := range lines {

		// Trim Characters from Line
		rArr := []rune(ln)
		filterArr := []rune{}
		inEscape := false
		counter := 0
		for _, c := range rArr {
			if inEscape {
				filterArr = append(filterArr, c)
				if unicode.IsLetter(c) {
					inEscape = false
				}
			} else if c == 0x1b {
				filterArr = append(filterArr, c)
				inEscape = true
			} else {
				counter += 1
				if counter >= xMin && counter < xMax {
					filterArr = append(filterArr, c)
				}
			}
		}
		textString := string(filterArr)

		// Atrrib
		fg := FgDefault
		bg := BgDefault
		textString = attrib.ReplaceAllStringFunc(textString, func(x string) string {
			f := numBits.FindAllStringSubmatch(x, -1)
			var atList AttributeList

			for _, a := range f {
				intVal, e := strconv.Atoi(a[0])

				if e != nil {
					fmt.Println(e)
				} else {
					attVal := Attribute(intVal)
					atList = append(atList, attVal)
				}
			}

			fg, bg = atList.ColourConsildate(fg, bg)
			return atList.ANSI()
		})

		lines[y] = textString
	}

	// Setup Screen
	ansRes := CSAVE + CurScrollUp(yMax-yMin) + CLOAD
	ansRes += strings.Join(lines, CDOWN+CurLeft(xMax-xMin)) + Set()
	return ansRes
}
