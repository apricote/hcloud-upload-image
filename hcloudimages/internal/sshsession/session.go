package sshsession

import (
	"io"

	"golang.org/x/crypto/ssh"
)

func Run(client *ssh.Client, cmd string, stdin io.Reader) ([]byte, error) {
	sess, err := client.NewSession()

	if err != nil {
		return nil, err
	}
	defer func() { _ = sess.Close() }()

	if stdin != nil {
		sess.Stdin = stdin
	}
	return sess.CombinedOutput(cmd)
}
