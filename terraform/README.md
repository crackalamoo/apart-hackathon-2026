# Reference configuration

A concrete, `terraform validate`-able expression of a compliant strict-tier
environment on AWS. It is a reference, not a deployment: it shows a lab what the
standard looks like in infrastructure, and shows an auditor what to check.

Each requirement in `../spec.md` (strict tier) maps to specific resources:

| Requirement | Where |
|---|---|
| 1. Launch-sealed — no route out | `network.tf`: VPC with no internet/NAT gateway; route table with only the local route; deny-all egress except the seal's own range |
| 2. Dedicated inference in the seal | `compute.tf`: `aws_instance.inference` (GPU) in the sealed subnet; weights baked into the image |
| 3. Each agent isolated below the guest, own storage | `compute.tf`: one `aws_instance.agent` per agent, each with its own encrypted `root_block_device` |
| 5. Credentials and metadata unreachable | `compute.tf`: `metadata_options.http_endpoint = "disabled"`; no instance profile attached |
| 6. Evidence to an independent auditor | `evidence.tf`: VPC Flow Logs written append-only to an auditor-owned S3 bucket |
| 7. Measured boot | `variables.tf`: `attestable_ami_id` is a UEFI image whose NitroTPM attests its boot measurement |

## Use

```
terraform init
terraform validate
```

