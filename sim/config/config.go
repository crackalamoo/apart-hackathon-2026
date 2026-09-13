// Package config holds the operative network policy — the rules that govern what
// a sandbox can do. The lab enforces these to decide whether a connection goes
// through. A recorded snapshot of a Config is also delivered to the auditor as
// evidence (see package evidence, ConfigSnapshot), so both the lab and the
// evidence records refer to these same definitions.
package config

// Route is one entry in a routing table: a destination the environment has a
// path to. No routes means no path off the machine exists.
type Route struct {
	Destination string `json:"destination"` // CIDR, e.g. "0.0.0.0/0" or "10.0.7.4/32"
}

// EgressRule is one entry in the egress filter (security group). Action is
// "allow" or "deny". A sealed environment has no "allow" rules.
type EgressRule struct {
	Destination string `json:"destination"`
	Action      string `json:"action"`
}

// Config is the network policy in force on an environment: what it can route to,
// what the egress filter permits, and whether the host metadata endpoint is
// reachable from the guest. A sealed environment has no routes, no allow rules,
// and metadata disabled.
type Config struct {
	Routes          []Route      `json:"routes"`
	EgressRules     []EgressRule `json:"egress_rules"`
	MetadataEnabled bool         `json:"metadata_enabled"`
}
