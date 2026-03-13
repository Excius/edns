package normalize

import (
	"errors"
	"strings"

	"golang.org/x/net/idna"
)

func NormalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("invalid email format")
	}

	loca := parts[0]
	domain := parts[1]

	asciiDomain, err := idna.Lookup.ToASCII(domain)
	if err != nil {
		return "", errors.New("invalid email domain")
	}

	return loca + "@" + asciiDomain, nil
}
