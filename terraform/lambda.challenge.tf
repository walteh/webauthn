resource "null_resource" "challenge" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      cmd = local.challenge_cmd
      tag = "${aws_ecr_repository.core.repository_url}:${local.challenge_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}

data "aws_ecr_image" "challenge" {
  depends_on      = [null_resource.challenge]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.challenge_tag
}


resource "aws_lambda_function" "challenge" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.challenge
  ]

  function_name    = "${local.app_stack}-challenge"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.challenge_tag}"
  role             = aws_iam_role.challenge.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.challenge.image_digest, "sha256:")

  environment {
    variables = {
      DYNAMO_CHALLENGE_TABLE_NAME = aws_dynamodb_table.challenge.name
    }
  }
}



resource "aws_iam_role" "challenge" {
  name               = "${local.app_stack}-challenge-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-challenge-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.challenge_lambda_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess",
  ]
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "challenge_lambda_inline" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.challenge.arn]
  }
}

