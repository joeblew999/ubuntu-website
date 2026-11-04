package main

import "os"

// email is genlogo email, the task genlogo:email: the email signature logo, static/images/email-logo.png.
func email() { os.Exit(run(append([]string{"-asset", "email"}, os.Args[2:]...))) }
