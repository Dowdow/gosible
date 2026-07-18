// Command fakesshd is a throwaway SSH server used only to record demo.gif
// against something other than a real homelab. It accepts any password,
// runs no real commands, and always reports success after a short delay
// (so the TUI's spinner is visible in the recording).
//
// It is not part of the gosible binary and is never imported by it.
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

const listenAddr = "127.0.0.1:2244"

func main() {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("generating host key: %v", err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		log.Fatalf("wrapping host key: %v", err)
	}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil // accept anything, this is a demo fixture
		},
	}
	config.AddHostKey(signer)

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("listening on %s: %v", listenAddr, err)
	}
	defer l.Close()

	fmt.Printf("fakesshd listening on %s\n", listenAddr)

	for {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		go handleConn(conn, config)
	}
}

func handleConn(conn net.Conn, config *ssh.ServerConfig) {
	sconn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		return
	}
	defer sconn.Close()
	go ssh.DiscardRequests(reqs)

	for newChan := range chans {
		if newChan.ChannelType() != "session" {
			newChan.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}

		channel, requests, err := newChan.Accept()
		if err != nil {
			continue
		}

		go handleSession(channel, requests)
	}
}

func handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	for req := range requests {
		if req.Type != "exec" {
			if req.WantReply {
				req.Reply(false, nil)
			}
			continue
		}

		req.Reply(true, nil)

		// Give the pending spinner a moment to be visible in the recording.
		time.Sleep(600 * time.Millisecond)

		var status struct{ Status uint32 }
		channel.SendRequest("exit-status", false, ssh.Marshal(&status))
		return
	}
}
