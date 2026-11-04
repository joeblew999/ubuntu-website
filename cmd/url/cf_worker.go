package main

// cfWorker is url cf-worker, the task url:cf-worker: the site's Worker in Cloudflare's dashboard (its deployments,
// builds, domains, logs).
func cfWorker() {
	open("https://dash.cloudflare.com/" + cfAccountID + "/workers/services/view/ubuntusoftware-net/production")
}
