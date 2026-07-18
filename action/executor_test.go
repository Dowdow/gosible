package action

import (
	"bytes"
	"io"
)

// fakeExecutor is an in-memory action.Executor test double: no SSH, no
// network, no external process.
type fakeExecutor struct {
	runs   []string
	starts []string
	stdin  *bytes.Buffer

	runOut, runErrOut string
	runErr            error

	startErr error
	waitErr  error
}

func (e *fakeExecutor) Run(cmd string) (string, string, error) {
	e.runs = append(e.runs, cmd)
	return e.runOut, e.runErrOut, e.runErr
}

func (e *fakeExecutor) Start(cmd string) (io.WriteCloser, func() (string, string, error), error) {
	e.starts = append(e.starts, cmd)
	if e.startErr != nil {
		return nil, nil, e.startErr
	}

	e.stdin = &bytes.Buffer{}
	wait := func() (string, string, error) {
		return e.runOut, e.runErrOut, e.waitErr
	}

	return nopWriteCloser{e.stdin}, wait, nil
}

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

func identityResolvePath(p string) string { return p }
func identityReplaceEnv(s string) string  { return s }
