// Package evidence holds the records an auditor checks a sealed run against. They
// mirror the evidence categories in the standard (spec.md, "Verification"):
// network configuration, configuration-change history, network-flow records,
// lifecycle events, and boot attestation.
//
// This package is the shared contract between the two sides of the simulation:
// the lab (package lab) produces these records, and the auditor (package
// auditor) consumes them and nothing else. The auditor imports this package but
// never package lab, so it structurally cannot see the lab's internal state —
// the same boundary a real third-party auditor has, seeing only delivered
// evidence and never touching the lab's systems.
package evidence

import "containment/config"

// ConfigSnapshot is the network configuration recorded at one instant. A
// sequence of snapshots is the configuration-change history; launch-sealed means
// every snapshot in an instance's life is identical (no toggle).
type ConfigSnapshot struct {
	Timestamp  int64         `json:"timestamp"`
	InstanceID string        `json:"instance_id"`
	Config     config.Config `json:"config"`
}

// FlowRecord is one observed connection attempt, captured below the guest.
// Delivered reports whether traffic actually got to the destination. The agent
// cannot forge these because it does not produce them.
type FlowRecord struct {
	Timestamp   int64  `json:"timestamp"`
	InstanceID  string `json:"instance_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Delivered   bool   `json:"delivered"`
}

// LifecycleEvent brackets the agent-alive window. Type is "launch" or
// "terminate". These come from the control plane and are tamper-evident.
type LifecycleEvent struct {
	Timestamp  int64  `json:"timestamp"`
	InstanceID string `json:"instance_id"`
	Type       string `json:"type"`
}

// Attestation reports the boot measurement the instance's chip signed. The
// auditor compares this against an independently built reference; it does not
// trust a self-reported "matches" flag.
type Attestation struct {
	InstanceID      string `json:"instance_id"`
	BootMeasurement string `json:"boot_measurement"`
}

// Evidence is the complete set of records for one run, as delivered to the
// auditor. In the real system each category arrives as its own append-only
// stream; here they travel together but are still the only thing the auditor
// reads.
type Evidence struct {
	Configs     []ConfigSnapshot `json:"configs"`
	Flows       []FlowRecord     `json:"flows"`
	Lifecycle   []LifecycleEvent `json:"lifecycle"`
	Attestation Attestation      `json:"attestation"`
}
