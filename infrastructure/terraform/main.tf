variable "employee_id" {
  default = "01541157"
}

variable "environment" {
  default = "sandbox"
}

variable "region" {
  default = "us-east-2"
}

variable "service_name" {
  default = "tech-search"
}

variable "redis" {

}

variable "vpc_id" {

}

data "aws_caller_identity" "self" {
}

terraform {
  required_version = "~> 1.3.0"
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  backend "s3" {
    bucket   = "fr-koichi-sandbox"
    key      = "terraform/terraform.tfstate.tech_search"
    region   = "ap-northeast-1"
    profile  = "frit"
    role_arn = "arn:aws:iam::882275506731:role/Switch-from-599453524280-frit-whitelist-user-2"
  }
}

provider "aws" {
  region = "us-east-2"
  assume_role {
    role_arn = "arn:aws:iam::882275506731:role/Switch-from-599453524280-frit-whitelist-user-2"
  }
  default_tags {
    tags = {
      Env         = var.environment
      Region      = var.region
      employee_id = var.employee_id
    }
  }
}