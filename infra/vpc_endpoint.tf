# ===========================
# VPC ENDPOINTS for ECR & CloudWatch Logs
# ===========================
resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.vpc_endpoint_sg.id] # Use vpc_endpoint_sg here

  tags = {
    Name = "${var.project}-ecr-api-endpoint"
  }
}

resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.dkr"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.vpc_endpoint_sg.id] # Use vpc_endpoint_sg here

  tags = {
    Name = "${var.project}-ecr-dkr-endpoint"
  }
}

resource "aws_vpc_endpoint" "logs" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.logs"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.vpc_endpoint_sg.id] # Use vpc_endpoint_sg here

  tags = {
    Name = "${var.project}-logs-endpoint"
  }
}

resource "aws_security_group" "vpc_endpoint_sg" {
  name        = "${var.project}-vpc-endpoint-sg"
  description = "Allow traffic to/from VPC endpoints"
  vpc_id      = aws_vpc.cv_c9_vpc.id

  # Allow inbound HTTPS traffic for the VPC endpoint from ECS tasks
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"] # Internal VPC access
  }

  # Allow outbound traffic for services to reach CloudWatch and ECR
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]  # Allow all outbound traffic (you may narrow this further if needed)
  }

  tags = {
    Name = "${var.project}-vpc-endpoint-sg"
  }
}

