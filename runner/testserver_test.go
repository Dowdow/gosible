package runner_test

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// testSSHServer is a minimal in-process SSH server, written in pure Go
// (no system sshd, no ssh-keygen), used to drive runner.Runner end-to-end.
// It understands plain exec commands and the "scp -tr <dir>" sink protocol
// used by action.CopyArgs, writing anything "uploaded" under root.
type testSSHServer struct {
	listener net.Listener
	config   *ssh.ServerConfig
	root     string
	password string
}

// newTestSSHServer starts the server and registers its shutdown on cleanup.
// When password is empty, any password or public key is accepted; auth
// itself is never skipped, so every test genuinely exercises the client's
// auth method construction in runner.Runner.Run.
func newTestSSHServer(t *testing.T, password string) *testSSHServer {
	t.Helper()

	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating host key: %v", err)
	}
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatalf("wrapping host key: %v", err)
	}

	srv := &testSSHServer{root: t.TempDir(), password: password}

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if srv.password != "" && string(password) != srv.password {
				return nil, fmt.Errorf("wrong password")
			}
			return nil, nil
		},
		PublicKeyCallback: func(c ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		},
	}
	config.AddHostKey(signer)
	srv.config = config

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listening: %v", err)
	}
	srv.listener = l
	t.Cleanup(func() { l.Close() })

	go srv.serve()

	return srv
}

func (s *testSSHServer) addr() string {
	return s.listener.Addr().String()
}

func (s *testSSHServer) serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}

func (s *testSSHServer) handleConn(conn net.Conn) {
	sconn, chans, reqs, err := ssh.NewServerConn(conn, s.config)
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

		go s.handleSession(channel, requests)
	}
}

func (s *testSSHServer) handleSession(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	for req := range requests {
		if req.Type != "exec" {
			if req.WantReply {
				req.Reply(false, nil)
			}
			continue
		}

		var payload struct{ Command string }
		ssh.Unmarshal(req.Payload, &payload)
		req.Reply(true, nil)

		exitStatus := s.runCommand(channel, payload.Command)

		var status struct{ Status uint32 }
		status.Status = exitStatus
		channel.SendRequest("exit-status", false, ssh.Marshal(&status))
		return
	}
}

func (s *testSSHServer) runCommand(channel ssh.Channel, cmd string) uint32 {
	if strings.HasPrefix(cmd, "scp -tr ") {
		return s.scpSink(channel, strings.TrimPrefix(cmd, "scp -tr "))
	}

	if cmd == "false-command" {
		fmt.Fprint(channel.Stderr(), "boom")
		return 1
	}

	fmt.Fprintf(channel, "ran: %s", cmd)
	return 0
}

// scpSink implements just enough of the scp sink protocol to receive what
// action.CopyArgs writes: C (file), D (start dir) and E (end dir) lines.
func (s *testSSHServer) scpSink(channel ssh.Channel, destDir string) uint32 {
	reader := bufio.NewReader(channel)

	// destDir ("scp -tr <destDir>") is where the client wants things placed,
	// e.g. path.Dir(a.Dest). Map it under the sandbox root as-is: for a
	// directory copy the stream itself carries a leading D line naming the
	// final directory, so destDir must not be collapsed to its basename.
	root := filepath.Join(s.root, strings.TrimPrefix(destDir, "/"))
	if err := os.MkdirAll(root, 0o755); err != nil {
		return 1
	}
	dirStack := []string{root}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && line == "" {
				return 0
			}
			if err != io.EOF {
				return 1
			}
		}
		line = strings.TrimRight(line, "\n")
		if line == "" {
			return 0
		}

		switch line[0] {
		case 'E':
			if len(dirStack) > 1 {
				dirStack = dirStack[:len(dirStack)-1]
			}
		case 'D':
			name, _, err := parseSCPHeader(line)
			if err != nil {
				return 1
			}
			dir := filepath.Join(dirStack[len(dirStack)-1], name)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return 1
			}
			dirStack = append(dirStack, dir)
		case 'C':
			name, size, err := parseSCPHeader(line)
			if err != nil {
				return 1
			}

			buf := make([]byte, int(size))
			if _, err := io.ReadFull(reader, buf); err != nil {
				return 1
			}
			if _, err := reader.ReadByte(); err != nil { // trailing zero byte
				return 1
			}

			path := filepath.Join(dirStack[len(dirStack)-1], name)
			if err := os.WriteFile(path, buf, 0o644); err != nil {
				return 1
			}
		default:
			return 1
		}
	}
}

func parseSCPHeader(line string) (name string, size int64, err error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return "", 0, fmt.Errorf("malformed scp header: %q", line)
	}
	size, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, err
	}
	return parts[2], size, nil
}
