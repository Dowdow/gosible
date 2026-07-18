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

// Args is implemented by every action type (copy, dir, docker, file, shell).
type Args interface {
	// Validate checks the action's own fields are structurally correct.
	// It must not depend on anything outside the action itself.
	Validate() error

	// Prepare resolves context-dependent values (relative paths, env(VAR)
	// placeholders) using the resolver functions provided by the config package.
	Prepare(resolvePath func(string) string, replaceEnv func(string) string) error

	// Run executes the action against executor and returns its output.
	Run(executor Executor) (stdout string, stderr string, err error)
}
