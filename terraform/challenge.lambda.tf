
resource "aws_ecr_repository" "challenge" {
  name = "${local.app_stack}-challenge-image"
}

resource "null_resource" "challenge" {
  depends_on = [null_resource.docker]
  provisioner "local-exec" {
    environment = {
      tag = "${aws_ecr_repository.challenge.repository_url}:${local.latest}"
    }
    command = <<EOF
			cd ${path.module}/../apple
		    docker build --platform=linux/arm64 --target challenge -t $tag .
			docker push $tag
		EOF
  }
}

data "aws_ecr_image" "challenge" {
  depends_on      = [null_resource.challenge]
  repository_name = aws_ecr_repository.challenge.name
  image_tag       = local.latest
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
