
resource "aws_ecr_repository" "apple" {
  name = "${local.app_stack}-apple-image"
}

resource "aws_lambda_function" "apple" {
  function_name    = "${local.app_stack}-apple"
  image_uri        = "${aws_ecr_repository.apple.repository_url}:${local.latest}"
  role             = aws_iam_role.apple.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.apple.image_digest, "sha256:")
  environment {
    variables = {
      DYNAMO_CHALLENGE_TABLE_NAME        = aws_dynamodb_table.challenge.name
      COGNITO_IDENTITY_POOL_ID           = aws_cognito_identity_pool.main.id
      APPLE_PUBLICKEY_ENDPOINT           = "https://appleid.apple.com/auth/keys"
      APPLE_TOKEN_ENDPOINT               = "https://appleid.apple.com/auth/token"
      SM_SIGNINWITHAPPLE_PRIVATEKEY_NAME = aws_secretsmanager_secret.apple_signinwithapple_privatekey.name
      APPLE_TEAM_ID                      = local.apple_team_id
      APPLE_KEY_ID                       = local.apple_key_id
      APPLE_SERVICE_NAME                 = local.apple_service_name
    }
  }
  /* tracing_config { mode = "Active" } */
  depends_on = [
    aws_ecr_repository.apple,
    data.aws_ecr_image.apple
  ]
}

data "archive_file" "apple" {
  type        = "zip"
  source_dir  = "../apple"
  excludes    = ["../apple/bin/**"]
  output_path = "bin/apple.zip"
}

resource "null_resource" "apple" {
  triggers = {
    src_hash = "${data.archive_file.apple.output_sha}"
  }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com
           cd ${path.module}/../apple
		   docker build --platform=linux/arm64 -t ${aws_ecr_repository.apple.repository_url}:${local.latest} .
           docker push ${aws_ecr_repository.apple.repository_url}:${local.latest}
       EOF
  }
}

data "aws_ecr_image" "apple" {
  depends_on = [
    null_resource.apple
  ]
  repository_name = aws_ecr_repository.apple.name
  image_tag       = local.latest
}

resource "aws_iam_role" "apple" {
  name               = "${local.app_stack}-apple-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-apple-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.apple_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
  ]
}


data "aws_iam_policy_document" "apple_inline" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:GetItem"]
    resources = [aws_dynamodb_table.challenge.arn]
  }

  statement {
    effect    = "Allow"
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.apple_signinwithapple_privatekey.arn]
  }

  statement {
    effect = "Allow"
    actions = [
      "cognito-identity:GetCredentialsForIdentity",
      "cognito-identity:GetId",
    ]
    resources = [aws_cognito_identity_pool.main.arn]
  }
}

