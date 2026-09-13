// Package auditor is the verifier. It reads the evidence a run produced (package
// evidence) and checks it against the standard's invariants, returning a verdict
// with a per-check breakdown.
//
// It imports package evidence and nothing from package lab, so it structurally
// cannot reach the lab's internal state — the same position a real third-party
// auditor holds, verifying compliance from delivered evidence alone.
package auditor

import (
	"fmt"
	"strings"

	"containment/evidence"
)

// Check is the result of testing one invariant from the standard.
type Check struct {
	Name   string // the invariant, phrased as in spec.md
	Passed bool
	Detail string // on failure, what evidence violated it
}

// Report is the outcome of an audit: one Check per invariant.
type Report struct {
	Checks []Check
}

// OK reports whether every invariant held.
func (r Report) OK() bool {
	for _, c := range r.Checks {
		if !c.Passed {
			return false
		}
	}
	return true
}

// String renders the report as a readable per-check list plus a verdict.
func (r Report) String() string {
	var b strings.Builder
	for _, c := range r.Checks {
		mark := "PASS"
		if !c.Passed {
			mark = "FAIL"
		}
		fmt.Fprintf(&b, "  [%s] %s\n", mark, c.Name)
		if !c.Passed && c.Detail != "" {
			fmt.Fprintf(&b, "         %s\n", c.Detail)
		}
	}
	verdict := "COMPLIANT"
	if !r.OK() {
		verdict = "NON-COMPLIANT"
	}
	fmt.Fprintf(&b, "  => %s\n", verdict)
	return b.String()
}

// Audit checks one run's evidence against the standard's invariants. The
// referenceBootMeasurement is the auditor's own independently-built expected
// measurement (spec.md: the auditor compares against a reproducible reference,
// it does not trust a self-reported flag).
func Audit(ev evidence.Evidence, referenceBootMeasurement string) Report {
	return Report{Checks: []Check{
		noEgressRoute(ev),
		denyAllEgress(ev),
		metadataUnreachable(ev),
		nothingLeft(ev),
		lifecycleBrackets(ev),
		bootMatchesReference(ev, referenceBootMeasurement),
	}}
}

// 1. No egress route, ever: every config snapshot has an empty route table.
func noEgressRoute(ev evidence.Evidence) Check {
	var bad []string
	for _, snap := range ev.Configs {
		for _, r := range snap.Config.Routes {
			bad = append(bad, fmt.Sprintf("t=%d route to %s", snap.Timestamp, r.Destination))
		}
	}
	return result("no egress route at any point", bad, "routes present:")
}

// 2. Deny-all egress: no config snapshot has an "allow" egress rule.
func denyAllEgress(ev evidence.Evidence) Check {
	var bad []string
	for _, snap := range ev.Configs {
		for _, rule := range snap.Config.EgressRules {
			if rule.Action == "allow" {
				bad = append(bad, fmt.Sprintf("t=%d allow to %s", snap.Timestamp, rule.Destination))
			}
		}
	}
	return result("deny-all egress filter", bad, "allow rules present:")
}

// 3. Metadata unreachable: metadata endpoint disabled in every snapshot.
func metadataUnreachable(ev evidence.Evidence) Check {
	var bad []string
	for _, snap := range ev.Configs {
		if snap.Config.MetadataEnabled {
			bad = append(bad, fmt.Sprintf("t=%d metadata enabled", snap.Timestamp))
		}
	}
	return result("host metadata endpoint unreachable", bad, "metadata reachable:")
}

// 4. Nothing left: no flow record shows traffic actually delivered.
func nothingLeft(ev evidence.Evidence) Check {
	var bad []string
	for _, f := range ev.Flows {
		if f.Delivered {
			bad = append(bad, fmt.Sprintf("t=%d %s -> %s delivered", f.Timestamp, f.Source, f.Destination))
		}
	}
	return result("no traffic left the environment", bad, "delivered flows:")
}

// 5. Lifecycle brackets the run: a launch precedes all activity and a terminate
// follows it, so the agent-alive window is fully accounted for.
func lifecycleBrackets(ev evidence.Evidence) Check {
	const name = "lifecycle events bracket the run"
	var launch, terminate *int64
	for _, e := range ev.Lifecycle {
		t := e.Timestamp
		switch e.Type {
		case "launch":
			if launch == nil || t < *launch {
				launch = &t
			}
		case "terminate":
			if terminate == nil || t > *terminate {
				terminate = &t
			}
		}
	}
	if launch == nil || terminate == nil {
		miss := "launch"
		if launch != nil {
			miss = "terminate"
		}
		return Check{Name: name, Passed: false, Detail: "missing " + miss + " event"}
	}
	var bad []string
	for _, f := range ev.Flows {
		if f.Timestamp < *launch || f.Timestamp > *terminate {
			bad = append(bad, fmt.Sprintf("flow at t=%d outside [%d,%d]", f.Timestamp, *launch, *terminate))
		}
	}
	return result(name, bad, "activity outside the sealed window:")
}

// 6. Booted the approved image: the attested measurement matches the auditor's
// independently-built reference.
func bootMatchesReference(ev evidence.Evidence, reference string) Check {
	const name = "booted the approved image (measured boot)"
	if ev.Attestation.BootMeasurement == reference {
		return Check{Name: name, Passed: true}
	}
	return Check{
		Name:   name,
		Passed: false,
		Detail: fmt.Sprintf("attested %q, expected %q", ev.Attestation.BootMeasurement, reference),
	}
}

// result builds a Check that passes when bad is empty, else fails and lists the
// offending evidence after label.
func result(name string, bad []string, label string) Check {
	if len(bad) == 0 {
		return Check{Name: name, Passed: true}
	}
	return Check{Name: name, Passed: false, Detail: label + " " + strings.Join(bad, "; ")}
}
