package main

import (
	"os/exec"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// storeSecret puts value in the keychain under fnox.toml's key (-p keychain: never inlined into fnox.toml).
func storeSecret(key, value string) error {
	cmd := exec.Command("fnox", "set", "-p", "keychain", "-k", "UBUNTU_WEBSITE_"+key, key)
	cmd.Stdin = strings.NewReader(value)
	return cmd.Run()
}

func openURL(url string) { _ = browser.OpenURL(url) }
