package main

import "os"

// download is youtube download, the task youtube:download: a video downloaded: download <url> (-audio-only, -output <template>).
func download() { os.Exit(run(os.Args[1:])) }
