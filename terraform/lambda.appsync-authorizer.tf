resource "null_resource" "appsync_authorizer" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      dir = local.appsync_dir
      tag = "${aws_ecr_repository.core.repository_url}:${local.appsync_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}

data "aws_ecr_image" "appsync_authorizer" {
  depends_on      = [null_resource.appsync_authorizer]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.appsync_tag
}

resource "aws_lambda_function" "appsync_authorizer" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.appsync_authorizer
  ]

  function_name    = "${local.app_stack}-appsync-authorizer"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.appsync_tag}"
  role             = aws_iam_role.appsync_authorizer.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.appsync_authorizer.image_digest, "sha256:")

  environment {
    variables = {
      COGNITO_DEVELOPER_PROVIDER_NAME    = aws_cognito_identity_pool.main.developer_provider_name
      DYNAMO_CEREMONIES_TABLE_NAME       = aws_dynamodb_table.ceremonies.name
      DYNAMO_USERS_TABLE_NAME            = aws_dynamodb_table.users.name
      DYNAMO_CREDENTIALS_TABLE_NAME      = aws_dynamodb_table.credentials.name
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

resource "aws_iam_role" "appsync_authorizer" {
  name               = "${local.app_stack}-appsync-authorizer-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-apple-appsync-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.appsync_authorizer_inline.json
  }
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
}

data "aws_iam_policy_document" "appsync_authorizer_inline" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:GetItem"]
    resources = [aws_dynamodb_table.credentials.arn]
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

