package main

import "fmt"

// emailRemove is cf email-remove <from>, the task cf:email-remove: the Email Routing rule for <from> removed.
func emailRemove(args []string) {
	if len(args) != 1 {
		fail("usage: cf email-remove <from>")
	}
	rs, err := rules()
	if err != nil {
		fail(err.Error())
	}
	for _, r := range rs {
		if len(r.Matchers) > 0 && r.Matchers[0].Value == args[0] {
			if err := api("DELETE", "/zones/"+zoneID+"/email/routing/rules/"+r.Tag, nil, nil); err != nil {
				fail(err.Error())
			}
			fmt.Println("✓ Removed the rule for " + args[0])
			return
		}
	}
	fail("no rule for " + args[0])
}
