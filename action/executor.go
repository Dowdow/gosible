package action

import "io"

// Executor abstracts what an action needs to run a command on a remote
// machine, without depending on a concrete SSH implementation.
type Executor interface {
	// Run executes cmd and waits for it to complete, returning its captured output.
	Run(cmd string) (stdout string, stderr string, err error)

	// Start begins executing cmd and returns a pipe to its stdin. The caller
	// must close stdin and call wait to obtain the captured output and let
	// the remote command complete.
	Start(cmd string) (stdin io.WriteCloser, wait func() (stdout string, stderr string, err error), err error)
}
