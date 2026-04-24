package filter

import "testing"

func TestNoRulesAllowsAll(t *testing.T) {
	f := New(nil)
	if !f.Relevant(80, "tcp") {
		t.Error("expected port 80/tcp to be relevant with no rules")
	}
	if !f.Relevant(443, "tcp") {
		t.Error("expected port 443/tcp to be relevant with no rules")
	}
}

func TestIgnoreRuleHidesPort(t *testing.T) {
	f := New([]Rule{
		{Port: 22, Protocol: "tcp", Allow: false},
	})
	if f.Relevant(22, "tcp") {
		t.Error("expected port 22/tcp to be hidden by ignore rule")
	}
	// Other ports should still be relevant (no whitelist mode).
	if !f.Relevant(80, "tcp") {
		t.Error("expected port 80/tcp to be relevant when not in whitelist mode")
	}
}

func TestWhitelistModeOnlyAllowsListed(t *testing.T) {
	f := New([]Rule{
		{Port: 80, Protocol: "tcp", Allow: true},
		{Port: 443, Protocol: "tcp", Allow: true},
	})

	if !f.Relevant(80, "tcp") {
		t.Error("expected 80/tcp to be relevant")
	}
	if !f.Relevant(443, "tcp") {
		t.Error("expected 443/tcp to be relevant")
	}
	if f.Relevant(8080, "tcp") {
		t.Error("expected 8080/tcp to be hidden in whitelist mode")
	}
}

func TestIgnoreOverridesAllowForSamePort(t *testing.T) {
	// An explicit ignore rule wins even when allow rules exist.
	f := New([]Rule{
		{Port: 80, Protocol: "tcp", Allow: true},
		{Port: 22, Protocol: "tcp", Allow: false},
	})

	if f.Relevant(22, "tcp") {
		t.Error("expected 22/tcp to be hidden by ignore rule")
	}
	if !f.Relevant(80, "tcp") {
		t.Error("expected 80/tcp to be relevant")
	}
}

func TestProtocolIsDistinct(t *testing.T) {
	f := New([]Rule{
		{Port: 53, Protocol: "tcp", Allow: true},
	})
	// udp/53 is not in the whitelist.
	if f.Relevant(53, "udp") {
		t.Error("expected 53/udp to be hidden; only 53/tcp is whitelisted")
	}
}

func TestString(t *testing.T) {
	f := New(nil)
	if f.String() != "filter: all ports" {
		t.Errorf("unexpected string for empty filter: %s", f.String())
	}

	f2 := New([]Rule{{Port: 80, Protocol: "tcp", Allow: true}})
	if f2.String() != "filter: 1 rule(s)" {
		t.Errorf("unexpected string: %s", f2.String())
	}
}
