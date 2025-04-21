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

    network_configuration {
    subnets          = var.private_subnet_ids
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
    var.vpc_id,
  ]
}