variable "project" {
  description = "Project prefix for resource naming"
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID"
  type        = string
}

variable "public_subnet_ids" {
  description = "Public Subnet IDs"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "Private Subnet IDs"
  type        = list(string)
}

variable "redis_endpoint" {
  description = "Redis Endpoint"
  type        = string
}

variable "redis_port" {
  description = "Redis Port"
  type        = number
}

variable "private_route_table_id" {
  description = "Private Route Table ID"
  type        = string
}

variable "public_route_table_id" {
  description = "Public Route Table ID"
  type        = string
}

