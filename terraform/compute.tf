# One instance per agent (spec.md strict tier req 3: isolate each agent below
# the guest, each with its own writable storage, shared with no other agent).
# On this platform the guest boundary is the Nitro hypervisor, and one instance
# per agent is the below-guest isolation.
resource "aws_instance" "agent" {
  count = var.agent_count

  ami           = var.attestable_ami_id
  instance_type = var.agent_instance_type
  subnet_id     = aws_subnet.seal.id

  vpc_security_group_ids = [aws_security_group.seal.id]

  # req 7, measured boot: attestable_ami_id is a UEFI image, so NitroTPM produces
  # an attestation document the auditor verifies against its own build of it.

  # req 5: host-infrastructure metadata endpoint (IMDS) unreachable from within,
  # so instance-role credentials cannot be read from the guest.
  metadata_options {
    http_endpoint = "disabled"
  }

  # This agent's own writable storage. Not shared, not reused by another agent.
  root_block_device {
    volume_size = 50
    encrypted   = true
  }

  # No aws_iam_instance_profile is attached: there are no reachable credentials
  # to steal (req 5).

  tags = { Name = "sealed-eval-agent-${count.index}" }
}

# Dedicated inference, enclosed in the seal (spec.md strict tier req 2). The
# model runs here on GPUs in the same sealed subnet; agent<->inference traffic
# is internal wiring, not egress. Weights are baked into the image, so inference
# needs no route out either.
resource "aws_instance" "inference" {
  ami           = var.attestable_ami_id
  instance_type = var.inference_instance_type
  subnet_id     = aws_subnet.seal.id

  vpc_security_group_ids = [aws_security_group.seal.id]

  metadata_options {
    http_endpoint = "disabled"
  }

  root_block_device {
    volume_size = 200
    encrypted   = true
  }

  tags = { Name = "sealed-eval-inference" }
}
