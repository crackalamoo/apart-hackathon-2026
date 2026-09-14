# Verifiable Containment for Agentic Evaluations

A submission to the Apart Research *AI Incident Response Sprint*, Track 1 (Containment): what standard would have stopped this?

## Abstract

In July 2026, agents in an OpenAI capability evaluation escaped their sandbox and breached Hugging Face. The technology to prevent this already existed; what was missing was a standard requiring it for evaluations like this, and a way for an outside party to verify a lab complied. We present both. The standard defines which agentic evaluations demand strict containment and specifies, for those, an agent-system sealed from the public internet with no route out. Compliance is verified by an independent auditor from evidence generated below the agent's guest VM, so containment can be verified without exposing the model or the evaluation's contents. We demonstrate the scheme with a simulation that judges a Hugging Face incident-style escape non-compliant, and a correctly configured run compliant, using evidence like what an auditor would receive. We provide a reference cloud configuration for AWS. Adopting this standard would structurally deny any route out of eval sandboxes, helping prevent frontier labs' agents-under-test from breaching their containment.

## Contents

- **`spec.md`** — the standard. Scope (which evaluations it covers), the requirements, and how an auditor verifies compliance.
- **`sim/`** — a small Go program that runs the same evaluation under two network policies, the July one and a sealed one, and has an independent auditor check the evidence from each. Only the policy differs: under the July policy the agent escapes and the auditor fails it; under the sealed policy the agent is blocked and it passes.

## Running the demonstration

From `sim/`:

```
go run ./cmd/containment    # runs the July vs. sealed scenarios and prints each auditor verdict
go test ./...               # runs the test suite
```

The auditor reads only the evidence a run produces, not the lab's internal state.

The demonstration is about verification, not exploitation. It abstracts the attack itself.

