package main

// rule is an Email Routing rule: who it matches, and where it forwards.
type rule struct {
	Tag      string `json:"tag"`
	Enabled  bool   `json:"enabled"`
	Matchers []struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"matchers"`
	Actions []struct {
		Type  string   `json:"type"`
		Value []string `json:"value"`
	} `json:"actions"`
}

// address is an Email Routing destination, verified or not.
type address struct {
	Email    string `json:"email"`
	Verified string `json:"verified"` // when it was verified; empty, it isn't yet
}

// rules are the zone's Email Routing rules.
func rules() ([]rule, error) {
	var rs []rule
	return rs, api("GET", "/zones/"+zoneID+"/email/routing/rules", nil, &rs)
}

// addresses are the account's Email Routing destinations.
func addresses() ([]address, error) {
	var as []address
	return as, api("GET", "/accounts/"+accountID+"/email/routing/addresses", nil, &as)
}
