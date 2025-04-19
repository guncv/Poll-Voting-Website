TERRAFORM_DIR := infra
AWS_ACCOUNT_ID := $(shell terraform -chdir=$(TERRAFORM_DIR) output -raw aws_account_id)
AWS_REGION := ap-southeast-1

.PHONY: info run run-prod build build-prod down clean logs restart ps rebuild rebuild-prod \
        tf-init tf-plan tf-apply tf-destroy tf-output help \
        deploy-redis build-backend push-backend build-frontend push-frontend deploy-ecs

# =======================
# Docker (Dev)
# =======================

info:
	docker compose ps

run:
	docker compose up

down:
	docker compose down

build:
	docker compose build

clean:
	docker compose down --rmi all --volumes --remove-orphans

logs:
	docker compose logs -f

restart:
	docker compose restart

ps:
	docker compose ps

rebuild: clean build run

# =======================
# Docker (Prod)
# =======================

run-prod:
	docker compose -f docker-compose.yml up

build-prod:
	docker compose -f docker-compose.yml build

rebuild-prod:
	docker compose -f docker-compose.yml down --remove-orphans && \
	docker compose -f docker-compose.yml build && \
	docker compose -f docker-compose.yml up

# =======================
# Terraform Commands
# =======================

TERRAFORM_DIR := infra

# Step 1: Deploy Redis infrastructure (ElastiCache)
deploy-redis: _apply-redis update-redis-host env-prod

_apply-redis:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=data.aws_caller_identity.current \
		-target=aws_elasticache_cluster.redis \
		-target=aws_vpc.cv_c9_vpc \
		-target=aws_subnet.private_1 \
		-target=aws_subnet.private_2 \
		-target=aws_subnet.public \
		-target=aws_internet_gateway.igw \
		-target=aws_nat_gateway.nat \
		-target=aws_route_table.public_rt \
		-target=aws_route_table.private_rt \
		-target=aws_security_group.redis_sg \
		-target=aws_elasticache_subnet_group.redis_subnet_group \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

update-redis-host:
	@echo "🔄 Updating REDIS_HOST in backend/.env.prod..."
	@sed -i '' "s/^REDIS_HOST=.*/REDIS_HOST=$(shell terraform -chdir=$(TERRAFORM_DIR) output -raw redis_endpoint)/" backend/.env.prod
	@echo "✅ REDIS_HOST updated to latest Redis endpoint."
	
# Step 2: Push Backend Docker Images
check-aws:
	@which aws >/dev/null || (echo "❌ AWS CLI not found. Please install it first." && exit 1)

create-ecr:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_ecr_repository.backend_repo \
		-target=aws_ecr_repository.frontend_repo \
		-target=aws_vpc_endpoint.ecr_api \
		-target=aws_vpc_endpoint.ecr_dkr \
		-target=aws_vpc_endpoint.logs \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

deploy-ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

deploy-backend : env-prod build-backend push-backend deploy-backend-ecs 

build-backend:
	docker build -t cv-c9-backend ./backend

push-backend: deploy-ecr-login
	docker tag cv-c9-backend $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-backend:latest
	docker push $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-backend:latest

# Step 3: Deploy Backend ECS Services
deploy-backend-ecs:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_ecs_cluster.cv_c9_cluster \
		-target=aws_iam_role.ecs_task_execution_role \
		-target=aws_iam_role_policy_attachment.ecs_task_execution_policy \
		-target=aws_ecs_task_definition.backend_task \
		-target=aws_ecs_service.backend_service \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

# Step 4: Push Frontend Docker Images
build-frontend:
	docker build -t cv-c9-frontend ./frontend

push-frontend: deploy-ecr-login
	docker tag cv-c9-frontend $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-frontend:latest
	docker push $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-frontend:latest
	
# Step 5: Deploy Frontend ECS Services
deploy-frontend-ecs:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_ecs_task_definition.frontend_task \
		-target=aws_ecs_service.frontend_service \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

# Full infra control
tf-init:
	cd $(TERRAFORM_DIR) && terraform init

tf-plan:
	cd $(TERRAFORM_DIR) && terraform plan \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars"

tf-apply:
	cd $(TERRAFORM_DIR) && terraform apply \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

tf-destroy:
	cd $(TERRAFORM_DIR) && terraform destroy \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

tf-output:
	cd $(TERRAFORM_DIR) && terraform output

# =======================
# Environment Switching
# =======================

env-dev:
	@echo "🔄 Switching to LOCAL .env..."
	@cp backend/.env.local backend/.env
	@echo "✅ Now using backend/.env.local → backend/.env"

env-prod:
	@echo "🔄 Switching to PRODUCTION .env..."
	@cp backend/.env.prod backend/.env
	@echo "✅ Now using backend/.env.prod → backend/.env"

# =======================
# Help Guide
# =======================

help:
	@echo ""
	@echo "🛠  Available Makefile Commands:"
	@echo "-------------------------------"
	@echo " Docker (Dev):"
	@echo "   run           - Start containers (dev mode)"
	@echo "   build         - Build all images (dev)"
	@echo "   down          - Stop containers"
	@echo "   clean         - Remove containers, images, volumes"
	@echo "   logs          - Show logs"
	@echo "   restart       - Restart containers"
	@echo "   ps            - Show container status"
	@echo "   rebuild       - Clean + Build + Run (dev)"
	@echo ""
	@echo " Docker (Prod):"
	@echo "   run-prod      - Start containers (prod only)"
	@echo "   build-prod    - Build images (prod)"
	@echo "   rebuild-prod  - Clean + Build + Run (prod)"
	@echo ""
	@echo " Terraform (Modular):"
	@echo "   deploy-redis     - Deploy Redis infra (ElastiCache, VPC)"
	@echo "   build-backend    - Build Docker backend image"
	@echo "   push-backend     - Push backend image to ECR"
	@echo "   build-frontend   - Build frontend Docker image"
	@echo "   push-frontend    - Push frontend image to ECR"
	@echo "   deploy-ecs       - Deploy ECS Cluster + Tasks"
	@echo ""
	@echo " Terraform (Full):"
	@echo "   tf-init       - Initialize Terraform"
	@echo "   tf-plan       - Show execution plan"
	@echo "   tf-apply      - Apply infrastructure (auto approve)"
	@echo "   tf-destroy    - Destroy infrastructure (auto approve)"
	@echo "   tf-output     - Show Terraform outputs"
	@echo ""
	@echo " Environment Switching:"
	@echo "   env-dev       - Use development .env"
	@echo "   env-prod      - Use production .env"
	@echo ""
	@echo " Pro Tip 💡: Use 'make help' to rediscover commands anytime!"
