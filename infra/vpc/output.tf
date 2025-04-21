output "vpc_id" {
  value = aws_vpc.cv_c9_vpc.id
}

output "public_subnet_ids" {
  value = [aws_subnet.public.id, aws_subnet.public_2.id, aws_subnet.public_3.id]
}

output "private_subnet_ids" {
  value = [aws_subnet.private_1.id, aws_subnet.private_2.id]
}

output "public_route_table_id" {
  value = aws_route_table.public_rt.id
}

output "private_route_table_id" {
  value = aws_route_table.private_rt.id
}
