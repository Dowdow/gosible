package args

import "golang.org/x/crypto/ssh"

type Args interface {
	Pty() bool
	Run(session *ssh.Session) error
}
