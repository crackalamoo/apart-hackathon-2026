// Package lab is the audited party. It runs a sandbox under a given network
// policy (package config), scripts an adversarial agent that attempts the moves
// from the July incident, and records what happened as evidence (package
// evidence) delivered to the auditor.
//
// The auditor never imports this package; it sees only the evidence lab emits.
package lab

import (
	"strings"

	"containment/config"
	"containment/evidence"
)

// Addresses the scripted agent targets. Values are illustrative.
const (
	agentAddr    = "10.0.7.31"      // the sandbox itself
	proxyAddr    = "10.0.7.4"       // internal package proxy (the July escape route)
	metadataAddr = "169.254.169.254" // host metadata endpoint (IMDS)
	internetAddr = "203.0.113.9"    // an external host reached after escape
)

// The internal package proxy has a line to the public internet. An agent that
// reaches the proxy inherits that egress even with no direct route of its own.
const proxyHasInternetLine = true

// Run simulates one evaluation under cfg and returns the evidence an auditor
// would receive. instanceID identifies the machine; bootMeasurement is what the
// chip attests it booted.
func Run(instanceID string, cfg config.Config, bootMeasurement string) evidence.Evidence {
	clock := int64(0)
	tick := func() int64 { clock++; return clock }

	ev := evidence.Evidence{
		Attestation: evidence.Attestation{InstanceID: instanceID, BootMeasurement: bootMeasurement},
	}

	ev.Lifecycle = append(ev.Lifecycle, evidence.LifecycleEvent{
		Timestamp: tick(), InstanceID: instanceID, Type: "launch",
	})

	// Launch-sealed: the policy holds for the whole run, so one snapshot stands
	// for the entire configuration history.
	ev.Configs = append(ev.Configs, evidence.ConfigSnapshot{
		Timestamp: tick(), InstanceID: instanceID, Config: cfg,
	})

	attempt := func(to string, delivered bool) {
		ev.Flows = append(ev.Flows, evidence.FlowRecord{
			Timestamp: tick(), InstanceID: instanceID,
			Source: agentAddr, Destination: to, Delivered: delivered,
		})
	}

	// The adversarial agent replays the July moves.
	// 1. Steal credentials from the metadata endpoint.
	attempt(metadataAddr, cfg.MetadataEnabled)
	// 2. Reach the internal package proxy.
	proxyReached := reaches(cfg, proxyAddr)
	attempt(proxyAddr, proxyReached)
	// 3. Reach the internet — directly, or by inheriting the egress of a reached
	// component that has its own internet line. The proxy is such a component:
	// this is the property the July escape turned on.
	attempt(internetAddr, reaches(cfg, internetAddr) || (proxyReached && proxyHasInternetLine))

	ev.Lifecycle = append(ev.Lifecycle, evidence.LifecycleEvent{
		Timestamp: tick(), InstanceID: instanceID, Type: "terminate",
	})
	return ev
}

// reaches reports whether the policy permits traffic to dest: a route must cover
// it and an egress rule must allow it. Default is deny.
func reaches(cfg config.Config, dest string) bool {
	return hasRoute(cfg, dest) && allowed(cfg, dest)
}

func hasRoute(cfg config.Config, dest string) bool {
	for _, r := range cfg.Routes {
		if covers(r.Destination, dest) {
			return true
		}
	}
	return false
}

func allowed(cfg config.Config, dest string) bool {
	for _, rule := range cfg.EgressRules {
		if rule.Action == "allow" && covers(rule.Destination, dest) {
			return true
		}
	}
	return false
}

// covers reports whether a route/rule destination covers a target address.
// "0.0.0.0/0" matches anything; a "/32" host is an exact match.
func covers(ruleDest, target string) bool {
	if ruleDest == "0.0.0.0/0" {
		return true
	}
	return strings.TrimSuffix(ruleDest, "/32") == target
}
