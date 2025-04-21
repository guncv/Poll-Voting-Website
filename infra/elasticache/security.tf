// ========================================
// Security Group for Redis
// ========================================

resource "aws_security_group" "redis_sg" {
  name        = "${var.project}-redis-sg"
  description = "Allow Redis traffic"
  vpc_id      = var.vpc_id

  # 🧠 Conditionally allow Redis access based on environment
  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = var.env == "dev" ? ["0.0.0.0/0"] : ["10.0.0.0/16"]
    description = var.env == "dev" ? "Allow from anywhere (dev)" : "Internal VPC access only (prod)"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-redis-sg"
  }

  depends_on = [var.vpc_id]
}