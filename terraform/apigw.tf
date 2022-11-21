
resource "aws_apigatewayv2_api" "auth" {
  name                       = "${local.app_stack}-exp-api"
  protocol_type              = "HTTP"
  route_selection_expression = "$request.method $request.path"
  cors_configuration {
    allow_origins  = ["*"]
    allow_methods  = ["POST"]
    allow_headers  = ["*"]
    expose_headers = ["*"]
  }
}

resource "aws_apigatewayv2_authorizer" "apple" {
  api_id                            = aws_apigatewayv2_api.auth.id
  authorizer_type                   = "REQUEST"
  authorizer_uri                    = aws_lambda_function.apigw_authorizer.invoke_arn
  identity_sources                  = ["$request.header.Authorization"]
  name                              = "${local.app_stack}-apple-apigw-authorizer"
  authorizer_payload_format_version = "2.0"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.auth.id
  name        = "default"
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.stack.arn
    format          = "$context.identity.sourceIp - $context.identity.caller - $context.identity.user [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage $context.error.message "
  }
}

resource "aws_apigatewayv2_api_mapping" "auth" {
  api_id          = aws_apigatewayv2_api.auth.id
  domain_name     = local.rs_mesh_apigw_host
  stage           = aws_apigatewayv2_stage.default.id
  api_mapping_key = local.app
}

resource "aws_apigatewayv2_deployment" "default" {
  depends_on = [
    aws_apigatewayv2_route.apple_passkey_attest,
    aws_apigatewayv2_route.init,
    aws_apigatewayv2_route.apple_passkey_assert,
  ]
  api_id = aws_apigatewayv2_api.auth.id
  lifecycle {
    create_before_destroy = true
  }
}


/* resource "aws_iam_role" "complete_integration_role" {
  name               = "${local.app_stack}-complete-IntegrationRole"
  assume_role_policy = data.aws_iam_policy_document.apigw_assume.json
  inline_policy {
    name   = "${local.app_stack}-complete-IntegrationRolePolicy"
    policy = data.aws_iam_policy_document.complete_integration_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
  ]
}

data "aws_iam_policy_document" "apigw_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "complete_integration_inline" {
  statement {
    effect    = "Allow"
    actions   = ["lambda:InvokeFunction"]
    resources = [aws_lambda_function.complete.arn]
  }
} */
