terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  required_version = ">= 1.3.0"
}

provider "aws" {
  region = "ap-southeast-1"
}

# ---------------------
# VPC A Configuration (Public)
# ---------------------
resource "aws_vpc" "vpc_a" {
  cidr_block           = "10.10.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = { Name = "VPC-A" }
}

resource "aws_subnet" "subnet_a_public" {
  vpc_id                  = aws_vpc.vpc_a.id
  cidr_block              = "10.10.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "ap-southeast-1a"
  tags = { Name = "Subnet-A-Public" }
}

resource "aws_internet_gateway" "igw_a" {
  vpc_id = aws_vpc.vpc_a.id
  tags = { Name = "IGW-A" }
}

resource "aws_route_table" "rt_a_public" {
  vpc_id = aws_vpc.vpc_a.id
  tags   = { Name = "RT-A-Public" }
}

resource "aws_route_table_association" "rta_a_public" {
  subnet_id      = aws_subnet.subnet_a_public.id
  route_table_id = aws_route_table.rt_a_public.id
}

resource "aws_route" "public_internet_route_a" {
  route_table_id         = aws_route_table.rt_a_public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.igw_a.id
}

# ---------------------
# VPC B Configuration (Private)
# ---------------------
resource "aws_vpc" "vpc_b" {
  cidr_block           = "10.20.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = { Name = "VPC-B" }
}

resource "aws_subnet" "subnet_b_private" {
  vpc_id            = aws_vpc.vpc_b.id
  cidr_block        = "10.20.1.0/24"
  availability_zone = "ap-southeast-1b"
  tags = { Name = "Subnet-B-Private" }
}

resource "aws_route_table" "rt_b_private" {
  vpc_id = aws_vpc.vpc_b.id
  tags   = { Name = "RT-B-Private" }
}

resource "aws_route_table_association" "rta_b_private" {
  subnet_id      = aws_subnet.subnet_b_private.id
  route_table_id = aws_route_table.rt_b_private.id
}

# ---------------------
# Transit Gateway
# ---------------------
resource "aws_ec2_transit_gateway" "tgw" {
  description = "TGW between VPC A and VPC B"
  tags = { Name = "TGW-A-B" }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "attach_a" {
  subnet_ids         = [aws_subnet.subnet_a_public.id]
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = aws_vpc.vpc_a.id
  tags = { Name = "TGW-Attach-A" }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "attach_b" {
  subnet_ids         = [aws_subnet.subnet_b_private.id]
  transit_gateway_id = aws_ec2_transit_gateway.tgw.id
  vpc_id             = aws_vpc.vpc_b.id
  tags = { Name = "TGW-Attach-B" }
}

# Add TGW routes
resource "aws_route" "tgw_route_a_to_b" {
  route_table_id         = aws_route_table.rt_a_public.id
  destination_cidr_block = aws_vpc.vpc_b.cidr_block
  transit_gateway_id     = aws_ec2_transit_gateway.tgw.id
}

resource "aws_route" "tgw_route_b_to_a" {
  route_table_id         = aws_route_table.rt_b_private.id
  destination_cidr_block = aws_vpc.vpc_a.cidr_block
  transit_gateway_id     = aws_ec2_transit_gateway.tgw.id
}

# ---------------------
# Security Groups
# ---------------------
resource "aws_security_group" "sg_public" {
  name        = "allow-ssh-icmp-public"
  description = "Allow SSH and ICMP"
  vpc_id      = aws_vpc.vpc_a.id

  ingress {
    description = "SSH from anywhere"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "ICMP (ping)"
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "SG-Public" }
}

# Security Group for Private EC2 Instance
resource "aws_security_group" "sg_b_private" {
  name        = "allow-icmp-ssh-from-vpc-a-by-tgw"
  description = "Allow ICMP and SSH from VPC A"
  vpc_id      = aws_vpc.vpc_b.id

  ingress {
    description = "Allow SSH and Ping from VPC A"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.vpc_a.cidr_block]
  }

  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = [aws_vpc.vpc_a.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}


# ---------------------
# EC2 Instance (Public)
# ---------------------
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_key_pair" "key" {
  key_name   = "transit-demo"
  public_key = file("~/.ssh/transit.pem.pub")
}

resource "aws_instance" "ec2_public" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.subnet_a_public.id
  vpc_security_group_ids = [aws_security_group.sg_public.id]
  key_name               = aws_key_pair.key.key_name
  associate_public_ip_address = true
  tags = { Name = "EC2-Public-A" }
}


resource "aws_instance" "ec2_private" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = "t3.micro"
  subnet_id              = aws_subnet.subnet_b_private.id
  vpc_security_group_ids = [aws_security_group.sg_b_private.id]
  key_name               = aws_key_pair.key.key_name
  associate_public_ip_address = false
  tags = { Name = "EC2-Private-B" }
}
