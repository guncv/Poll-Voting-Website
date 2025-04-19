variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
  default     = "ap-southeast-1"  # Singapore
}
variable "availability_zones" {
  description = "AZs used for subnet placement in Singapore region"
  type        = list(string)
  default     = ["ap-southeast-1a", "ap-southeast-1b"]
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
