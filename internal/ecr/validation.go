package ecr

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrUsernameInvalid = errors.New("invalid docker username")
	ErrRegistryInvalid = errors.New("invalid registry")
)

// allowedRunes in a docker username
// it is a best guess as I could not find any reference
func allowedRunes(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '.' || r == '_' || r == '-'
}

func every(s string, fn func(rune) bool) bool {
	for _, r := range s {
		if !fn(r) {
			return false
		}
	}
	return true
}

func IsCredentialValid(c Credential) error {
	if !every(c.Username, allowedRunes) {
		return fmt.Errorf("username `%s` is not valid: %w", c.Username, ErrUsernameInvalid)
	}

	if strings.TrimSpace(c.RegistryURL) == "" {
		return fmt.Errorf("empty registry url: %w", ErrRegistryInvalid)
	}

	if _, err := url.Parse(c.RegistryURL); err != nil {
		return fmt.Errorf("registry `%s` not valid: %w", c.RegistryURL, ErrRegistryInvalid)
	}

	return nil
}
