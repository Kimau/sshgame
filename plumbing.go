package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

// parseDims extracts terminal dimensions (width x height) from the provided buffer.
func parseDims(b []byte) (int, int) {
	w := int(binary.BigEndian.Uint32(b))
	h := int(binary.BigEndian.Uint32(b[4:]))
	return w, h
}

func loginAuthFunc(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	// if c.User() == "testuser" && string(pass) == "tiger" { return nil, nil}
	// return nil, fmt.Errorf("password rejected for %q", c.User())

	return nil, nil
}

type GameChan struct {
	netChan ssh.Channel
	req     <-chan *ssh.Request
}

func starServer(hostaddr string) <-chan GameChan {

	gcChan := make(chan GameChan)

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{PasswordCallback: loginAuthFunc}

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
	listener, err := net.Listen("tcp", hostaddr)
	if err != nil {
		panic("failed to listen for connection")
	}

	fmt.Println("STARTING SERVER")

	////////////////////////////////////////
	/// LISTEN LOOP
	go func() {
		for {

			nConn, err := listener.Accept()
			if err != nil {
				panic("failed to accept incoming connection")
			}

			fmt.Println("Accepted Connection")

			// Before use, a handshake must be performed on the incoming
			// net.Conn.
			_, newConnChan, reqs, err := ssh.NewServerConn(nConn, config)
			if err != nil {
				panic("failed to handshake")
			}
			// The incoming Request channel must be serviced.
			go ssh.DiscardRequests(reqs)

			nChan := <-newConnChan

			// Channels have a type, depending on the application level
			// protocol intended. In the case of a shell, the type is
			// "session" and ServerShell may be used to present a simple
			// terminal interface.
			if nChan.ChannelType() != "session" {
				nChan.Reject(ssh.UnknownChannelType, "unknown channel type")
				return
			}
			channel, requests, err := nChan.Accept()
			if err != nil {
				panic("could not accept channel.")
			}

			fmt.Println("Channel Connection")

			newGCobj := GameChan{
				netChan: channel,
				req:     requests,
			}

			gcChan <- newGCobj
		}
	}()
	//////////////////////////////////////////

	return gcChan

}
