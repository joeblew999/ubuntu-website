package main

import "os"

// refresh is youtube refresh, the task youtube:refresh: yt-dlp and ffmpeg fetched again (newest).
func refresh() { os.Exit(run(os.Args[1:])) }
