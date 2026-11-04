package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// secretValue is a keychain secret as fnox resolves it for this repo (fnox.toml).
func secretValue(key string) (string, error) {
	out, err := exec.Command("fnox", "get", key).Output()
	if err != nil {
		return "", fmt.Errorf("fnox get %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// storeSecret puts value in the keychain under fnox.toml's key (-p keychain: never inlined into fnox.toml).
func storeSecret(key, value string) error {
	cmd := exec.Command("fnox", "set", "-p", "keychain", "-k", "UBUNTU_WEBSITE_"+key, key)
	cmd.Stdin = strings.NewReader(value)
	return cmd.Run()
}

func openURL(url string) { _ = browser.OpenURL(url) }
