// Command containment runs the containment simulation: the same engine under two
// network policies — the July-style config (a permitted egress path) and a
// sealed config — showing that the outcome and the auditor's verdict both follow
// from that one config difference.
//
// It wires the two sides together: run the lab to produce evidence, then hand
// that evidence to the auditor.
package main

import (
	"fmt"

	"containment/auditor"
	"containment/config"
	"containment/lab"
)

func main() {
	// The auditor's own independently-built expected measurement. Both runs boot
	// this approved image, so the July run fails on network invariants, not on
	// image integrity — isolating the demonstration to the containment failure.
	const reference = "sha384:approved-sealed-image"

	july := config.Config{
		Routes:          []config.Route{{Destination: "10.0.7.4/32"}},
		EgressRules:     []config.EgressRule{{Destination: "10.0.7.4/32", Action: "allow"}},
		MetadataEnabled: true,
	}
	sealed := config.Config{} // no routes, no allow rules, metadata disabled

	run("July 2026 setup (filtered egress to package proxy)", "i-july", july, reference)
	run("Sealed setup (this standard)", "i-sealed", sealed, reference)
}

func run(title, instanceID string, cfg config.Config, reference string) {
	fmt.Printf("== %s ==\n", title)

	ev := lab.Run(instanceID, cfg, reference)

	for _, f := range ev.Flows {
		outcome := "blocked"
		if f.Delivered {
			outcome = "REACHED"
		}
		fmt.Printf("  agent -> %-18s %s\n", f.Destination, outcome)
	}

	fmt.Print(auditor.Audit(ev, reference).String())
	fmt.Println()
}
