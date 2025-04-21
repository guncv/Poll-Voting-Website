
resource "aws_cloudwatch_log_group" "backend_log_group" {
  name = "/ecs/${var.project}-backend"

  tags = {
    Name = "${var.project}-backend-log-group"
  }
}

resource "aws_cloudwatch_log_group" "frontend_log_group" {
  name = "/ecs/${var.project}-frontend"

  tags = {
    Name = "${var.project}-frontend-log-group"
  }
}

