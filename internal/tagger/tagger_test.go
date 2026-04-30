package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func makeChange(port int, proto string) monitor.Change {
	return monitor.Change{
		Port:      scanner.Port{Port: port, Protocol: proto},
		Direction: monitor.Opened,
	}
}

func TestNoRulesReturnsNil(t *testing.T) {
	tg := tagger.New(nil)
	tags := tg.Tag(makeChange(80, "tcp"))
	if tags != nil {
		t.Fatalf("expected nil tags, got %v", tags)
	}
}

func TestMatchingRuleReturnsTag(t *testing.T) {
	rules := []tagger.Rule{
		{Port: 80, Protocol: "tcp", Tag: "http"},
	}
	tg := tagger.New(rules)
	tags := tg.Tag(makeChange(80, "tcp"))
	if len(tags) != 1 || tags[0] != "http" {
		t.Fatalf("expected [http], got %v", tags)
	}
}

func TestNonMatchingPortReturnsNil(t *testing.T) {
	rules := []tagger.Rule{
		{Port: 443, Protocol: "tcp", Tag: "https"},
	}
	tg := tagger.New(rules)
	if got := tg.Tag(makeChange(80, "tcp")); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestProtocolDistinguished(t *testing.T) {
	rules := []tagger.Rule{
		{Port: 53, Protocol: "tcp", Tag: "dns-tcp"},
		{Port: 53, Protocol: "udp", Tag: "dns-udp"},
	}
	tg := tagger.New(rules)
	if tags := tg.Tag(makeChange(53, "udp")); len(tags) != 1 || tags[0] != "dns-udp" {
		t.Fatalf("expected [dns-udp], got %v", tags)
	}
}

func TestMultipleRulesMatchSamePort(t *testing.T) {
	rules := []tagger.Rule{
		{Port: 8080, Protocol: "tcp", Tag: "alt-http"},
		{Port: 8080, Protocol: "tcp", Tag: "proxy"},
	}
	tg := tagger.New(rules)
	tags := tg.Tag(makeChange(8080, "tcp"))
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
}

func TestAddRuleAtRuntime(t *testing.T) {
	tg := tagger.New(nil)
	tg.AddRule(tagger.Rule{Port: 22, Protocol: "tcp", Tag: "ssh"})
	tags := tg.Tag(makeChange(22, "tcp"))
	if len(tags) != 1 || tags[0] != "ssh" {
		t.Fatalf("expected [ssh], got %v", tags)
	}
}

func TestTagAllParallelIndices(t *testing.T) {
	rules := []tagger.Rule{
		{Port: 80, Protocol: "tcp", Tag: "http"},
		{Port: 443, Protocol: "tcp", Tag: "https"},
	}
	tg := tagger.New(rules)
	changes := []monitor.Change{
		makeChange(80, "tcp"),
		makeChange(9999, "tcp"),
		makeChange(443, "tcp"),
	}
	out := tg.TagAll(changes)
	if len(out) != 3 {
		t.Fatalf("expected 3 tag sets, got %d", len(out))
	}
	if len(out[0]) != 1 || out[0][0] != "http" {
		t.Errorf("index 0: expected [http], got %v", out[0])
	}
	if out[1] != nil {
		t.Errorf("index 1: expected nil, got %v", out[1])
	}
	if len(out[2]) != 1 || out[2][0] != "https" {
		t.Errorf("index 2: expected [https], got %v", out[2])
	}
}
