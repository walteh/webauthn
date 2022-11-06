terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.23.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9.0"
    }
  }
  backend "s3" {
    bucket = "nugg.xyz-terraform"
    key    = "crypto.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      env       = local.env
      stack     = local.app_stack
      namespace = local.namespace
      app       = local.app
    }
  }
}
