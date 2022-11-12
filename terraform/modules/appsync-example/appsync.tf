variable "appsync_authorizer_function_name" { type = string }

terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "aws" {
  region = "us-east-1"
  default_tags {
    tags = {
      stack = local.app_stack
    }
  }
}


data "aws_lambda_function" "authorizer" {
  function_name = var.appsync_authorizer_function_name
}

data "aws_caller_identity" "current" {}

resource "random_string" "random" {
  length           = 16
  special          = true
  override_special = "/@£$"
}

locals {
  app_stack = random_string.random.result
}


output "appsync_graphql_api_endpoint" {
  value = aws_appsync_graphql_api.appsync.uris["GRAPHQL"]
}

resource "aws_appsync_graphql_api" "appsync" {
  authentication_type = "AWS_IAM"
  name                = "${local.app_stack}-auth-tmp-GraphQLAPI"
  schema              = file("${path.module}/schema.graphql")
  log_config {
    cloudwatch_logs_role_arn = aws_iam_role.appsync.arn
    field_log_level          = "ALL"
    exclude_verbose_content  = false
  }
  additional_authentication_provider {
    authentication_type = "AWS_LAMBDA"
    lambda_authorizer_config {
      authorizer_uri = data.aws_lambda_function.authorizer.arn
    }
  }
}

resource "aws_lambda_permission" "apple" {
  statement_id  = "${local.app_stack}-auth-tmp-authorizer-AllowExecutionFromAppSync"
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.apple_authorizer.function_name
  principal     = "appsync.amazonaws.com"
  source_arn    = aws_appsync_graphql_api.appsync.arn
}

resource "aws_iam_role" "appsync" {
  name               = "${local.app_stack}-auth-tmp-AppsyncExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.appsync_assume.json
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSAppSyncPushToCloudWatchLogs",
  ]
}

data "aws_iam_policy_document" "appsync_assume" {
  version = "2012-10-17"
  statement {
    sid     = ""
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["appsync.amazonaws.com"]
    }
  }
}

resource "aws_appsync_resolver" "listen" {
  api_id            = aws_appsync_graphql_api.appsync.id
  type              = "Subscription"
  field             = "listen"
  data_source       = aws_appsync_datasource.none.name
  kind              = "UNIT"
  request_template  = file("${path.module}/velocity/listen.req.vtl")
  response_template = file("${path.module}/velocity/listen.res.vtl")
}

resource "aws_appsync_resolver" "send" {
  api_id            = aws_appsync_graphql_api.appsync.id
  type              = "Mutation"
  field             = "send"
  data_source       = aws_appsync_datasource.none.name
  kind              = "UNIT"
  request_template  = file("${path.module}/velocity/send.req.vtl")
  response_template = file("${path.module}/velocity/send.res.vtl")
}

resource "aws_appsync_datasource" "none" {
  api_id = aws_appsync_graphql_api.appsync.id
  name   = replace("${local.app_stack}-auth-tmp-none", "-", "_")
  type   = "NONE"
}

resource "aws_iam_policy" "current" {
  name   = "${local.app_stack}-auth-tmp-CurrentUserAppsyncExecutionRole"
  users  = [data.aws_caller_identity.current.user_id]
  policy = data.aws_iam_policy_document.current_user_send.json
}

data "aws_iam_policy_document" "current_user_send" {
  statement {
    effect  = "Allow"
    actions = ["appsync:GraphQL"]
    resources = [
      "${aws_appsync_graphql_api.appsync.arn}/types/Mutation/fields/send",
    ]
  }
}

