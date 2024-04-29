package sshsession

import "golang.org/x/crypto/ssh"

func Run(client *ssh.Client, cmd string) ([]byte, error) {
	sess, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()
	return sess.CombinedOutput(cmd)
}
