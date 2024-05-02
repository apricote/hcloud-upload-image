package sshkey

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSSHKeyPair(t *testing.T) {
	pubBytes, privBytes, err := GenerateKeyPair()
	assert.Nil(t, err)

	pub := string(pubBytes)
	priv := string(privBytes)

	if !(strings.HasPrefix(priv, "-----BEGIN OPENSSH PRIVATE KEY-----\n") &&
		strings.HasSuffix(priv, "-----END OPENSSH PRIVATE KEY-----\n")) {
		assert.Fail(t, "private key is invalid", priv)
	}

	if !strings.HasPrefix(pub, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA") {
		assert.Fail(t, "public key is invalid", pub)
	}
}
