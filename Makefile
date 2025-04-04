.PHONY: info run run-prod build build-prod down clean logs restart ps rebuild rebuild-prod \
        tf-init tf-plan tf-apply tf-destroy tf-output help

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
	@echo " Terraform:"
	@echo "   tf-init       - Initialize Terraform"
	@echo "   tf-plan       - Show execution plan"
	@echo "   tf-apply      - Apply infrastructure (auto approve)"
	@echo "   tf-destroy    - Destroy infrastructure (auto approve)"
	@echo "   tf-output     - Show Terraform outputs"
	@echo ""
