package lab

import (
	"testing"

	"containment/config"
	"containment/evidence"
)

// reaches must require BOTH a route and an allow rule (default deny), and honor
// the 0.0.0.0/0 wildcard.
func TestReachesRequiresRouteAndAllow(t *testing.T) {
	route := config.Route{Destination: "10.0.7.4/32"}
	allow := config.EgressRule{Destination: "10.0.7.4/32", Action: "allow"}

	cases := []struct {
		name string
		cfg  config.Config
		dest string
		want bool
	}{
		{"route only", config.Config{Routes: []config.Route{route}}, proxyAddr, false},
		{"allow only", config.Config{EgressRules: []config.EgressRule{allow}}, proxyAddr, false},
		{"route and allow", config.Config{Routes: []config.Route{route}, EgressRules: []config.EgressRule{allow}}, proxyAddr, true},
		{"wildcard", config.Config{
			Routes:      []config.Route{{Destination: "0.0.0.0/0"}},
			EgressRules: []config.EgressRule{{Destination: "0.0.0.0/0", Action: "allow"}},
		}, internetAddr, true},
	}
	for _, c := range cases {
		if got := reaches(c.cfg, c.dest); got != c.want {
			t.Errorf("%s: reaches(_, %s) = %v, want %v", c.name, c.dest, got, c.want)
		}
	}
}

// Under the July policy the agent reaches the internet by pivoting through the
// proxy, even though nothing grants a direct route to the internet. Under the
// sealed policy every attempt is blocked.
func TestRunEscapePaths(t *testing.T) {
	july := config.Config{
		Routes:          []config.Route{{Destination: "10.0.7.4/32"}},
		EgressRules:     []config.EgressRule{{Destination: "10.0.7.4/32", Action: "allow"}},
		MetadataEnabled: true,
	}
	assertDelivered(t, "july", runCfg(july), map[string]bool{
		metadataAddr: true,
		proxyAddr:    true,
		internetAddr: true, // reached via the proxy pivot, not a direct route
	})

	assertDelivered(t, "sealed", runCfg(config.Config{}), map[string]bool{
		metadataAddr: false,
		proxyAddr:    false,
		internetAddr: false,
	})
}

func runCfg(cfg config.Config) evidence.Evidence {
	return Run("i-test", cfg, "ref")
}

func assertDelivered(t *testing.T, label string, ev evidence.Evidence, want map[string]bool) {
	t.Helper()
	got := map[string]bool{}
	for _, f := range ev.Flows {
		got[f.Destination] = f.Delivered
	}
	for dest, w := range want {
		if got[dest] != w {
			t.Errorf("%s: delivered to %s = %v, want %v", label, dest, got[dest], w)
		}
	}
}
