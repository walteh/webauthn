resource "null_resource" "apple_passkey_register" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      dir = local.apple_passkey_register_dir
      tag = "${aws_ecr_repository.core.repository_url}:${local.apple_passkey_register_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}

data "aws_ecr_image" "apple_passkey_register" {
  depends_on      = [null_resource.apple_passkey_register]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.apple_passkey_register_tag
}


resource "aws_lambda_function" "apple_passkey_register" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.apple_passkey_register
  ]

  function_name    = "${local.app_stack}-apple-passkey-register"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.apple_passkey_register_tag}"
  role             = aws_iam_role.apple_passkey_register.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.apple_passkey_register.image_digest, "sha256:")

  environment {
    variables = {
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



resource "aws_iam_role" "apple_passkey_register" {
  name               = "${local.app_stack}-apple-passkey-register-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-apple-passkey-register-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.apple_passkey_register_lambda_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
}


data "aws_iam_policy_document" "apple_passkey_register_lambda_inline" {
  /* statement {
    effect    = "Allow"
    actions   = ["secretsmanager:GetSecretValue"]
    resources = [aws_secretsmanager_secret.apple_signinwithapple_privatekey.arn]
  } */

  statement {
    effect    = "Allow"
    actions   = ["dynamodb:TransactWriteItems", "dynamodb:TransactGetItems"]
    resources = [aws_dynamodb_table.credentials.arn, aws_dynamodb_table.users.arn, aws_dynamodb_table.ceremonies.arn]
  }
}

/*/////////////////////////
		API Gateway
////////////////////////*/

resource "aws_apigatewayv2_route" "apple_passkey_register" {
  api_id             = aws_apigatewayv2_api.auth.id
  route_key          = "POST /apple/passkey/register"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.apple_passkey_register_lambda.id}"
}

resource "aws_lambda_permission" "apple_passkey_register" {
  statement_id  = "${local.app_stack}-apple-passkey-register-AllowExecutionFromApiGatewayRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.apple_passkey_register.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.auth.execution_arn}/*/*/register"
}

resource "aws_apigatewayv2_integration" "apple_passkey_register_lambda" {
  api_id                 = aws_apigatewayv2_api.auth.id
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.apple_passkey_register.invoke_arn
  payload_format_version = "2.0"
}
