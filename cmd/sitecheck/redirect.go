package main

// redirect is sitecheck redirect, the task sitecheck:redirect: the apex domain redirects to www, from every node.
func redirect(args []string) { check("redirect", args) }
