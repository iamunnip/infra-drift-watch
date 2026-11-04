terraform {
  required_version = ">= 1.13.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.71.0"
    }
  }
  # Configure an existing backend separately. Do not commit state to Git.
}

variable "region" {
  type    = string
  default = null
}

variable "bucket_name" {
  type        = string
  description = "An existing or separately provisioned, globally unique demo bucket name"
}

provider "aws" {
  region = var.region
}

resource "aws_s3_bucket" "demo" {
  bucket = var.bucket_name
  tags = {
    Project     = "infra-drift-watch"
    Environment = "development"
  }
}
