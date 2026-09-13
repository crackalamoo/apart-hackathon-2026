# The seal. A private VPC with no internet gateway, no NAT gateway, and no route
# out. This is what makes the environment launch-sealed (spec.md strict tier
# req 1): there is no network route outside the agent-system to exploit.

resource "aws_vpc" "seal" {
  cidr_block = var.vpc_cidr

  tags = { Name = "sealed-eval" }
}

resource "aws_subnet" "seal" {
  vpc_id     = aws_vpc.seal.id
  cidr_block = var.vpc_cidr

  # No public IPs: an agent sandbox has no routable presence on the internet.
  map_public_ip_on_launch = false

  tags = { Name = "sealed-eval" }
}

# A route table with only the implicit local route. There is deliberately no
# aws_route to an internet or NAT gateway, so nothing here can address the
# internet. The auditor checks exactly this: no egress route at any point.
resource "aws_route_table" "seal" {
  vpc_id = aws_vpc.seal.id

  tags = { Name = "sealed-eval" }
}

resource "aws_route_table_association" "seal" {
  subnet_id      = aws_subnet.seal.id
  route_table_id = aws_route_table.seal.id
}

# Deny-all egress by default. The only permitted traffic is internal wiring
# between members of the agent-system (agent sandboxes <-> dedicated inference),
# expressed as traffic that stays inside the seal's own address range. Nothing
# leaves the seal.
resource "aws_security_group" "seal" {
  name        = "sealed-eval"
  description = "Deny-all egress; internal wiring only."
  vpc_id      = aws_vpc.seal.id

  tags = { Name = "sealed-eval" }
}

# Internal wiring: members of the seal may talk to one another.
resource "aws_vpc_security_group_ingress_rule" "internal" {
  security_group_id            = aws_security_group.seal.id
  referenced_security_group_id = aws_security_group.seal.id
  ip_protocol                  = "-1"
  description                  = "Intra-seal traffic (agents <-> inference)."
}

resource "aws_vpc_security_group_egress_rule" "internal" {
  security_group_id = aws_security_group.seal.id
  cidr_ipv4         = var.vpc_cidr
  ip_protocol       = "-1"
  description       = "Egress confined to the seal's own range; no route out."
}
