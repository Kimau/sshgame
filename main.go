package main

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"

	"./ansi"
)

func main() {
	newConnChan := starServer("0.0.0.0:2022")

	for {
		select {
		case c := <-newConnChan:
			go startClient(c)
		}
	}
}

func startClient(gc GameChan) {
	chanWidth, chanHeight := 80, 100

	// Sessions have out-of-band requests such as "shell",
	// "pty-req" and "env".  Here we handle only the
	// "shell" request.
	go func() {
		for req := range gc.req {
			fmt.Println("Request:", req)

			switch req.Type {
			case "shell":
				if len(req.Payload) > 0 {
					// We don't accept any commands, only the default shell.
					req.Reply(false, nil)
					continue
				}
				req.Reply(true, nil)

			case "pty-req":
				termLen := req.Payload[3]
				chanWidth, chanHeight = parseDims(req.Payload[termLen+4:])
				req.Reply(true, nil)

			case "window-change":
				chanWidth, chanHeight = parseDims(req.Payload)
				req.Reply(true, nil)

			default:
				req.Reply(false, nil)
			}

		}
	}()

	term := terminal.NewTerminal(gc.netChan, "> ")

	fmt.Fprintf(term, "%s %s [%d,%d] %s%s%s                  Login                 %s \n\r", ansi.CLEAR_SCREEN, ansi.Pos(chanWidth-10, chanHeight), chanWidth, chanHeight, ansi.GOTO_TL, ansi.CLEAR_LINE, ansi.Set(ansi.FgBlack, ansi.BgYellow), ansi.Set())
	fmt.Fprint(term, ansi.Set(ansi.FgRed, ansi.BgRed)+"                                                                                      \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"                                                                                    "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"  ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"     █    ██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████▄  █    ██  ███▄ ▄███"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓   ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████▄  ▄▄▄       ██▀███  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████  "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"     ██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀ ██▌ ██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"▀█▀ ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀ ██▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████▄    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+" ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█   ▀  "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░    ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██   █▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░   ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██   █▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  ▀█▄  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"▄█ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"███    "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░    ▓▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█▄   ▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█▄   ▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▄▄▄▄██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀▀█▄  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  ▄  "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓  ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█   "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒ "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░ ▒░▓  ░░▒▓▒ ▒ ▒  ▒▒▓  ▒ ░▒▓▒ ▒ ▒ ░ ▒░   ░  ░    ▒▒▓  ▒  ▒▒   ▓▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░ ▒▓ ░▒▓░░░ ▒░ ░ "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░ ░ ▒  ░░░▒░ ░ ░  ░ ▒  ▒ ░░▒░ ░ ░ ░  ░      ░    ░ ▒  ▒   ▒   ▒▒ ░  ░▒ ░ ▒░ ░ ░  ░ "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"   ░ ░    ░░░ ░ ░  ░ ░  ░  ░░░ ░ ░ ░      ░       ░ ░  ░   ░   ▒     ░░   ░    ░    "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"     ░  ░   ░        ░       ░            ░         ░          ░  ░   ░        ░  ░ "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"                                                                                    "+ansi.Set(ansi.BgRed, ansi.BgRed)+" \n\r"+
		"                                                                                      \n\r"+ansi.Set())

	go func() {
		defer gc.netChan.Close()
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}

			switch line {
			case "br":
				fmt.Fprintf(term, "%s%sX", ansi.CLEAR_SCREEN, ansi.Pos(chanWidth, chanHeight))

			case "clear":
				fmt.Fprintf(term, "%s%s", ansi.CLEAR_SCREEN, ansi.GOTO_TL)

			case "border":
				fmt.Fprint(term, ansi.CLEAR_SCREEN+ansi.Set(ansi.FgBlack, ansi.BgWhite)+ansi.GOTO_TL+ansi.CLEAR_LINE+ansi.Pos(chanWidth, chanHeight-1)+ansi.CLEAR_LINE)

				for y := 1; y < chanHeight; y += 1 {
					fmt.Fprintf(term, ansi.Pos(chanWidth, y)+"  ")
				}

				fmt.Fprint(term, ansi.Set()+ansi.Pos(1, chanHeight))

			default:
				fmt.Fprintf(term, "Uknown command: [%s]\n\r", line)
			}

		}
	}()

}
