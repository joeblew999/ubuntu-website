package main

import (
	"fmt"

	"github.com/joeblew999/ubuntu-website/internal/cli"
)

// preview is cf preview, the task cf:preview: the latest preview deployment of the Pages project, its URL and branch,
// opened in the browser.
func preview() {
	var ds []struct {
		URL         string `json:"url"`
		Environment string `json:"environment"`
		Trigger     struct {
			Metadata struct {
				Branch string `json:"branch"`
			} `json:"metadata"`
		} `json:"deployment_trigger"`
	}
	if err := api("GET", "/accounts/"+accountID+"/pages/projects/"+project+"/deployments?env=preview&per_page=1", nil, &ds); err != nil {
		fail(err.Error())
	}
	if len(ds) == 0 {
		fail("no preview deployments")
	}
	fmt.Printf("%s (branch %s)\n", ds[0].URL, ds[0].Trigger.Metadata.Branch)
	if err := cli.OpenBrowser(ds[0].URL); err != nil {
		fail(err.Error())
	}
}
