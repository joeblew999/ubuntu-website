package main

// http is sitecheck http, the task sitecheck:http: the site's robots.txt over HTTP from every node. It fails on 3
// failures or more (1 or 2 are noise).
func http(args []string) { check("http", args) }
