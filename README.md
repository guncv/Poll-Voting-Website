# 🗳️ Poll-Voting-Website

A full-stack web application allowing users to vote on a daily "hot" question, featuring two simple, fun, or thought-provoking answer choices. Designed to be lightweight, interactive, and engaging.

**Created by CV On Cloud9**  
**Cloud Computing Final Project (2110524)**  

- จิรวัฒน์ เล่งน้อย (6431307421)  
- ชนกันต์ วิริยะสถาพรพงศ์ (6431309721)  
- ณัฐศิษฐ์ วิริยะโยธิน (6431318321)

---

## 🛠️ Tech Stack

- **Frontend**: Next.js + Typescript  
- **Backend**: Golang (Fiber)
- **Containerization**: Docker
- **Infrastructure as Code**: Terraform
- **Deployment**: AWS ECS (Fargate)

---

## ☁️ AWS Services Used

| Service                        | Description                                                  |
|--------------------------------|--------------------------------------------------------------|
| **Amazon ECS (Fargate)**       | Hosts the frontend and backend with serverless scalability.  |
| **Amazon ElastiCache (Redis)** | Real-time caching for voting data (24-hour TTL, auto-reset). |
| **Amazon SNS**                 | Sends notifications for user and admin milestones.           |
| **Amazon RDS (PostgreSQL)**    | Stores user data, metadata, and question content.            |
| **Amazon ECR**                 | Stores Docker images for backend and frontend.               |

---

## 🔐 IAM User and Policies Setup

### Terraform IAM User Policies (Full Access):

To provision AWS resources via Terraform, create an IAM user with these policies:

- `AmazonEC2ContainerRegistryFullAccess`
- `AmazonEC2FullAccess`
- `AmazonECS_FullAccess`
- `AmazonElastiCacheFullAccess`
- `AmazonS3FullAccess`
- `CloudWatchLogsFullAccess`
- `ElasticLoadBalancingFullAccess`
- `IAMFullAccess`
- `AmazonSNSFullAccess`

After creating this IAM user:

- Configure the AWS CLI locally with the IAM user's **ACCESS KEY** and **SECRET KEY**.
- Copy these keys into `infra/private.tfvars`.

### SNS IAM User (Limited):

Create a separate IAM user with minimal SNS permissions (**Publish only**).  
- Retrieve its **ACCESS KEY** and **SECRET KEY**.
- Store these credentials in the backend `.env` file (`SNS_ACCESS_KEY`, `SNS_SECRET_KEY`).
- this user is only used within the app for sending notifications, not for deploying anything.

---

## 🚀 Deployment Guide (AWS - Terraform + ECS + ECR + Redis)

This guide explains how to deploy the infrastructure using **Terraform**, **Docker**, and **AWS**.

---

### ⚙️ Prerequisites

