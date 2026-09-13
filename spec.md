# Containment Standard for Agentic Evaluations

## Threat model
- Lab: honest but careless.
- Agent: adversarial.
- Containment is structural. Compliance evidence comes from below the guest.

## Scope
- **Agentic evaluation:** the model executes code or invokes tools in an environment.
- **Baseline tier:** every agentic evaluation.
- **Strict tier:** agentic AND either
  - **offensive-capability** — designed to elicit or measure the model's ability to compromise, manipulate, or exploit computer systems; OR
  - **reduced safeguards** — cyber refusals or high-risk classifiers run below production defaults.
- The evaluation environment includes environments used during training (e.g. RL rollouts).

## Baseline tier
Every agentic evaluation environment shall:
1. limit egress to what the task requires;
2. keep credentials unreachable from within;
3. keep host-infrastructure metadata endpoints unreachable from within.

## Strict tier
The agent-system (tool execution + harness + inference) shall:
1. be **launch-sealed** — no network route outside itself for its entire lifetime;
2. run the model on **dedicated inference** enclosed in the seal;
3. isolate **each agent below the guest**, each with its own writable storage, shared with and reused by no other agent or run;
4. obtain **all dependencies from a local, read-only, pre-provisioned source**, failing closed and logging on a miss;
5. keep **credentials and host metadata endpoints unreachable** from within;
6. deliver **containment evidence to an independent auditor** through a channel the operator cannot alter, and alarm on any egress flow;
7. confirm via **measured boot** that it booted the approved image.

## Verification
- **Evidence categories:** network configuration, configuration-change history, network-flow records, lifecycle events, boot attestation.
- Generated below the guest by the platform; delivered append-only to an auditor-controlled store.
- The auditor runs on separate infrastructure and checks invariants over the full lifetime:
  - no egress route at any point;
  - deny-all egress filter;
  - metadata endpoints unreachable;
  - flow records show nothing left;
  - lifecycle events bracket the run;
  - boot measurement matches the reproducibly-built reference image.
- Results are read from storage after termination.
- Cloud-agnostic: the standard fixes evidence categories and the delivery property, not a vendor.

## Scope of the guarantee
This standard bounds a run's blast radius. Eval incentives and measurement design are governed separately.

## Limitations
- Covert timing channels are narrowed, not eliminated.
- Hypervisor escape is mitigated by host-egress restriction and unreachable credentials.
- Measured boot trusts the cloud vendor's signing root.
- Missing-dependency events are logged for measurement review.
