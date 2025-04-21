variable "project" {
  description = "Project prefix for resource naming"
  type        = string
}

variable "aws_region" {
  description = "AWS region to deploy resources"
  type        = string
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
}

