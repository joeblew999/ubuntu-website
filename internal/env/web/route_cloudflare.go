package web

import (
	"github.com/go-via/via"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// cloudflarePage redirects to step 1 of the setup wizard
func cloudflarePage(c *via.Context, cfg *env.EnvConfig, mockMode bool) {
	cloudflareStep1Page(c, cfg, mockMode)
}
