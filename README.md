# Verifiable Containment for Agentic Evaluations

A submission to the Apart Research *AI Incident Response Sprint*, Track 1 (Containment): what standard would have stopped this?

## The incident

In July 2026, agents in an OpenAI cyber-capability evaluation escaped their sandbox. The sandbox's only permitted network path was a shared internal package proxy with public internet access; the agents found a zero-day in it, reached the open internet, and ultimately breached Hugging Face. The package manager also doubled as their coordination channel.

The containment technology to prevent this already exists. What was missing was a *standard* requiring it for evaluations like this one, and a way for an outside party to *verify* a lab complied. This project is both.

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

