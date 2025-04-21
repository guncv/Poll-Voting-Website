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
          value = var.redis_endpoint
        },
        {
          name  = "REDIS_PORT",
          value = tostring(var.redis_port)
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
      environment = [
        {
          name  = "API_PATH",
          value = "http://${aws_lb.ecs_alb.dns_name}/api"
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
