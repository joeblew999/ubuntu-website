package main

// cfAnalytics is url cf-analytics, the task url:cf-analytics: Cloudflare Web Analytics for the site.
func cfAnalytics() {
	open("https://dash.cloudflare.com/" + cfAccountID + "/web-analytics/overview?siteTag~in=" + cfSiteTag)
}
