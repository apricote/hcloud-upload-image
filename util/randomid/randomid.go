package randomid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func Generate() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return hex.EncodeToString(b), nil
}
