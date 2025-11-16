package runner

import (
	"bytes"
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
	Validate() error
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

	for _, action := range r.config.Actions {
		session, err := client.NewSession()
		if err != nil {
			ch <- ErrorMsg{Error: fmt.Errorf("cannot create ssh session : %v", err)}
			return
		}
		defer session.Close()

		ch <- ActionStartMsg{Name: action.Name}

		var stdoutBuf, stderrBuf bytes.Buffer
		session.Stdout = &stdoutBuf
		session.Stderr = &stderrBuf

		err = action.Args.Run(session, ch)
		if err != nil {
			ch <- StdoutMsg{Msg: stdoutBuf.String()}
			ch <- StderrMsg{Msg: stderrBuf.String()}
			ch <- ActionEndMsg{Success: false}
			ch <- ErrorMsg{Error: err}
			return
		}

		ch <- StdoutMsg{Msg: stdoutBuf.String()}
		ch <- StderrMsg{Msg: stderrBuf.String()}
		ch <- ActionEndMsg{Success: true}
	}

	ch <- EndMsg{}
}
