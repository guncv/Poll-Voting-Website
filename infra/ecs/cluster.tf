# ===========================
# ECS CLUSTER
# ===========================
resource "aws_ecs_cluster" "cv_c9_cluster" {
  name = "${var.project}-cluster"
}