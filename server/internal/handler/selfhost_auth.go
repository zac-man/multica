package handler

import (
	"crypto/subtle"
	"os"
	"strings"
)

// selfhostPasswordLogin returns the configured single-user email and password when
// SELFHOST_LOGIN_EMAIL and SELFHOST_LOGIN_PASSWORD are both set (non-empty after trim).
func selfhostPasswordLogin() (email string, password string, ok bool) {
	email = strings.ToLower(strings.TrimSpace(os.Getenv("SELFHOST_LOGIN_EMAIL")))
	password = strings.TrimSpace(os.Getenv("SELFHOST_LOGIN_PASSWORD"))
	if email == "" || password == "" {
		return "", "", false
	}
	return email, password, true
}

func selfhostPasswordEqual(input, expected string) bool {
	if len(input) != len(expected) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(input), []byte(expected)) == 1
}
