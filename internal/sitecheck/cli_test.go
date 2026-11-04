package sitecheck

import (
	"testing"
	"time"
)

// TestReportNoise: a node or two failing is check-host.net's noise, which the check itself passes; the report must not
// raise an issue for it either (it did: 205 open issues from 2025-12 to 2026-11, nearly all for one or two nodes).
func TestReportNoise(t *testing.T) {
	state := func(failed int) *State {
		return &State{Timestamp: time.Now(), TotalNodes: 48, OKCount: 48 - failed, FailedCount: failed}
	}
	for _, tc := range []struct {
		name      string
		previous  *State
		failed    int
		wantIssue bool
	}{
		{"none", state(0), 0, false},
		{"one node, first run", nil, 1, false},
		{"one node, after none", state(0), 1, false},
		{"two nodes, after none", state(0), 2, false},
		{"three nodes", state(0), 3, true},
		{"three nodes, first run", nil, 3, true},
		{"up by three beyond the noise", state(3), 6, true},
	} {
		if _, got := generateMarkdownReport(state(tc.failed), tc.previous, nil, "id"); got != tc.wantIssue {
			t.Errorf("%s: issue = %v, want %v", tc.name, got, tc.wantIssue)
		}
	}
}
