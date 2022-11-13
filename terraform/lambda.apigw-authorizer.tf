


resource "null_resource" "apigw_authorizer" {
  triggers = { src_hash = "${data.archive_file.core.output_sha}" }
  provisioner "local-exec" {
    environment = {
      cmd = local.apigw_cmd
      tag = "${aws_ecr_repository.core.repository_url}:${local.apigw_tag}"
    }
    command = local.lambda_docker_deploy_command
  }
}


data "aws_ecr_image" "apigw_authorizer" {
  depends_on      = [null_resource.apigw_authorizer]
  repository_name = aws_ecr_repository.core.name
  image_tag       = local.apigw_tag
}

resource "aws_lambda_function" "apigw_authorizer" {
  depends_on = [
    aws_ecr_repository.core,
    data.aws_ecr_image.apigw_authorizer
  ]

  function_name    = "${local.app_stack}-apigw-authorizer"
  image_uri        = "${aws_ecr_repository.core.repository_url}:${local.apigw_tag}"
  role             = aws_iam_role.apigw_authorizer.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.apigw_authorizer.image_digest, "sha256:")

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

}

resource "aws_iam_role" "apigw_authorizer" {
  name               = "${local.app_stack}-apigw-authorizer-ExecutionRole"
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

