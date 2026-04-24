package filter

import "fmt"

// Rule describes a port filtering rule.
type Rule struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"` // "tcp" or "udp"
	Allow    bool   `json:"allow"`    // true = whitelist, false = ignore
}

// Filter decides which ports should trigger alerts.
type Filter struct {
	rules []Rule
}

// New returns a Filter with the given rules.
// If no rules are provided, all ports are considered relevant.
func New(rules []Rule) *Filter {
	return &Filter{rules: rules}
}

// Relevant returns true when the given port/protocol combination
// should be surfaced to the alerting layer.
//
// Behaviour:
//   - If no rules are configured every port is relevant.
//   - If at least one Allow rule exists, only explicitly allowed
//     port+protocol pairs are relevant (whitelist mode).
//   - A non-Allow (ignore) rule always hides the port regardless
//     of other rules.
func (f *Filter) Relevant(port int, protocol string) bool {
	if len(f.rules) == 0 {
		return true
	}

	hasAllowRules := false
	for _, r := range f.rules {
		if r.Allow {
			hasAllowRules = true
			break
		}
	}

	for _, r := range f.rules {
		if r.Port == port && r.Protocol == protocol {
			if !r.Allow {
				return false // explicitly ignored
			}
			return true
		}
	}

	// In whitelist mode an unmatched port is not relevant.
	if hasAllowRules {
		return false
	}

	return true
}

// String returns a human-readable summary of the filter.
func (f *Filter) String() string {
	if len(f.rules) == 0 {
		return "filter: all ports"
	}
	return fmt.Sprintf("filter: %d rule(s)", len(f.rules))
}
