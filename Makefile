TERRAFORM_DIR := infra
AWS_ACCOUNT_ID := 913524943390
ALB_DNS_NAME := $(shell terraform -chdir=$(TERRAFORM_DIR) output -raw alb_dns_name)
AWS_REGION := us-west-2


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

# Step 1: Refresh Terraform state and output values
refresh-outputs:
	cd $(TERRAFORM_DIR) && terraform refresh

# Step 2: Deploy Network (VPC, Subnets, Internet Gateway, Route Tables)
deploy-network:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_vpc.cv_c9_vpc \
		-target=aws_subnet.public \
		-target=aws_subnet.public_2 \
		-target=aws_subnet.public_3 \
		-target=aws_subnet.private_1 \
		-target=aws_subnet.private_2 \
		-target=aws_subnet.private_3 \
		-target=aws_internet_gateway.igw \
		-target=aws_route_table.public_rt \
		-target=aws_route_table.private_rt \
		-target=aws_route_table_association.public_assoc \
		-target=aws_route_table_association.private_assoc_1 \
		-target=aws_route_table_association.private_assoc_2 \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve

# Step 3: Deploy Redis infrastructure (if needed)
deploy-redis:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_elasticache_cluster.redis \
		-target=aws_vpc.cv_c9_vpc \
		-target=aws_subnet.private_1 \
		-target=aws_subnet.private_2 \
		-target=aws_subnet.private_3 \
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

# Step 4: Deploy ECR repositories (if not created)
create-ecr:
	aws ecr describe-repositories --repository-name cv-c9-backend --region $(AWS_REGION) || aws ecr create-repository --repository-name cv-c9-backend --region $(AWS_REGION)
	aws ecr describe-repositories --repository-name cv-c9-frontend --region $(AWS_REGION) || aws ecr create-repository --repository-name cv-c9-frontend --region $(AWS_REGION)

deploy-backend: deploy-ecr-login build-backend push-backend

deploy-frontend: deploy-ecr-login build-frontend push-frontend

# Step 5: Deploy ECR login (for Docker login)
deploy-ecr-login:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com

# Step 6: Build Backend Docker Image for ECS
build-backend:
	docker buildx build --platform linux/amd64 -t cv-c9-backend ./backend

# Step 7: Push Backend Docker Image to ECR (with correct platform)
push-backend: 
	docker buildx build --platform linux/amd64 -t $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-backend:latest --push ./backend

# Step 8: Build Frontend Docker Image for ECS
build-frontend:
	docker buildx build --platform linux/amd64 -t cv-c9-frontend ./frontend

# Step 9: Push Frontend Docker Image to ECR (with correct platform)
push-frontend: 
	docker buildx build --platform linux/amd64 -t $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/cv-c9-frontend:latest --push ./frontend

update-frontend-path:
	echo "Fetching ALB DNS name..."
	ALB_DNS_NAME=$(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)
	echo "Updating frontend API path with ALB DNS..."
	sed -i '' "s|NEXT_PUBLIC_API_PATH=http://.*|NEXT_PUBLIC_API_PATH=http://$(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)/api|g" ./frontend/.env.local
	echo "Frontend API path updated with ALB DNS: $(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)"

update-cors-domain:
	echo "Updating CORS domain..."
	sed -i '' "s|CORS_ECS_DOMAIN=.*|CORS_ECS_DOMAIN=http://$(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)|g" ./backend/.env
	sed -i '' "s|CORS_ECS_DOMAIN=.*|CORS_ECS_DOMAIN=http://$(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)|g" ./backend/.env.prod
	echo "CORS domain updated with ALB DNS: $(shell aws elbv2 describe-load-balancers --names cv-c9-alb --query "LoadBalancers[0].DNSName" --output text)"

update-ecr: \
	deploy-ecr-login \
	push-backend \
	push-frontend

# Step 10: Deploy ECS Services (Combined Backend and Frontend)
deploy-ecs:
	cd $(TERRAFORM_DIR) && terraform apply \
		-target=aws_ecs_cluster.cv_c9_cluster \
		-target=aws_ecs_task_definition.combined_task \
		-target=aws_ecs_service.combined_service \
		-var-file="terraform.tfvars" \
		-var-file="private.tfvars" \
		-auto-approve \
		-refresh=true

restart-ecs:
	@echo "Updating ECS service and forcing a new deployment..."
	@aws ecs update-service \
		--cluster cv-c9-cluster \
		--service cv-c9-combined-service \
		--force-new-deployment \
		--no-cli-pager > /dev/null 2>&1
	@echo "ECS service has been updated and new containers are running."

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

tf-refresh:
	cd $(TERRAFORM_DIR) && terraform refresh

reset-repo:
	@echo "Resetting ECR repositories..."
	# Delete only if it already exists
	@aws ecr describe-repositories \
	    --repository-names cv-c9-backend \
	    --region $(AWS_REGION) >/dev/null 2>&1 \
	  && aws ecr delete-repository \
	       --repository-name cv-c9-backend \
	       --region $(AWS_REGION) \
	       --force --no-cli-pager > /dev/null 2>&1
	@aws ecr describe-repositories \
	    --repository-names cv-c9-frontend \
	    --region $(AWS_REGION) >/dev/null 2>&1 \
	  && aws ecr delete-repository \
	       --repository-name cv-c9-frontend \
	       --region $(AWS_REGION) \
	       --force --no-cli-pager > /dev/null 2>&1

	# Now (re)create them—ignore "already exists" on create
	@aws ecr create-repository \
	    --repository-name cv-c9-backend \
	    --region $(AWS_REGION) >/dev/null 2>&1 || true
	@aws ecr create-repository \
	    --repository-name cv-c9-frontend \
	    --region $(AWS_REGION) >/dev/null 2>&1 || true

	@echo "ECR repositories reset successfully."


deploy-all: \
	env-prod \
	tf-init \
	tf-apply \
	tf-output \
	reset-repo \
	update-frontend-path \
	update-cors-domain \
	update-ecr \
	restart-ecs

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
	@echo "   deploy-network  - Deploy network infra (VPC, Subnets)"
	@echo "   build-backend    - Build backend Docker image (linux/amd64)"
	@echo "   push-backend     - Build & Push backend image to ECR (linux/amd64)"
	@echo "   build-frontend   - Build frontend Docker image (linux/amd64)"
	@echo "   push-frontend    - Build & Push frontend image to ECR (linux/amd64)"

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
