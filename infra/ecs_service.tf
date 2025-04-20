# ===========================
# ECS SERVICE (Combining Backend and Frontend)
# ===========================
resource "aws_ecs_service" "combined_service" {
  name                 = "${var.project}-combined-service"
  cluster              = aws_ecs_cluster.cv_c9_cluster.id
  task_definition      = aws_ecs_task_definition.combined_task.arn
  desired_count        = 1
  launch_type          = "FARGATE"
  force_new_deployment = true


  # Network Configuration
#   network_configuration {
#     subnets          = [aws_subnet.public.id, aws_subnet.public_2.id, aws_subnet.public_3.id]
#     security_groups  = [aws_security_group.ecs_sg.id]
#     assign_public_ip = true
#   }
network_configuration {
  subnets          = [aws_subnet.private_1.id, aws_subnet.private_2.id]
  security_groups  = [aws_security_group.ecs_sg.id]
  assign_public_ip = false
}



  # Load Balancer configuration for Frontend
  load_balancer {
    target_group_arn = aws_lb_target_group.frontend_tg.arn
    container_name   = "frontend"
    container_port   = 3000
  }

  # Load Balancer configuration for Backend
  load_balancer {
    target_group_arn = aws_lb_target_group.backend_tg.arn
    container_name   = "backend"
    container_port   = 8080
  }

  depends_on = [
    aws_lb_listener_rule.frontend_rule,
    aws_lb_listener_rule.backend_rule,
    aws_nat_gateway.nat,
    aws_route_table.private_rt,
    aws_route_table.public_rt,
    aws_vpc_endpoint.ecr_api,
    aws_vpc_endpoint.ecr_dkr,
    aws_vpc_endpoint.s3_gateway
  ]
}

# ===========================
# ALB TARGET GROUP FOR FRONTEND
# ===========================
resource "aws_lb_target_group" "frontend_tg" {
  name        = "${var.project}-frontend-tg"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.cv_c9_vpc.id
  target_type = "ip"
}

# ===========================
# ALB TARGET GROUP FOR BACKEND
# ===========================
resource "aws_lb_target_group" "backend_tg" {
  name        = "${var.project}-backend-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = aws_vpc.cv_c9_vpc.id
  target_type = "ip"
}

# ===========================
# ALB LISTENER RULES
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

# ===========================
# ECS TASK DEFINITION (Multi-container for Backend and Frontend)
# ===========================
resource "aws_ecs_task_definition" "combined_task" {
  family                   = "${var.project}-combined"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_execution_role.arn

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
      ],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-group"         = "/ecs/${var.project}-backend",
          "awslogs-region"        = var.aws_region,
          "awslogs-stream-prefix" = "backend"
        }
      }
    },
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
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          "awslogs-group"         = "/ecs/${var.project}-frontend",
          "awslogs-region"        = var.aws_region,
          "awslogs-stream-prefix" = "frontend"
        }
      }
    }
  ])
}
