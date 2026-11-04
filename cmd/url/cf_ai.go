package main

// cfAi is url cf-ai, the task url:cf-ai: Cloudflare's AI Crawl Control (which AI crawlers may read the site).
func cfAi() { open("https://dash.cloudflare.com/" + cfAccountID + "/ubuntusoftware.net/ai") }
