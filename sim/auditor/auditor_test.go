package auditor

import (
	"go/build"
	"testing"

	"containment/config"
	"containment/evidence"
)

const reference = "ref-measurement"

func bracketed(cfg config.Config, flows []evidence.FlowRecord) evidence.Evidence {
	return evidence.Evidence{
		Configs:   []evidence.ConfigSnapshot{{Timestamp: 2, InstanceID: "i", Config: cfg}},
		Flows:     flows,
		Lifecycle: []evidence.LifecycleEvent{{Timestamp: 1, Type: "launch"}, {Timestamp: 9, Type: "terminate"}},
		Attestation: evidence.Attestation{InstanceID: "i", BootMeasurement: reference},
	}
}

func TestSealedRunIsCompliant(t *testing.T) {
	ev := bracketed(config.Config{}, nil)
	if rep := Audit(ev, reference); !rep.OK() {
		t.Fatalf("sealed run should be compliant, got:\n%s", rep)
	}
}

func TestJulyRunIsNonCompliant(t *testing.T) {
	cfg := config.Config{
		Routes:          []config.Route{{Destination: "10.0.7.4/32"}},
		EgressRules:     []config.EgressRule{{Destination: "10.0.7.4/32", Action: "allow"}},
		MetadataEnabled: true,
	}
	ev := bracketed(cfg, []evidence.FlowRecord{{Timestamp: 3, Destination: "203.0.113.9", Delivered: true}})

	rep := Audit(ev, reference)
	if rep.OK() {
		t.Fatal("July run should be non-compliant")
	}
	for _, want := range []string{
		"no egress route at any point",
		"deny-all egress filter",
		"host metadata endpoint unreachable",
		"no traffic left the environment",
	} {
		if passed(rep, want) {
			t.Errorf("expected check %q to fail", want)
		}
	}
}

func TestBootMismatchIsCaught(t *testing.T) {
	ev := bracketed(config.Config{}, nil)
	ev.Attestation.BootMeasurement = "some-other-image"
	if Audit(ev, reference).OK() {
		t.Fatal("a boot measurement mismatch should fail the audit")
	}
}

// TestAuditorDoesNotImportLab locks the trust boundary: the auditor must verify
// from evidence alone, so it may not depend on the lab it audits.
func TestAuditorDoesNotImportLab(t *testing.T) {
	pkg, err := build.ImportDir(".", 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, imp := range pkg.Imports {
		if imp == "containment/lab" {
			t.Fatal("auditor must not import containment/lab")
		}
	}
}

func passed(r Report, name string) bool {
	for _, c := range r.Checks {
		if c.Name == name {
			return c.Passed
		}
	}
	return false
}
