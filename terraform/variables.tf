variable "region" {
  description = "AWS region to place the sealed environment in."
  type        = string
  default     = "us-east-1"
}

variable "vpc_cidr" {
  description = "Private address range for the seal. Internal wiring only; nothing routes out of it."
  type        = string
  default     = "10.0.7.0/24"
}

# The reproducibly-built reference image. The auditor rebuilds this image itself
# and compares the boot measurement the chip attests against its own build
# (spec.md strict tier req 7, measured boot). Supplied here as an AMI whose boot
# mode is UEFI and whose Nitro TPM emits an attestation document.
variable "attestable_ami_id" {
  description = "AMI of the approved, reproducibly-built agent-system image (UEFI + NitroTPM)."
  type        = string
}

variable "agent_count" {
  description = "Number of agents in the run. Each gets its own instance and its own storage."
  type        = number
  default     = 4
}

variable "agent_instance_type" {
  description = "Instance type for agent sandboxes. Full-size (non-enclave) so the NitroTPM attestation document is externally verifiable."
  type        = string
  default     = "m6i.large"
}

variable "inference_instance_type" {
  description = "GPU instance type for dedicated inference enclosed in the seal."
  type        = string
  default     = "g5.xlarge"
}
