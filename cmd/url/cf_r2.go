package main

// cfR2 is url cf-r2, the task url:cf-r2: the R2 dashboard (the bucket of large assets, ubuntu-website-assets).
func cfR2() { open("https://dash.cloudflare.com/" + cfAccountID + "/r2/overview") }
