package main

import "os"

// avatar is genlogo avatar, the task genlogo:avatar: the Bluesky avatar, static/images/bluesky-avatar.png.
func avatar() { os.Exit(run(append([]string{"-asset", "avatar"}, os.Args[2:]...))) }
