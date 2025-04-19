# ===========================
# ECS CLUSTER
# ===========================
resource "aws_ecs_cluster" "cv_c9_cluster" {
  name = "${var.project}-cluster"
}

# ===========================
# IAM ROLE FOR ECS TASKS
# ===========================
resource "aws_iam_role" "ecs_task_execution_role" {
  name               = "${var.project}-ecs-task-execution-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect   = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
        Action   = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Name = "${var.project}-ecs-task-execution-role"
  }
}

# Custom Policy to allow access to CloudWatch, VPC Endpoints, and ECR
resource "aws_iam_policy" "ecs_task_execution_policy_custom" {
  name        = "ecs-task-execution-policy-custom"
  description = "Custom policy to allow CloudWatch, ECR access, and VPC Endpoint usage"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      # Allow ECS Task to pull images from ECR
      {
        Effect   = "Allow",
        Action   = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchGetImage",
          "ecr:BatchCheckLayerAvailability"
        ],
        Resource = "arn:aws:ecr:${var.aws_region}:${var.aws_account_id}:repository/*"
      },
      # Allow ECS Task to write logs to CloudWatch
      {
        Effect   = "Allow",
        Action   = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams",
          "logs:DescribeLogGroups"
        ],
        Resource = "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:/ecs/*"
      },
      # Allow ECS Task to connect to VPC Endpoint services
      {
        Effect   = "Allow",
        Action   = [
          "ec2:DescribeVpcEndpoints",
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSubnets"
        ],
        Resource = "*"
      }
    ]
  })
}

# Attach the custom policy to the ECS Task Execution Role
resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy_custom_attachment" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = aws_iam_policy.ecs_task_execution_policy_custom.arn 
}

# ===========================
# SHARED APPLICATION LOAD BALANCER
# ===========================
resource "aws_lb" "ecs_alb" {
  name               = "${var.project}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = [aws_subnet.public_1.id,aws_subnet.public_2.id]

  tags = {
    Name = "${var.project}-alb"
  }
}

resource "aws_security_group" "alb_sg" {
  name        = "${var.project}-alb-sg"
  description = "Allow HTTP traffic to ALB"
  vpc_id      = aws_vpc.cv_c9_vpc.id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.project}-alb-sg"
  }
}


resource "aws_lb_listener" "alb_listener" {
  load_balancer_arn = aws_lb.ecs_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "text/plain"
      message_body = "404 Not Found"
      status_code  = "404"
    }
  }
}