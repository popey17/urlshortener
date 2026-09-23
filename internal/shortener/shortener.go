package shortener

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"strings"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateCode() (string, error) {

	b := make([]byte, 7)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i, v := range b {
		b[i] = alphabet[v%byte(len(alphabet))]
	}

	return string(b), nil
}

func ValidateURL(originalUrl string) error {
	rawUrl := strings.TrimSpace(originalUrl)
	if rawUrl == "" {
		return fmt.Errorf("url is required")
	}

	u, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q", u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("invalid url:")
	}

	return nil
}
