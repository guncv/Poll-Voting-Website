
# ===========================
# Custom Policy to allow access to CloudWatch, VPC Endpoints, and ECR
# ===========================
resource "aws_cloudwatch_log_group" "ecs_log_group" {
  name = "/ecs/${var.project}-logs"

  tags = {
    Name = "${var.project}-ecs-log-group"
  }
}

resource "aws_iam_policy" "ecs_task_execution_policy_custom" {
  name        = "ecs-task-execution-policy-custom"
  description = "Custom policy to allow CloudWatch, ECR access, and VPC Endpoint usage"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      # Allow ECS Task to pull images from ECR
      {
        Effect = "Allow",
        Action = [
          "ecr:GetAuthorizationToken",
          "ecr:BatchGetImage",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchCheckLayerAvailability"
        ],
        Resource = "arn:aws:ecr:${var.aws_region}:${var.aws_account_id}:repository/*"
      },
      # Allow ECS Task to write logs to CloudWatch for frontend and backend log groups
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "logs:DescribeLogStreams",
          "logs:DescribeLogGroups",
          "logs:TagResource"
        ],
        Resource = [
          "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:/ecs/${var.project}-backend:*",
          "arn:aws:logs:${var.aws_region}:${var.aws_account_id}:log-group:/ecs/${var.project}-frontend:*"
        ]
      },
      # Allow ECS Task to connect to VPC Endpoint services
      {
        Effect = "Allow",
        Action = [
          "ec2:DescribeVpcEndpoints",
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSubnets"
        ],
        Resource = "*"
      }
    ]
  })
}

# ===========================
# IAM ROLE FOR ECS TASKS
# ===========================
resource "aws_iam_role" "ecs_task_execution_role" {
  name       = "${var.project}-ecs-task-execution-role"
  depends_on = [var.vpc_id]

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect    = "Allow",
        Action    = "sts:AssumeRole",
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

# ===========================
# Attach IAM Policies to ECS Task Role
# ===========================
resource "aws_iam_role_policy_attachment" "ecs_task_execution_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "ecs_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonECS_FullAccess"
}

resource "aws_iam_role_policy_attachment" "elasticache_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess"
}

resource "aws_iam_role_policy_attachment" "vpc_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonVPCFullAccess"
}

resource "aws_iam_role_policy_attachment" "cloudwatch_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"
}

# ✅ Attach **EC2 Full Access**
resource "aws_iam_role_policy_attachment" "ec2_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2FullAccess"
}

# ✅ Attach **S3 Full Access**
resource "aws_iam_role_policy_attachment" "s3_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}

# ✅ Add IAMFullAccess Policy to Create Role Permissions
resource "aws_iam_role_policy_attachment" "iam_full_access" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/IAMFullAccess"
}

# Attach the custom ECS task execution policy
resource "aws_iam_role_policy_attachment" "ecr_pull_policy_custom" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = aws_iam_policy.ecs_task_execution_policy_custom.arn
}

# ===========================
# IAM Instance Profile for ECS Tasks
# ===========================
resource "aws_iam_instance_profile" "ecs_task_instance_profile" {
  name = "${var.project}-ecs-task-instance-profile"
  role = aws_iam_role.ecs_task_execution_role.name
}
