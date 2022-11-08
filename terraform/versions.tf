terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    time = {
      source = "hashicorp/time"
    }
    null = {
      source = "hashicorp/null"
    }
    archive = {
      source = "hashicorp/archive"
    }
  }
  backend "s3" {
    bucket = "nugg.xyz-terraform"
    key    = "auth.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      env       = local.env
      stack     = local.app_stack
      namespace = local.rs_mesh_namespace
      app       = local.app
    }
  }
}



