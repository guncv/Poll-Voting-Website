# ===========================
# ECR REPOSITORY
# ===========================
resource "aws_ecr_repository" "backend_repo" {
    name = "${var.project}-backend"
    force_delete = true
}

# ===========================
# ECR REPOSITORY
# ===========================
resource "aws_ecr_repository" "frontend_repo" {
    name = "${var.project}-frontend"
    force_delete = true
}