# 🗳️ Poll-Voting-Website

A full-stack web application where users vote on a daily "hot" question with two simple, fun, or thought-provoking answer choices. Designed to be lightweight, interactive, and engaging.

---

## 🛠️ Tech Stack

- **Frontend**: Next.js
- **Backend**: Golang (Gin)
- **Containerization**: Docker
- **Infrastructure as Code**: Terraform
- **Deployment**: AWS ECS (Fargate)

---

## ☁️ AWS Services Used

- **Amazon ECS (Fargate)** – Hosts the frontend and backend in containers with serverless scalability.
- **Amazon ElastiCache (Redis)** – Caches real-time voting data with a 24-hour TTL and auto-reset.
- **Amazon SNS** – Sends notifications for milestones and trending questions.
- **Amazon RDS (PostgreSQL)** – Stores user data, metadata, and question content.

---

## 🚀 Deployment Guide (AWS - Terraform + ECS + ECR + Redis)

This guide explains how to deploy the infrastructure using **Terraform**, **Docker**, and **AWS**.

---

### ⚙️ Prerequisites

- ✅ Terraform installed
- ✅ Docker installed and running
- ✅ AWS CLI configured (`aws configure`)
- ✅ Access to ECR, ECS, and ElastiCache via IAM
- ✅ `.env.prod` prepared for `backend` and `frontend`

---

### 📦 Project Structure

```
.
├── backend/                   # Go server
├── frontend/                  # Next.js app
├── infra/                     # Terraform configurations
│   ├── ecs/                   # ECS-related resources (ALB, cluster, tasks, etc.)
│   │   ├── alb.tf
│   │   ├── alb_target_group.tf
│   │   ├── cloud_logs.tf
│   │   ├── cluster.tf
│   │   ├── iam.tf
│   │   ├── output.tf
│   │   ├── security.tf
│   │   ├── service.tf
│   │   ├── task_definition.tf
│   │   ├── variable.tf
│   │   └── vpc_endpoint.tf
│   ├── elasticache/           # ElastiCache Redis config
│   │   ├── cluster.tf
│   │   ├── output.tf
│   │   ├── security.tf
│   │   ├── subnet_group.tf
│   │   └── variable.tf
│   ├── vpc/                   # VPC networking setup
│   │   ├── internet_gateway.tf
│   │   ├── nat.tf
│   │   ├── output.tf
│   │   ├── route_table.tf
│   │   ├── security.tf
│   │   ├── subnet.tf
│   │   ├── variable.tf
│   │   └── vpc.tf
│   ├── main.tf
│   ├── outputs.tf
│   ├── provider.tf
│   ├── terraform.tfvars       # Shared variables
│   └── variables.tf
├── docker-compose.yml         # Docker multi-service config
├── docker-compose.override.yml
├── .dockerignore
├── .gitignore
├── Makefile                   # Deployment and utility commands
└── README.md
```

---

### 🔐 Setting Up `private.tfvars`

Create a file called `private.tfvars` inside the `/infra` directory:

```bash
touch infra/private.tfvars
```

Then add the following contents (replace with your actual values):

```hcl
aws_access_key = "YOUR_AWS_ACCESS_KEY"
aws_secret_key = "YOUR_AWS_SECRET_KEY"
aws_account_id = "YOUR_AWS_ACCOUNT_ID"
```
Sure! Based on your Makefile, here's a breakdown of the `make deploy-all` and its sub-methods that are executed, which you can include in your README file. The explanation will help users understand how each step in the Makefile contributes to the full deployment process.

---

### 🛠 Step-by-Step Explanation of `make deploy-all`

The `make deploy-all` command is a combination of various steps to initialize, deploy, and configure the entire infrastructure and application stack. Below is a detailed explanation of what happens when you run `make deploy-all`:

#### 1. 🧑‍💻 **Switch to Production Environment**
   ```bash
   make env-prod
   ```
   - This step ensures that the system is using the production environment variables. It copies the production `.env.prod` file to `.env`, which is used by both the backend and frontend services to connect to the appropriate services and databases in the production environment.

