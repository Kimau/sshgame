package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"./ansi"
)

func main() {
	testSSHServer()
}

// parseDims extracts terminal dimensions (width x height) from the provided buffer.
func parseDims(b []byte) (int, int) {
	w := int(binary.BigEndian.Uint32(b))
	h := int(binary.BigEndian.Uint32(b[4:]))
	return w, h
}

func testSSHServer() {
	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			if c.User() == "testuser" && string(pass) == "tiger" {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic("Failed to load private key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "0.0.0.0:2022")
	if err != nil {
		panic("failed to listen for connection")
	}
	nConn, err := listener.Accept()
	if err != nil {
		panic("failed to accept incoming connection")
	}

	// Before use, a handshake must be performed on the incoming
	// net.Conn.
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		panic("failed to handshake")
	}
	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level
		// protocol intended. In the case of a shell, the type is
		// "session" and ServerShell may be used to present a simple
		// terminal interface.
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			panic("could not accept channel.")
		}

		chanWidth, chanHeight := 80, 100

		// Sessions have out-of-band requests such as "shell",
		// "pty-req" and "env".  Here we handle only the
		// "shell" request.
		go func(in <-chan *ssh.Request) {
			for req := range in {
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
		}(requests)

		term := terminal.NewTerminal(channel, "> ")

		fmt.Fprintf(term, "%s %s [%d,%d] %s%s                  Login                 %s \n\r", ansi.CLEAR_SCREEN, ansi.Pos(chanWidth-10, chanHeight), chanWidth, chanHeight, ansi.GOTO_TL, ansi.Set(ansi.FgBlack, ansi.BgYellow), ansi.Set())

		go func() {
			defer channel.Close()
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
					fmt.Fprintf(term, "[")
					for x := 3; x < chanWidth; x += 1 {
						fmt.Fprintf(term, "=")
					}
					fmt.Fprintf(term, "]\n\r")
					for y := 3; y < chanHeight; y += 1 {
						fmt.Fprintf(term, "[")
						for x := 3; x < chanWidth; x += 1 {
							fmt.Fprintf(term, ".")
						}
						fmt.Fprintf(term, "]\n\r")
					}

					fmt.Fprintf(term, "[")
					for x := 3; x < chanWidth; x += 1 {
						fmt.Fprintf(term, "=")
					}
					fmt.Fprintf(term, "]\n\r")

				default:
					fmt.Fprintf(term, "Uknown command: [%s]\n\r", line)
				}

			}
		}()
	}
}
