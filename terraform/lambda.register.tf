resource "null_resource" "register" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      dir = local.register_dir
      tag = "${aws_ecr_repository.core.repository_url}:${local.register_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}

data "aws_ecr_image" "register" {
  depends_on      = [null_resource.register]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.register_tag
}


resource "aws_lambda_function" "register" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.register
  ]

  function_name    = "${local.app_stack}-register"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.register_tag}"
  role             = aws_iam_role.register.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.register.image_digest, "sha256:")

  environment {
    variables = {
      DYNAMO_USER_TABLE_NAME             = aws_dynamodb_table.user.name
      COGNITO_IDENTITY_POOL_ID           = aws_cognito_identity_pool.main.id
      APPLE_PUBLICKEY_ENDPOINT           = "https://appleid.apple.com/auth/keys"
      APPLE_TOKEN_ENDPOINT               = "https://appleid.apple.com/auth/token"
      SM_SIGNINWITHAPPLE_PRIVATEKEY_NAME = aws_secretsmanager_secret.apple_signinwithapple_privatekey.name
      APPLE_TEAM_ID                      = local.apple_team_id
      SIGNIN_WITH_APPLE_PRIVATE_KEY_ID   = local.apple_key_id
      APPLE_SERVICE_NAME                 = local.apple_service_name
    }
  }
}



resource "aws_iam_role" "register" {
  name               = "${local.app_stack}-register-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-register-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.register_lambda_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
}


data "aws_iam_policy_document" "register_lambda_inline" {
  statement {
    effect    = "Allow"
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.apple_signinwithapple_privatekey.arn]
  }

  statement {
    effect    = "Allow"
    actions   = ["dynamodb:PutItem"]
    resources = [aws_dynamodb_table.user.arn]
  }
}

