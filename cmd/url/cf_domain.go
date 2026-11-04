package main

// cfDomain is url cf-domain, the task url:cf-domain: the domain's Cloudflare dashboard (DNS, rules, email).
func cfDomain() { open("https://dash.cloudflare.com/" + cfAccountID + "/ubuntusoftware.net") }