- ✅ [Terraform](https://developer.hashicorp.com/terraform/downloads)
- ✅ [Docker](https://docs.docker.com/get-docker/)
- ✅ [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)
- ✅ `make` utility (`brew install make` on MacOS or `apt install make` on Ubuntu)
- ✅ Access to ECR, ECS, and ElastiCache via IAM
- ✅ `.env.local` prepared for `backend` and `frontend`
- ✅ `.env.prod` prepared for `backend` 

---

### 📦 Project Structure

```
.
├── backend/                     # Backend Golang server (Fiber)
│   ├── .env                     # Development/production config
│   ├── .env.prod                # Production environment config
│   └── Dockerfile               # used in production deployment
│   └── Dockerfile.dev           # for development
├── frontend/                    # Frontend Next.js application
│   ├── .env                     # Default environment variables
│   ├── .env.local               # Local overrides for development
│   └── Dockerfile               # used in production deployment
│   └── Dockerfile.dev           # for development
├── infra/                       # Terraform infrastructure code
│   ├── ecs/                     # ECS & related resources
│   ├── elasticache/             # Redis caching layer
│   ├── vpc/                     # VPC and networking setup
│   └── private.tfvars           # AWS credentials for Terraform
│   ├── main.tf
│   ├── outputs.tf
│   ├── provider.tf
│   ├── terraform.tfvars       # Shared variables
│   └── variables.tf
├── docker-compose.yml         # Docker multi-service config
├── docker-compose.override.yml # for development purpose used to override the default Compose settings for development
├── .dockerignore
├── .gitignore
├── Makefile                   # Deployment and utility commands
└── README.md
```


---

## 🔧 Environment Variables (.env Example)

Create these environment files based on the examples below.

### Backend (`backend/.env` and `backend/.env.prod`):

```bash
# Database Configuration
DB_DRIVER=postgres
DB_HOST=<your-db-host>
DB_PORT=5432
DB_USER=<your-db-user>
DB_PASSWORD=<your-db-password>
DB_NAME=poll_app
DB_SSLMODE=require

# AWS & SNS Configuration
SNS_ACCESS_KEY=<sns-access-key>
SNS_SECRET_KEY=<sns-secret-key>
SNS_SESSION_TOKEN=
AWS_REGION=ap-southeast-1
ADMIN_TOPIC_ARN=<admin-sns-topic-arn>
USER_TOPIC_ARN=<user-sns-topic-arn>

# Application Config
APP_ENV=dev
SERVER_ADDRESS=:8080

# Redis Configuration
REDIS_HOST=<redis-host>
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Frontend URL (for CORS security)
CORS_ECS_DOMAIN=<frontend-ecs-domain>
```

> Example
```bash
DB_DRIVER=postgres
DB_HOST=<SECRET>
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<SECRET>
DB_NAME=poll_app
DB_SSLMODE=require

SNS_ACCESS_KEY=<SECRET>
SNS_SECRET_KEY=<SECRET>
SNS_SESSION_TOKEN=
AWS_REGION=ap-southeast-1
ADMIN_TOPIC_ARN=arn:aws:sns:ap-southeast-1:<SECRET>:CloudProjAdminNotification
USER_TOPIC_ARN=arn:aws:sns:ap-southeast-1:<SECRET>:CloudProjUserNotification

APP_ENV=dev
SERVER_ADDRESS=:8080

REDIS_HOST=<SECRET>
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

CORS_ECS_DOMAIN=http://cv-c9-alb-267285815.us-west-2.elb.amazonaws.com
```

### Frontend (`frontend/.env` and `frontend/.env.local`):

```bash
# Backend API Endpoint
NEXT_PUBLIC_API_PATH=<backend-api-url>/api
```

> Example
```bash
# Backend API Endpoint
NEXT_PUBLIC_API_PATH=http://cv-c9-alb-267285815.us-west-2.elb.amazonaws.com/api
```


For current configuration, the `frontend-ecs-domain` and `backend-api-url` will  be the same which are `ECS ALB DNS` that you can see in the `ECS` configuration or in the output after finish applying terraform

Both `CORS_ECS_DOMAIN` and `NEXT_PUBLIC_API_PATH` resolve to the ALB DNS. This ensures secure, internal-only communication between frontend and backend.

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

---


### 🧑‍💼 Summary of Key `Makefile` Commands

Here is the table for important commands in `Makefile` so that you can understand more about what it did

| Command                     | Description                                                    |
|-----------------------------|----------------------------------------------------------------|
| `make deploy-all`           | Full deployment of the entire stack, including backend, frontend, and infrastructure. |
| `make tf-init`              | Initializes Terraform and sets up backend state.               |
| `make tf-apply`             | Applies Terraform configurations to create all infrastructure resources. |
| `make tf-output`            | Fetches output values like ALB DNS and Redis endpoint.        |
| `make update-frontend-path` | Updates the frontend with the ALB DNS name for API communication. |
| `make update-cors-domain`   | Configures CORS settings for the backend to accept frontend requests. |
| `make update-ecr`           | Builds and pushes backend/ frontend Docker images to ECR.     |
| `make restart-ecs`          | Forces a new deployment of the ECS services.                  |
| `make tf-destroy`           | Destroys all infrastructure created by Terraform.             |

---

### 🛠 Step-by-Step Explanation of `make deploy-all`

The `make deploy-all` command combines multiple steps to initialize, deploy, and configure the entire infrastructure and application stack. Below is a detailed breakdown of each step.

> 💡 **Note:** You don't have to use `make` if you prefer not to install the `make` utility. All commands used in the Makefile are standard CLI commands — feel free to copy and run them directly from the `Makefile`.

#### 1. 🧑‍💻 **Switch to Production Environment**
   ```bash
   make env-prod
   ```
   - This step ensures that the system is using the production environment variables. It copies the production `.env.prod` file to `.env`, which is used by backend service to connect to the appropriate services and databases in the production environment.

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
   make push-backend && make deploy-ecs
   ```
   This will rebuild the backend Docker image, push it to ECR, and redeploy it to ECS.

- **Redeploy Frontend:**
   ```bash
   make push-frontend && make deploy-ecs
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

## 🗂 Helpful Resources & Documentation

- [Terraform Documentation](https://developer.hashicorp.com/terraform/docs)
- [AWS CLI User Guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-welcome.html)
- [Docker Documentation](https://docs.docker.com)