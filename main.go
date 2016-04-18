package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"./ansi"
)

type GamePlayer struct {
	login         string
	pass          string
	charactername string
	team          int
	currChand     *GameChan
}

var playerMap map[string]*GamePlayer

func main() {
	data, _ := ioutil.ReadFile("gpolice.ans")
	str := ansi.AnsFileToStr(data)
	fmt.Println(str)
	fmt.Println("----------------------------------------------------------------")
	trimStr, ans := ansi.AnsFileTrim(str, 20, 9)
	fmt.Println("12345678901234567890")
	fmt.Println(trimStr)
	fmt.Println(ans)

	playerMap = make(map[string]*GamePlayer)
	newConnChan := starServer("0.0.0.0:2022")

	for {
		select {
		case c := <-newConnChan:
			go startClient(c)
		}
	}
}

func startClient(gc GameChan) {

	gc.chanWidth, gc.chanHeight = 80, 100

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
				gc.chanWidth, gc.chanHeight = parseDims(req.Payload[termLen+4:])
				req.Reply(true, nil)

			case "window-change":
				gc.chanWidth, gc.chanHeight = parseDims(req.Payload)
				req.Reply(true, nil)

			default:
				req.Reply(false, nil)
			}

		}
	}()

	term := terminal.NewTerminal(gc.netChan, "> ")

	PrintWelcome(term)
	fmt.Fprint(term, ansi.Set(ansi.FgBlack, ansi.BgYellow)+ansi.CLEAR_LINE+"LOGIN"+ansi.Set()+"\n\r Warning :: Pease do not reuse a password I'm sorting these plaintext this is a jam game. Your SSH connection is secure but I'm not salting and hashing your pass (yet) \n\r")

	fmt.Fprint(term, "Character Login:")
	login, _ := term.ReadLine()
	p, ok := playerMap[login]
	if !ok {
		// New User
		MakeNewUser(term, login)
		p = playerMap[login]
		p.currChand = &gc
	} else {
		p.currChand = &gc
		LoginUser(term, p)
	}

	go GameLoopForPlayer(term, p)
}

func GameLoopForPlayer(term *terminal.Terminal, p *GamePlayer) {
	defer p.currChand.netChan.Close()
	for {
		line, err := term.ReadLine()
		if err != nil {
			break
		}

		words := strings.Split(line, " ")

		switch words[0] {
		case "welcome":
			PrintWelcome(term)

		case "police":
			data, _ := ioutil.ReadFile("gpolice.ans")
			str := ansi.AnsFileToStr(data)
			fmt.Fprint(term, str)

		case "test":
			data, _ := ioutil.ReadFile("gpolice.ans")
			str := ansi.AnsFileToStr(data)

			x, _ := strconv.Atoi(words[1])
			y, _ := strconv.Atoi(words[2])

			_, str = ansi.AnsFileTrim(str, x, y)
			fmt.Fprint(term, str)

		case "br":
			fmt.Fprintf(term, "%s%sX", ansi.CLEAR_SCREEN, ansi.Pos(p.currChand.chanWidth, p.currChand.chanHeight))

		case "clear":
			fmt.Fprintf(term, "%s%s", ansi.CLEAR_SCREEN, ansi.GOTO_TL)

		case "border":
			fmt.Fprint(term, ansi.CLEAR_SCREEN+ansi.Set(ansi.FgBlack, ansi.BgWhite)+ansi.GOTO_TL+ansi.CLEAR_LINE+ansi.Pos(p.currChand.chanWidth, p.currChand.chanHeight-1)+ansi.CLEAR_LINE)

			for y := 1; y < p.currChand.chanHeight; y += 1 {
				fmt.Fprintf(term, ansi.Pos(p.currChand.chanWidth, y)+"  ")
			}

			fmt.Fprint(term, ansi.Set()+ansi.Pos(1, p.currChand.chanHeight))

		default:
			fmt.Fprintf(term, "Uknown command: [%s]\n\r", line)
		}

	}
}

func LoginUser(term *terminal.Terminal, p *GamePlayer) {
	// Check Password
	for i := 0; i < 5; i += 1 {
		pass, _ := term.ReadPassword("Password:")
		if p.pass == pass {
			return
		}

		fmt.Fprint(term, ansi.Set(ansi.FgRed)+"*** Invalid PASS ***"+ansi.Set()+"\n\r")
	}

	fmt.Fprint(term, "CONNECTION DENIED")
	p.currChand.netChan.Close()
}

func MakeNewUser(term *terminal.Terminal, l string) {

	fmt.Fprint(term, ansi.Set(ansi.FgBlue, ansi.Underline)+"  Creating New Character  "+ansi.Set()+"\n\r")

	p := GamePlayer{login: l}
	playerMap[l] = &p

	// Set Password
	for p.pass = ""; p.pass == ""; {
		pass1, _ := term.ReadPassword("Password:")
		pass2, _ := term.ReadPassword("Password Again:")

		if pass1 != pass2 {
			fmt.Fprint(term, "**** Passwords don't match ***** \n\r")
		} else {
			p.pass = pass1
		}
	}

	//
}

func PrintWelcome(term *terminal.Terminal) {
	fmt.Fprint(term, ansi.CLEAR_SCREEN+ansi.GOTO_TL)

	startL := ansi.CHOME + ansi.Set(ansi.FgRed, ansi.BgBlack) + ansi.CLEAR_LINE + "█" + ansi.CEND + "█ " + ansi.CUP

	fmt.Fprint(term, ansi.Set(ansi.FgRed, ansi.BgRed)+ansi.CLEAR_LINE+ansi.CDOWN+
		startL+"\n"+
		startL+" "+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"  ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"     █    ██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████▄  █    ██  ███▄ ▄███"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓   ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████▄  ▄▄▄       ██▀███  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████  \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"     ██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀ ██▌ ██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"▀█▀ ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀ ██▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████▄    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+" ██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█   ▀ \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░    ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██   █▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░   ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██   █▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██  ▀█▄  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"▄█ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"███ \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░    ▓▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█▄   ▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██    "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█▄   ▌"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▄▄▄▄██ "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██▀▀█▄  "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█  ▄ \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒   ░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓  ▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█   "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒░"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▓ ▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"██"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒░▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"████"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"▒\n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░ ▒░▓  ░░▒▓▒ ▒ ▒  ▒▒▓  ▒ ░▒▓▒ ▒ ▒ ░ ▒░   ░  ░    ▒▒▓  ▒  ▒▒   ▓▒"+ansi.Set(ansi.FgHiRed, ansi.BgBlack)+"█"+ansi.Set(ansi.FgRed, ansi.BgBlack)+"░░ ▒▓ ░▒▓░░░ ▒░ ░\n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+" ░ ░ ▒  ░░░▒░ ░ ░  ░ ▒  ▒ ░░▒░ ░ ░ ░  ░      ░    ░ ▒  ▒   ▒   ▒▒ ░  ░▒ ░ ▒░ ░ ░  ░ \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"   ░ ░    ░░░ ░ ░  ░ ░  ░  ░░░ ░ ░ ░      ░       ░ ░  ░   ░   ▒     ░░   ░    ░    \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"     ░  ░   ░        ░       ░            ░         ░          ░  ░   ░        ░  ░ \n\r"+
		startL+" "+ansi.Set(ansi.FgRed, ansi.BgBlack)+"                                                                                    \n\r"+
		ansi.Set(ansi.FgRed, ansi.BgRed)+ansi.CLEAR_LINE+ansi.CDOWN+ansi.Set()+ansi.CHOME)
}
