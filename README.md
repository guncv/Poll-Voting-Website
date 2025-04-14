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
├── backend/              # Go server
├── frontend/             # Next.js app
├── infra/                # Terraform configurations
│   ├── elasticache.tf
│   ├── ecs.tf
│   ├── backend_ecs.tf
│   ├── frontend_ecs.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── terraform.tfvars / private.tfvars
├── Makefile              # Deployment commands
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

---

### 🛠  Step-by-Step Deployment Workflow

#### 1. 🧱 Initialize Terraform

```bash
make tf-init
```

---

### 2. 🌱 Deploy Redis Infrastructure (VPC, Subnets, ElastiCache)

```bash
make deploy-redis
```

This command provisions:
- 🟢 VPC + public and private subnets  
- 🌐 Internet & NAT Gateways  
- 🔒 Redis (ElastiCache) in a secure private subnet

After it's deployed, retrieve the Redis host:

```bash
make tf-output
```

Then copy the `redis_endpoint` and set it inside:

```
backend/.env.prod
```

Example:
```env
REDIS_HOST=your-elasticache-host.amazonaws.com
REDIS_PORT=6379
```

---

### 3. 🛠️ Build & Push Backend Docker Image to ECR

```bash
make build-backend
make push-backend
```

✅ This builds your Go backend and pushes it to **Amazon ECR** using your AWS Account ID and region.

---

### 4. 🚀 Deploy Backend Service to ECS (Fargate)

```bash
make deploy-ecs
```

✅ This sets up:
- ECS Cluster
- Backend ECS Fargate Service
- ALB (Application Load Balancer)
- Target Group and Listener for `/api*` routes

---

### 5. 🌐 Get the Backend URL from ALB

After the ECS deployment completes, run:

```bash
make tf-output
```

Then copy the `alb_dns_name` and paste it in:

```
frontend/.env.prod
```

Example:
```env
NEXT_PUBLIC_API_BASE_URL=http://<alb_dns_name>/api
```

---

### 6. 🛠️ Build & Push Frontend Image to ECR

```bash
make build-frontend
make push-frontend
```

✅ This builds your Next.js app and uploads it to ECR.

---

### 7. 🚀 Deploy Frontend to ECS

```bash
make deploy-ecs
```

✅ This reuses the existing ECS cluster and deploys the frontend container via Fargate behind the same ALB (routing on `/`).

---

## 🔁 Redeploy Only What You Changed

### 🌀 If you change only the backend:

```bash
make build-backend && make push-backend && make deploy-ecs
```

### 🌀 If you change only the frontend:

```bash
make build-frontend && make push-frontend && make deploy-ecs
```

---

## 🔥 Cleanup (Destroy All Infrastructure)

```bash
make tf-destroy
```

This tears down everything — ECS, Redis, subnets, VPC, and related AWS infra.

