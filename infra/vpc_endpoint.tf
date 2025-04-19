# ===========================
# VPC ENDPOINTS for ECR & CloudWatch Logs
# ===========================
resource "aws_vpc_endpoint" "ecr_api" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.api"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.backend_sg.id]

  tags = {
    Name = "${var.project}-ecr-api-endpoint"
  }
}

resource "aws_vpc_endpoint" "ecr_dkr" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.ecr.dkr"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.backend_sg.id]

  tags = {
    Name = "${var.project}-ecr-dkr-endpoint"
  }
}

resource "aws_vpc_endpoint" "logs" {
  vpc_id            = aws_vpc.cv_c9_vpc.id
  service_name      = "com.amazonaws.${var.aws_region}.logs"
  vpc_endpoint_type = "Interface"
  subnet_ids        = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_group_ids = [aws_security_group.backend_sg.id]

  tags = {
    Name = "${var.project}-logs-endpoint"
  }
}
