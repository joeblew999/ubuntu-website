package main

import "fmt"

// emailStatus is cf email-status, the task cf:email-status: whether Email Routing is on for ubuntusoftware.net, its
// rules (from → to), and its destinations, verified or not.
func emailStatus() {
	var settings struct {
		Enabled bool `json:"enabled"`
	}
	if err := api("GET", "/zones/"+zoneID+"/email/routing", nil, &settings); err != nil {
		fail(err.Error())
	}
	fmt.Printf("Email Routing for ubuntusoftware.net: enabled=%v\n\nRules:\n", settings.Enabled)
	rs, err := rules()
	if err != nil {
		fail(err.Error())
	}
	for _, r := range rs {
		from, to := "catch-all", "drop"
		if len(r.Matchers) > 0 && r.Matchers[0].Value != "" {
			from = r.Matchers[0].Value
		}
		if len(r.Actions) > 0 && len(r.Actions[0].Value) > 0 {
			to = r.Actions[0].Value[0]
		}
		state := "enabled"
		if !r.Enabled {
			state = "disabled"
		}
		fmt.Printf("  %s → %s [%s]\n", from, to, state)
	}
	as, err := addresses()
	if err != nil {
		fail(err.Error())
	}
	fmt.Println("\nDestinations:")
	for _, a := range as {
		state := "verified"
		if a.Verified == "" {
			state = "pending"
		}
		fmt.Printf("  %s [%s]\n", a.Email, state)
	}
}
