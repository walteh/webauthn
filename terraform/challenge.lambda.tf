resource "aws_ecr_repository" "challenge" {
  name = "${local.app_stack}-challenge-image"
}

resource "aws_lambda_function" "challenge" {
  function_name    = "${local.app_stack}-challenge"
  image_uri        = "${aws_ecr_repository.challenge.repository_url}:${local.latest}"
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
  /* tracing_config { mode = "Active" } */
}

data "archive_file" "challenge" {
  type        = "zip"
  source_dir  = "../challenge"
  excludes    = ["../challenge/bin/**"]
  output_path = "bin/challenge.zip"
}

resource "null_resource" "challenge" {
  triggers = {
    src_hash = "${data.archive_file.challenge.output_sha}"
  }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com
           cd ${path.module}/../challenge
		   docker build --platform=linux/arm64 -t ${aws_ecr_repository.challenge.repository_url}:${local.latest} .
           docker push ${aws_ecr_repository.challenge.repository_url}:${local.latest}
       EOF
  }
}

data "aws_ecr_image" "challenge" {
  depends_on = [
    null_resource.challenge
  ]
  repository_name = aws_ecr_repository.challenge.name
  image_tag       = local.latest
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