#### 2. ⚙️ **Initialize Terraform**
   ```bash
   make tf-init
   ```
   - This command initializes Terraform and prepares it to manage the infrastructure. It sets up the backend state and makes sure that all dependencies are in place for applying infrastructure configurations.

#### 3. 🚀 **Apply Terraform Configurations**
   ```bash
   make tf-apply
   ```
   - This step applies the Terraform configuration to provision all necessary AWS infrastructure resources, such as the VPC, subnets, security groups, ECS, ECR, and ElastiCache. It sets up the environment as per your specifications in the `terraform.tfvars` and `private.tfvars` files.

#### 4. 🗂 **Display Terraform Outputs**
   ```bash
   make tf-output
   ```
   - Once the infrastructure is applied, this command fetches the output values from Terraform, such as DNS names for the ALB (Application Load Balancer) and Redis endpoints. These outputs are used in the next steps to configure the backend and frontend services.

#### 5. 🔄 **Reset ECR Repositories (Optional)**
   ```bash
   make reset-repo
   ```
   - This command ensures that the necessary Elastic Container Registry (ECR) repositories are set up for the backend and frontend images. If the repositories already exist, they are deleted and recreated to ensure a clean slate for the deployment.

#### 6. 🌍 **Update Frontend API Path**
   ```bash
   make update-frontend-path
   ```
   - This step updates the frontend configuration with the ALB's DNS name, ensuring that the frontend knows where to send API requests. It modifies the `.env.local` file in the frontend directory to reference the correct API base URL.

#### 7. 🔐 **Update CORS Domain**
   ```bash
   make update-cors-domain
   ```
   - This step ensures that the backend allows CORS (Cross-Origin Resource Sharing) requests from the frontend, which is necessary for handling requests between the backend and frontend hosted on different domains or subdomains.

#### 8. 🚢 **Build & Push Backend and Frontend Docker Images**
   ```bash
   make update-ecr
   ```
   - This combines the backend and frontend image build and push steps. It first logs into AWS ECR, builds Docker images for both services, and then pushes them to the correct ECR repositories for backend and frontend containers. These images are used in the ECS services.

#### 9. 🔄 **Restart ECS Services**
   ```bash
   make restart-ecs
   ```
   - This forces a new deployment of the ECS services, updating them with the latest Docker images and configurations. It's used to ensure the backend and frontend services run with the most recent changes.

---

### 🔄 Redeploying Specific Services

If you only need to redeploy one service, you can run the following commands:

- **Redeploy Backend:**
   ```bash
   make build-backend && make push-backend && make deploy-ecs
   ```
   This will rebuild the backend Docker image, push it to ECR, and redeploy it to ECS.

- **Redeploy Frontend:**
   ```bash
   make build-frontend && make push-frontend && make deploy-ecs
   ```
   This will rebuild the frontend Docker image, push it to ECR, and redeploy it to ECS.

---

### 🧹 Cleanup (Destroy All Infrastructure)

To tear down everything (ECS, Redis, VPC, and related AWS infrastructure), you can run:

```bash
make tf-destroy
```

This command will safely remove all the resources provisioned by Terraform, including ECS, Redis, and networking components.

---

### 🧑‍💼 Summary of Key `Makefile` Commands

| Command                | Description                                                    |
|------------------------|----------------------------------------------------------------|
| `make deploy-all`       | Full deployment of the entire stack, including backend, frontend, and infrastructure. |
| `make tf-init`          | Initializes Terraform and sets up backend state.               |
| `make tf-apply`         | Applies Terraform configurations to create all infrastructure resources. |
| `make tf-output`        | Fetches output values like ALB DNS and Redis endpoint.        |
| `make update-frontend-path` | Updates the frontend with the ALB DNS name for API communication. |
| `make update-cors-domain` | Configures CORS settings for the backend to accept frontend requests. |
| `make update-ecr`       | Builds and pushes backend/ frontend Docker images to ECR.     |
| `make restart-ecs`      | Forces a new deployment of the ECS services.                  |
| `make tf-destroy`       | Destroys all infrastructure created by Terraform.             |
