# ===========================
# ECR REPOSITORY
# ===========================
resource "aws_ecr_repository" "backend_repo" {
  name = "${var.project}-backend"
}

# ===========================
# BACKEND TASK DEFINITION
# ===========================
resource "aws_ecs_task_definition" "backend_task" {
  family                   = "${var.project}-backend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "256"
  memory                   = "512"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "backend",
      image     = "${var.aws_account_id}.dkr.ecr.${var.aws_region}.amazonaws.com/${var.project}-backend:latest",
      essential = true,
      portMappings = [
        {
          containerPort = 8080,
          hostPort      = 8080
        }
      ],
      environment = [
        {
          name  = "REDIS_HOST",
          value = aws_elasticache_cluster.redis.cache_nodes[0].address
        },
        {
          name  = "REDIS_PORT",
          value = "6379"
        }
      ]
    }
  ])
}

# ===========================
# ALB TARGET GROUP AND LISTENER RULE FOR BACKEND
# ===========================
resource "aws_lb_target_group" "backend_tg" {
  name     = "${var.project}-backend-tg"
  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.cv_c9_vpc.id
  target_type = "ip"
}

resource "aws_lb_listener_rule" "backend_rule" {
  listener_arn = aws_lb_listener.alb_listener.arn
  priority     = 200

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.backend_tg.arn
  }

  condition {
    path_pattern {
      values = ["/api*"]
    }
  }
}

# ===========================
# BACKEND ECS SERVICE
# ===========================
resource "aws_ecs_service" "backend_service" {
  name            = "${var.project}-backend-service"
  cluster         = aws_ecs_cluster.cv_c9_cluster.id
  task_definition = aws_ecs_task_definition.backend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.private_1.id, aws_subnet.private_2.id]
    security_groups = [aws_security_group.redis_sg.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.backend_tg.arn
    container_name   = "backend"
    container_port   = 8080
  }

  depends_on = [aws_lb_listener_rule.backend_rule]
}
