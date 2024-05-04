// SPDX-License-Identifier: MIT
// From https://gitlab.com/hetznercloud/fleeting-plugin-hetzner/-/blob/0f60204582289c243599f8ca0f5be4822789131d/internal/utils/ssh_key.go
// Copyright (c) 2024 Hetzner Cloud GmbH

package sshkey

import (
	"crypto/ed25519"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
)

func GenerateKeyPair() ([]byte, []byte, error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := encodePublicKey(pub)
	if err != nil {
		return nil, nil, err
	}

	privBytes, err := encodePrivateKey(priv)
	if err != nil {
		return nil, nil, err
	}

	return pubBytes, privBytes, nil
}

func encodePublicKey(pub ed25519.PublicKey) ([]byte, error) {
	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		return nil, err
	}

	return ssh.MarshalAuthorizedKey(sshPub), nil
}

func encodePrivateKey(priv ed25519.PrivateKey) ([]byte, error) {
	privPem, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(privPem), nil
}
