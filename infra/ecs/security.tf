# ===========================
# ECS Security Group
# ===========================
resource "aws_security_group" "ecs_sg" {
  name        = "${var.project}-ecs-sg"
  description = "Allow traffic to ECS service"
  vpc_id      = var.vpc_id

  # ALB access to frontend
  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow frontend (3000) from anywhere"
  }

  # ALB access to backend
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow backend (8080) from anywhere"
  }

  # Default HTTP and HTTPS
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow HTTP traffic"
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow HTTPS traffic"
  }

  # Outbound (required for ECS, ECR, Redis, etc.)
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-ecs-sg"
  }

  depends_on = [var.vpc_id]
}

