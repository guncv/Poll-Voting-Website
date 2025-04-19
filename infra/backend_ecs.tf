# ===========================
# SECURITY GROUP: BACKEND
# ===========================
resource "aws_security_group" "backend_sg" {
  name   = "${var.project}-backend-sg"
  vpc_id = aws_vpc.cv_c9_vpc.id


  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
    description = "Internal VPC access to backend service"
  }

  # Allow ECS task to pull from ECR (via VPC Endpoint on port 443)
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
    description = "HTTPS access for ECR and CloudWatch VPC Endpoints"
  }

  # Egress rule for allowing outbound HTTPS traffic (port 443)
  egress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]  # Allow outbound traffic on port 443 to any destination
    description = "Outbound HTTPS traffic (for ECR, CloudWatch, etc.)"
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-backend-sg"
  }
}

# ===========================
# TASK DEFINITION: BACKEND
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
        },
        {
          name  = "APP_ENV",
          value = var.env
        }
      ]
     logConfiguration = {
        logDriver = "awslogs"
        options = {
          awslogs-group         = "/ecs/${var.project}-backend"  
          awslogs-region        = var.aws_region                
          awslogs-stream-prefix = "ecs"                         
        }
      }
    }
  ])

  tags = {
    Name = "${var.project}-backend-task"
  }
}

# ===========================
# ECS SERVICE: BACKEND
# ===========================
resource "aws_ecs_service" "backend_service" {
  name            = "${var.project}-backend-service"
  cluster         = aws_ecs_cluster.cv_c9_cluster.id
  task_definition = aws_ecs_task_definition.backend_task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets         = [aws_subnet.private_1.id, aws_subnet.private_2.id]
    security_groups = [aws_security_group.backend_sg.id]
    assign_public_ip = false
  }

  tags = {
    Name = "${var.project}-backend-service"
  }
}
