output "ec2_public_ip" {
  description = "Public IP of EC2 instance in VPC A"
  value       = aws_instance.ec2_public.public_ip
}

output "tgw_id" {
  description = "Transit Gateway ID"
  value       = aws_ec2_transit_gateway.tgw.id
}

output "vpc_a_id" {
  value = aws_vpc.vpc_a.id
}

output "vpc_b_id" {
  value = aws_vpc.vpc_b.id
}
