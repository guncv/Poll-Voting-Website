variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "us-west-2"
}

variable "availability_zones" {
  type    = list(string)
  default = ["us-west-2a", "us-west-2b", "us-west-2c"]
}

variable "aws_account_id" {
  description = "AWS Account ID"
  type        = string
}

variable "aws_access_key" {
  description = "AWS Access Key"
  type        = string
  sensitive   = true
}

variable "aws_secret_key" {
  description = "AWS Secret Key"
  type        = string
  sensitive   = true
}

variable "project" {
  description = "Project prefix for resource naming"
  type        = string
  default     = "cv-c9"
}

variable "env" {
  description = "Environment name (dev or prod)"
  type        = string
  default     = "dev"
}
