package main

import "fmt"

// purge is cf purge, the task cf:purge: everything in the CDN cache for ubuntusoftware.net purged.
func purge() {
	if err := api("POST", "/zones/"+zoneID+"/purge_cache", map[string]bool{"purge_everything": true}, nil); err != nil {
		fail(err.Error())
	}
	fmt.Println("✓ Cache purged (ubuntusoftware.net)")
}
