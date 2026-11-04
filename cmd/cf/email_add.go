package main

import "fmt"

// emailAdd is cf email-add <from> <to>, the task cf:email-add: mail to <from> (an address of ubuntusoftware.net)
// forwarded to <to>. A destination Cloudflare hasn't verified yet is added first, and it emails <to> to verify: run
// it again once that's done.
func emailAdd(args []string) {
	if len(args) != 2 {
		fail("usage: cf email-add <from> <to>")
	}
	from, to := args[0], args[1]
	as, err := addresses()
	if err != nil {
		fail(err.Error())
	}
	verified, known := false, false
	for _, a := range as {
		if a.Email == to {
			known, verified = true, a.Verified != ""
		}
	}
	if !verified {
		if !known {
			if err := api("POST", "/accounts/"+accountID+"/email/routing/addresses", map[string]string{"email": to}, nil); err != nil {
				fail(err.Error())
			}
		}
		fmt.Printf("%s isn't verified yet: Cloudflare emailed it a link. Run this again once it's verified.\n", to)
		return
	}
	body := map[string]any{
		"name":     "Forward " + from + " to " + to,
		"enabled":  true,
		"matchers": []map[string]string{{"type": "literal", "field": "to", "value": from}},
		"actions":  []map[string]any{{"type": "forward", "value": []string{to}}},
	}
	if err := api("POST", "/zones/"+zoneID+"/email/routing/rules", body, nil); err != nil {
		fail(err.Error())
	}
	fmt.Printf("✓ %s → %s\n", from, to)
}
