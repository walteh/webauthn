resource "null_resource" "passkey_assert" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      dir = local.passkey_assert_dir
      tag = "${aws_ecr_repository.core.repository_url}:${local.passkey_assert_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}

data "aws_ecr_image" "passkey_assert" {
  depends_on      = [null_resource.passkey_assert]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.passkey_assert_tag
}


resource "aws_lambda_function" "passkey_assert" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.passkey_assert
  ]

  function_name    = "${local.app_stack}-passkey-assert"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.passkey_assert_tag}"
  role             = aws_iam_role.passkey_assert.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.passkey_assert.image_digest, "sha256:")

  environment {
    variables = {
      COGNITO_DEVELOPER_PROVIDER_NAME  = aws_cognito_identity_pool.main.developer_provider_name
      DYNAMO_CEREMONIES_TABLE_NAME     = aws_dynamodb_table.ceremonies.name
      DYNAMO_USERS_TABLE_NAME          = aws_dynamodb_table.users.name
      DYNAMO_CREDENTIALS_TABLE_NAME    = aws_dynamodb_table.credentials.name
      COGNITO_IDENTITY_POOL_ID         = aws_cognito_identity_pool.main.id
      APPLE_PUBLICKEY_ENDPOINT         = "https://appleid.apple.com/auth/keys"
      APPLE_TOKEN_ENDPOINT             = "https://appleid.apple.com/auth/token"
      APPLE_TEAM_ID                    = local.apple_team_id
      SIGNIN_WITH_APPLE_PRIVATE_KEY_ID = local.apple_key_id
      APPLE_SERVICE_NAME               = local.apple_service_name
    }
  }
}



resource "aws_iam_role" "passkey_assert" {
  name               = "${local.app_stack}-passkey-assert-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume.json
  inline_policy {
    name   = "${local.app_stack}-passkey-assert-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.passkey_assert_lambda_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
  ]
}


data "aws_iam_policy_document" "passkey_assert_lambda_inline" {
  statement {
    effect    = "Allow"
    actions   = ["cognito-identity:GetOpenIdTokenForDeveloperIdentity"]
    resources = [aws_cognito_identity_pool.main.arn]
  }

  statement {
    effect    = "Allow"
    actions   = ["dynamodb:TransactWriteItems", "dynamodb:TransactGetItems", "dynamodb:GetItem", "dynamodb:PutItem", "dynamodb:UpdateItem"]
    resources = [aws_dynamodb_table.credentials.arn, aws_dynamodb_table.users.arn, aws_dynamodb_table.ceremonies.arn]
  }
}

/*/////////////////////////
		API Gateway
////////////////////////*/

resource "aws_apigatewayv2_route" "passkey_assert" {
  api_id             = aws_apigatewayv2_api.auth.id
  route_key          = "POST /passkey/assert"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.passkey_assert_lambda.id}"
}

resource "aws_lambda_permission" "passkey_assert" {
  statement_id  = "${local.app_stack}-passkey-assert-AllowExecutionFromApiGatewayRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.passkey_assert.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.auth.execution_arn}/*/*/assert"
}

resource "aws_apigatewayv2_integration" "passkey_assert_lambda" {
  api_id                 = aws_apigatewayv2_api.auth.id
  integration_type       = "AWS_PROXY"
  integration_method     = "POST"
  integration_uri        = aws_lambda_function.passkey_assert.invoke_arn
  payload_format_version = "2.0"
}
