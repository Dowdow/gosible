package runner

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Machine Machine
	Actions []Action
}

type Machine struct {
	Name     string
	Address  string
	User     string
	Password string
	Ssh      string
	Become   string
}

type Action struct {
	Name string
	Type string
	Args Args
}

type Args interface {
	Pty() bool
	Run(session *ssh.Session, ch chan tea.Msg) error
}

type ActionStartMsg struct {
	Name string
}

type ActionEndMsg struct {
	Success bool
}

type EndMsg struct{}

type StdoutMsg struct {
	Msg string
}

type StderrMsg struct {
	Msg string
}

type ErrorMsg struct {
	Error error
}

type Runner struct {
	config Config
}

func NewRunner(config *Config) *Runner {
	return &Runner{
		config: *config,
	}
}

func (r *Runner) Run(ch chan tea.Msg) {
	authMethods := []ssh.AuthMethod{}

	// Private key auth method
	if r.config.Machine.Ssh != "" {
		key, err := os.ReadFile(r.config.Machine.Ssh)
		if err != nil {
			ch <- ErrorMsg{Error: fmt.Errorf("cannot read private key : %v", err)}
			return
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			ch <- ErrorMsg{Error: fmt.Errorf("cannot parse private key : %v", err)}
			return
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// Password auth method
	if r.config.Machine.Password != "" {
		authMethods = append(authMethods, ssh.Password(r.config.Machine.Password))
	}

	clientConfig := &ssh.ClientConfig{
		User:            r.config.Machine.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // dev
	}

	client, err := ssh.Dial("tcp", r.config.Machine.Address, clientConfig)
	if err != nil {
		ch <- ErrorMsg{Error: fmt.Errorf("cannot ssh dial : %v", err)}
		return
	}
	defer client.Close()

	writerStdout := newTeaWriterStdout(ch)
	writerStderr := newTeaWriterStderr(ch)

	for _, action := range r.config.Actions {
		session, err := client.NewSession()
		if err != nil {
			ch <- ErrorMsg{Error: fmt.Errorf("cannot create ssh session : %v", err)}
			return
		}
		defer session.Close()

		ch <- ActionStartMsg{Name: action.Name}

		/*		if action.Args.Pty() {
				modes := ssh.TerminalModes{
					ssh.ECHO:          0,
					ssh.ICANON:        0,
					ssh.ISIG:          0,
					ssh.TTY_OP_ISPEED: 14400,
					ssh.TTY_OP_OSPEED: 14400,
				}

				fd := int(os.Stdin.Fd())
				width, height, _ := term.GetSize(fd)

				if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
					ch <- ActionEndMsg{Success: false}
					ch <- ErrorMsg{Error: fmt.Errorf("cannot request pty : %v", err)}
					return
				}
				} */

		session.Stdin = nil
		session.Stdout = writerStdout
		session.Stderr = writerStderr

		err = action.Args.Run(session, ch)
		if err != nil {
			ch <- ActionEndMsg{Success: false}
			ch <- ErrorMsg{Error: err}
			return
		}

		ch <- ActionEndMsg{Success: true}
	}

	ch <- EndMsg{}
}
