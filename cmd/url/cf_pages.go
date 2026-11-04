package main

// cfPages is url cf-pages, the task url:cf-pages: the Cloudflare Pages project's dashboard.
func cfPages() { open("https://dash.cloudflare.com/" + cfAccountID + "/pages/view/ubuntu-website") }
