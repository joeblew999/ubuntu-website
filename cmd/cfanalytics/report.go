package main

import "os"

// report is cfanalytics report, the task cfanalytics:report: the traffic report (--days, -v, --webhook to post it).
func report() { os.Exit(run(os.Args[2:])) }
