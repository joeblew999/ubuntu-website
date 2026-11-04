package main

// tcp is sitecheck tcp, the task sitecheck:tcp: the site's port 443 answers, from every node.
func tcp(args []string) { check("tcp", args) }
