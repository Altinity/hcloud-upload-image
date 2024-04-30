package randomid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// From https://gitlab.com/hetznercloud/fleeting-plugin-hetzner/-/blob/0f60204582289c243599f8ca0f5be4822789131d/internal/utils/random.go

func Generate() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}
	return hex.EncodeToString(b), nil
}
