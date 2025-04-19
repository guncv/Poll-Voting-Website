# ===========================
# SECURITY GROUP: FRONTEND
# ===========================
resource "aws_security_group" "frontend_sg" {
  name        = "${var.project}-frontend-sg"
  description = "Allow HTTP access to frontend"
  vpc_id      = aws_vpc.cv_c9_vpc.id

  ingress {
    from_port   = 3000
    to_port     = 3000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Optional: Restrict this in production
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-frontend-sg"
  }
}

# ===========================
# TASK DEFINITION: FRONTEND
# ===========================
resource "aws_ecs_task_definition" "frontend_task" {
  family                   = "${var.project}-frontend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "frontend",
      image     = "${var.aws_account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.project}-frontend:latest",
      essential = true,
      portMappings = [
        {
          containerPort = 3000,
          hostPort      = 3000
        }
      ],
      environment = [
        {
          name  = "BACKEND_URL",
          value = "http://${var.project}-backend-service.local:8080" # optional
        }
      ]
    }
  ])
}

# ===========================
# ALB TARGET GROUP: FRONTEND
# ===========================
resource "aws_lb_target_group" "frontend_tg" {
  name        = "${var.project}-frontend-tg"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.cv_c9_vpc.id
  target_type = "ip"

  health_check {
    path                = "/"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
    matcher             = "200"
  }

  tags = {
    Name = "${var.project}-frontend-tg"
  }
}

# ===========================
# ALB LISTENER RULE: FRONTEND
# ===========================
resource "aws_lb_listener_rule" "frontend_rule" {
  listener_arn = aws_lb_listener.alb_listener.arn
  priority     = 100

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.frontend_tg.arn
  }

  condition {
    path_pattern {
      values = ["/"]
    }
  }
}

# ===========================
# ECS SERVICE: FRONTEND
# ===========================
resource "aws_ecs_service" "frontend_service" {
  name            = "${var.project}-frontend-service"
  cluster         = aws_ecs_cluster.cv_c9_cluster.id
  task_definition = aws_ecs_task_definition.frontend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.public_1.id, aws_subnet.public_2.id]
    security_groups = [aws_security_group.frontend_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.frontend_tg.arn
    container_name   = "frontend"
    container_port   = 3000
  }

  depends_on = [aws_lb_listener_rule.frontend_rule]

  tags = {
    Name = "${var.project}-frontend-service"
  }
}
