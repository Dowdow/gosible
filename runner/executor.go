package runner

import (
	"bytes"
	"io"

	"golang.org/x/crypto/ssh"
)

type sshExecutor struct {
	client *ssh.Client
}

func (e *sshExecutor) Run(cmd string) (string, string, error) {
	session, err := e.client.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(cmd)
	return stdoutBuf.String(), stderrBuf.String(), err
}

func (e *sshExecutor) Start(cmd string) (io.WriteCloser, func() (string, string, error), error) {
	session, err := e.client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		return nil, nil, err
	}

	if err := session.Start(cmd); err != nil {
		session.Close()
		return nil, nil, err
	}

	wait := func() (string, string, error) {
		defer session.Close()
		err := session.Wait()
		return stdoutBuf.String(), stderrBuf.String(), err
	}

	return stdin, wait, nil
}
